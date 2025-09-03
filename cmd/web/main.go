// main.go contains our initalization processes and url patterns.
// This initializes our server loads in any of our .env files, establishes our connection to the database,
// allows for our DB models to have access to our DB connection for queries, and contains our URL patterns/handlers.

package main

import (
	"database/sql"
	"ezkitchen/internal/models"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	logger        *slog.Logger
	estimates     *models.EstimateModel
	users         *models.UserModel
	templateCache map[string]*template.Template
}

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
	}

	// executing flags and loading environment variables.
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	flag.Parse()

	err := godotenv.Load("secrets.env")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	//Database connection
	psqlStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("HOST"), os.Getenv("PORT"),
		os.Getenv("DBUSER"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"))

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

	app.estimates = &models.EstimateModel{DB: db}
	app.users = &models.UserModel{DB: db}
	app.templateCache = templateCache

	// HTTP Routes
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /estimate/view/{id}", app.estimateView)
	mux.HandleFunc("GET /estimate/create", app.estimateCreate)
	mux.HandleFunc("GET /estimate/edit/{id}", app.estimateEditView)

	mux.HandleFunc("POST /estimate/create", app.estimateCreatePost)
	mux.HandleFunc("POST /estimate/update", app.estimateUpdate)

	mux.HandleFunc("DELETE /estimate/delete/{id}", app.estimateDelete)

	logger.Info("Starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, mux)
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
