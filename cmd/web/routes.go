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
	protected := dynamic.Append(app.requireAuthentication)

	// --------------- Estimates ---------------
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /estimate/view/{id}", protected.ThenFunc(app.estimateView))
	mux.Handle("GET /estimate/create", protected.ThenFunc(app.estimateCreateView))
	mux.Handle("GET /estimate/edit/{id}", protected.ThenFunc(app.estimateEditView))
	mux.Handle("POST /estimate/create", protected.ThenFunc(app.estimateCreatePost))
	mux.Handle("POST /estimate/update", protected.ThenFunc(app.estimateUpdate))
	mux.Handle("POST /estimate/{id}/items/", protected.ThenFunc(app.estimateAddItem))
	mux.Handle("POST /estimate/{id}/progress", protected.ThenFunc(app.progressEstimate))
	mux.Handle("PUT /estimate/items/{id}", protected.ThenFunc(app.estimateUpdateItem))
	mux.Handle("DELETE /estimate/items/{id}", protected.ThenFunc(app.estimateDeleteItem))

	mux.Handle("DELETE /estimate/delete/{id}", protected.ThenFunc(app.estimateDelete))

	// --------------- Products ---------------
	mux.Handle("GET /product/get/{id}", protected.ThenFunc(app.productGet))
	mux.Handle("GET /product/get/", protected.ThenFunc(app.fetchProductsByFilters))
	mux.Handle("POST /product/create", protected.ThenFunc(app.productCreate))
	mux.Handle("POST /product/update", protected.ThenFunc(app.productUpdate))
	mux.Handle("DELETE /product/delete", protected.ThenFunc(app.productDelete))

	// --------------- Invoices ---------------

	mux.Handle("GET /invoice/sign", dynamic.ThenFunc(app.signInvoiceView))
	mux.Handle("POST /invoice/sign", dynamic.ThenFunc(app.submitSignature))
	mux.Handle("GET /invoice/signature/{id}", protected.ThenFunc(app.getInvoiceSignature))

	// --------------- Users ---------------
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLoginView))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogout))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders, preventCSRF)

	return standard.Then(mux)
}
