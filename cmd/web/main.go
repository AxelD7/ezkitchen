package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	logger *slog.Logger
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

	defer db.Close()

	// HTTP Server Start
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)

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
