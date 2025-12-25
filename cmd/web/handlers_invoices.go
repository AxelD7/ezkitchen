package main

import (
	"ezkitchen/internal/models"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (app *application) signInvoiceView(w http.ResponseWriter, r *http.Request) {

	rawToken := r.URL.Query().Get("token")
	if rawToken == "" {
		app.render(w, r, http.StatusGone, "invalidInvoice.tmpl", templateData{})
		return
	}

	it, err := app.invoiceToken.GetByRawToken(rawToken)
	if err != nil {
		app.render(w, r, http.StatusGone, "invalidInvoice.tmpl", templateData{})
		return
	}

	if time.Now().After(it.ExpiresAt) || it.UsedAt.Valid {
		app.render(w, r, http.StatusGone, "invalidInvoice.tmpl", templateData{})
		return
	}

	estimate, err := app.estimates.Get(it.EstimateID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	estimateProducts, err := app.estimateItems.GetByEstimateID(estimate.EstimateID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	estimateTotals := app.estimates.CalculateEstimateTotals(estimateProducts)

	customer, err := app.users.Get(estimate.CustomerID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Estimate = estimate
	data.Customer = customer
	data.Products = estimateProducts
	data.EstimateTotals = estimateTotals
	data.Token = rawToken

	fmt.Printf("ESTIMATE OBJECT: %+v\n", estimate)

	app.render(w, r, http.StatusOK, "customerInvoice.tmpl", data)

}
func (app *application) submitSignature(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rawToken := r.FormValue("token")

	it, err := app.invoiceToken.GetByRawToken(rawToken)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		app.render(w, r, http.StatusGone, "invalidInvoice.tmpl", templateData{})
		return
	}

	if time.Now().After(it.ExpiresAt) || it.UsedAt.Valid {
		data := app.newTemplateData(r)
		app.render(w, r, http.StatusGone, "invalidInvoice.tmpl", data)
		return
	}

	estimate, err := app.estimates.Get(it.EstimateID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if estimate.SignatureObjectKey.Valid {
		app.clientError(w, r, http.StatusConflict)
		return
	}

	err = r.ParseMultipartForm(1 << 20)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	file, header, err := r.FormFile("signature")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer file.Close()

	if header.Size > 512*1024 {
		app.clientError(w, r, http.StatusRequestEntityTooLarge)
		return
	}

	if header.Header.Get("Content-Type") != "image/png" {
		app.clientError(w, r, http.StatusUnsupportedMediaType)
		return
	}

	err = app.storage.UploadSignature(ctx, estimate.EstimateID, file, "image/png")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	tx, err := app.estimates.DB.Begin()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer tx.Rollback()

	err = app.estimates.SetSignatureKeyTx(tx, estimate.EstimateID, fmt.Sprintf("signatures/%d.png", estimate.EstimateID))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.estimates.UpdateStatusTx(tx, estimate.EstimateID, models.StatusInProgress)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.invoiceToken.MarkUsedTx(tx, it.InvoiceTokenID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, http.StatusOK, "invoiceSignatureSuccess.tmpl", app.newTemplateData(r))

}

func (app *application) getInvoiceSignature(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	estimateID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || estimateID < 1 {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	estimate, err := app.estimates.Get(estimateID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	if !estimate.SignatureObjectKey.Valid {
		app.clientError(w, r, http.StatusNotFound)
		return
	}

	obj, err := app.storage.Get(ctx, estimate.SignatureObjectKey.String)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	defer obj.Body.Close()

	w.Header().Set("Content-Type", obj.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(obj.Size, 10))
	w.Header().Set("Cache-Control", "private, max-age=3600")

	_, err = io.Copy(w, obj.Body)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

}
