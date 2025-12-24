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

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// --------------- Estimates ---------------
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /estimate/view/{id}", dynamic.ThenFunc(app.estimateView))
	mux.Handle("GET /estimate/create", dynamic.ThenFunc(app.estimateCreate))
	mux.Handle("GET /estimate/edit/{id}", dynamic.ThenFunc(app.estimateEditView))
	mux.Handle("POST /estimate/create", dynamic.ThenFunc(app.estimateCreatePost))
	mux.Handle("POST /estimate/update", dynamic.ThenFunc(app.estimateUpdate))
	mux.Handle("POST /estimate/{id}/items/", dynamic.ThenFunc(app.estimateAddItem))
	mux.Handle("POST /estimate/{id}/progress", dynamic.ThenFunc(app.progressEstimate))
	mux.Handle("PUT /estimate/items/{id}", dynamic.ThenFunc(app.estimateUpdateItem))
	mux.Handle("DELETE /estimate/items/{id}", dynamic.ThenFunc(app.estimateDeleteItem))

	mux.Handle("DELETE /estimate/delete/{id}", dynamic.ThenFunc(app.estimateDelete))

	// --------------- Products ---------------
	mux.Handle("GET /product/get/{id}", dynamic.ThenFunc(app.productGet))
	mux.Handle("GET /product/get/", dynamic.ThenFunc(app.fetchProductsByFilters))
	mux.Handle("POST /product/create", dynamic.ThenFunc(app.productCreate))
	mux.Handle("POST /product/update", dynamic.ThenFunc(app.productUpdate))
	mux.Handle("DELETE /product/delete", dynamic.ThenFunc(app.productDelete))

	// --------------- Invoices ---------------

	mux.Handle("GET /invoice/sign", dynamic.ThenFunc(app.signInvoiceView))
	mux.Handle("POST /invoice/sign", dynamic.ThenFunc(app.submitSignature))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
