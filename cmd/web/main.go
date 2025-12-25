// main.go contains our initalization processes and url patterns.
// This initializes our server loads in any of our .env files, establishes our connection to the database,
// allows for our DB models to have access to our DB connection for queries, and contains our URL patterns/handlers.

package main

import (
	"context"
	"database/sql"
	"encoding/gob"
	"ezkitchen/internal/mailer"
	"ezkitchen/internal/models"
	"ezkitchen/internal/storage"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/form/v4"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	logger         *slog.Logger
	estimates      *models.EstimateModel
	products       *models.ProductModel
	estimateItems  *models.EstimateItemModel
	users          *models.UserModel
	invoiceToken   *models.InvoiceTokenModel
	storage        *storage.R2Storage
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	mailer         *mailer.Mailer
}

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
	}

	// executing flags and loading environment variables.
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	//Database connection
	psqlStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := openDB(psqlStr)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info(psqlStr)

	defer db.Close()
	// R2 Object Storage
	r2AccessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	r2SecretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	r2Endpoint := os.Getenv("R2_ENDPOINT")
	r2Bucket := os.Getenv("R2_BUCKET")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(r2AccessKeyID, r2SecretKey, "")),
		config.WithRegion("auto"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(r2Endpoint)
	})

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	gob.Register(FlashMessage{})

	mailer, err := mailer.NewMailer()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app.estimates = &models.EstimateModel{DB: db}
	app.products = &models.ProductModel{DB: db}
	app.estimateItems = &models.EstimateItemModel{DB: db}
	app.users = &models.UserModel{DB: db}
	app.invoiceToken = &models.InvoiceTokenModel{DB: db}
	app.storage = storage.NewR2Storage(client, r2Bucket)
	app.templateCache = templateCache
	app.formDecoder = formDecoder
	app.sessionManager = sessionManager
	app.mailer = mailer

	logger.Info("Starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(psqlStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlStr)
	if err != nil {
		db.Close()
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
