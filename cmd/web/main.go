// main.go contains our initalization processes and url patterns.
// This initializes our server loads in any of our .env files, establishes our connection to the database,
// allows for our DB models to have access to our DB connection for queries, and contains our URL patterns/handlers.

package main

import (
	"database/sql"
	"encoding/gob"
	"ezkitchen/internal/models"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
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
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	app.estimates = &models.EstimateModel{DB: db}
	app.products = &models.ProductModel{DB: db}
	app.estimateItems = &models.EstimateItemModel{DB: db}
	app.users = &models.UserModel{DB: db}
	app.templateCache = templateCache
	app.formDecoder = formDecoder
	app.sessionManager = sessionManager

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
