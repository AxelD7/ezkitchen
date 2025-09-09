package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// routes contains all handlerfunc paths in our app.
func (app *application) routes() http.Handler {
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

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
