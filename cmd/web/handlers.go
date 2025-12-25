package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"ezkitchen/internal/mailer"
	"ezkitchen/internal/models"
	"ezkitchen/internal/validator"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type productRequestBody struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
	validator.Validator
}

type estimateCreateForm struct {
	Name                string  `form:"customerName"`
	StreetAddress       string  `form:"streetAddress"`
	City                string  `form:"city"`
	State               string  `form:"state"`
	Zip                 string  `form:"zip"`
	Email               string  `form:"email"`
	Phone               string  `form:"phone"`
	Length              float32 `form:"kitchenLength"`
	Width               float32 `form:"kitchenWidth"`
	Height              float32 `form:"kitchenHeight"`
	DoorWidth           float32 `form:"doorwayWidth"`
	DoorHeight          float32 `form:"doorwayHeight"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, http.StatusOK, "home.tmpl", templateData{})
}

func (app *application) estimateView(w http.ResponseWriter, r *http.Request) {

	var estimate models.Estimate

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	estimate, err = app.estimates.Get(id)
	if err != nil {
		app.serverError(w, r, err)
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

	fmt.Printf("ESTIMATE OBJECT: %+v\n", estimate)

	app.render(w, r, http.StatusOK, "viewEstimate.tmpl", data)
}

func (app *application) estimateCreate(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, http.StatusOK, "createEstimate.tmpl", templateData{
		Form: estimateCreateForm{},
	})
}

func (app *application) estimateCreatePost(w http.ResponseWriter, r *http.Request) {
	var form estimateCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "customerName", "This field cannot be blank.")
	form.CheckField(validator.MaxChars(form.Name, 50), "customerName", "This field cannot be more than 50 characters long.")

	form.CheckField(validator.NotBlank(form.StreetAddress), "streetAddress", "This field cannot be blank.")
	form.CheckField(validator.MaxChars(form.StreetAddress, 50), "streetAddress", "This field cannot be more than 50 characters long.")

	form.CheckField(validator.NotBlank(form.City), "city", "This field cannot be blank.")
	form.CheckField(validator.MaxChars(form.City, 30), "city", "This field cannot be more than 30 characters long.")

	form.CheckField(validator.NotBlank(form.Zip), "zip", "This field cannot be blank.")
	form.CheckField(validator.MaxChars(form.Zip, 10), "zip", "This field cannot be more than 10 characters long.")

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank.")
	form.CheckField(validator.IsValidEmail(form.Email), "email", "Email must be in proper format. Ex: john.doe@example.com")

	form.CheckField(validator.NotBlank(form.Phone), "phone", "This field cannot be blank.")
	form.CheckField(validator.NotBlank(form.State), "state", "Please select a state.")

	form.CheckField(validator.GreaterThanN(form.Length, float32(0)), "kitchenLength", "This value cannot be zero")
	form.CheckField(validator.GreaterThanN(form.Width, float32(0)), "kitchenWidth", "This value cannot be zero")
	form.CheckField(validator.GreaterThanN(form.Height, float32(0)), "kitchenHeight", "This value cannot be zero")
	form.CheckField(validator.GreaterThanN(form.DoorWidth, float32(0)), "doorwayWidth", "This value cannot be zero")
	form.CheckField(validator.GreaterThanN(form.DoorHeight, float32(0)), "doorwayHeight", "This value cannot be zero")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "createEstimate.tmpl", data)
		return
	}

	customer := models.User{
		Name:         form.Name,
		Email:        form.Email,
		PasswordHash: "",
		Phone:        form.Phone,
		Role:         models.RoleCustomer,
		CreatedAt:    time.Now(),
	}

	err = app.users.Insert(&customer)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	estimate := models.Estimate{
		CustomerID:        customer.UserID,
		CreatedBy:         1,
		Status:            models.StatusDraft,
		CreatedAt:         time.Now(),
		KitchenLengthInch: form.Length,
		KitchenWidthInch:  form.Width,
		KitchenHeightInch: form.Height,
		DoorWidthInch:     form.DoorWidth,
		DoorHeightInch:    form.DoorHeight,
		Street:            form.StreetAddress,
		City:              form.City,
		State:             form.State,
		Zip:               form.Zip,
	}

	err = app.estimates.Insert(&estimate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", FlashMessage{
		Type:    "success",
		Message: "Estimate creation was successful!",
	})

	http.Redirect(w, r, fmt.Sprintf("/estimate/edit/%d", estimate.EstimateID), http.StatusSeeOther)

	app.logger.Info(fmt.Sprintf("The id of the new estimate is %v", estimate.EstimateID))
}

func (app *application) estimateEditView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	estimate, err := app.estimates.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
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

	app.render(w, r, http.StatusOK, "editEstimate.tmpl", data)

	app.logger.Info(fmt.Sprintf("Viewing and editting the estimate with id %v", estimate.EstimateID))

}

func (app *application) progressEstimate(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	estimate, err := app.estimates.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	switch estimate.Status {
	case models.StatusDraft:
		ei, err := app.estimateItems.GetByEstimateID(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		if len(ei) == 0 {
			app.sessionManager.Put(r.Context(), "flash", FlashMessage{
				Type:    "error",
				Message: "You must add at least one product to the estimate before submitting!",
			})
			http.Redirect(
				w, r,
				fmt.Sprintf("/estimate/edit/%d", estimate.EstimateID),
				http.StatusSeeOther,
			)
			return
		}

		err = app.estimates.UpdateStatus(id, estimate.Status.Next())
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		expiresAt := time.Now().Add(72 * time.Hour)
		rawToken, err := app.invoiceToken.Insert(id, expiresAt)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		customer, err := app.users.Get(estimate.CustomerID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		signURL := os.Getenv("APP_BASE_URL") + "/invoice/sign?token=" + rawToken

		invoiceData := mailer.InvoiceLinkData{
			CustomerName:   customer.Name,
			EstimateNumber: estimate.EstimateID,
			SignURL:        signURL,
			ExpiresAt:      expiresAt.Format("Jan 2, 2006 3:04 PM")}

		err = app.mailer.SendInvoiceLink(customer.Email, invoiceData)
		if err != nil {
			app.logger.Error("email send failed", "error", err)
			app.serverError(w, r, err)
			return
		}

		app.sessionManager.Put(r.Context(), "flash", FlashMessage{
			Type:    "success",
			Message: "Estimate submission was successful!",
		})
	}

	http.Redirect(
		w, r,
		fmt.Sprintf("/estimate/view/%d", estimate.EstimateID),
		http.StatusSeeOther,
	)
}

func (app *application) estimateUpdate(w http.ResponseWriter, r *http.Request) {

	estimate := models.Estimate{
		EstimateID:        1,
		CustomerID:        2,
		CreatedBy:         2,
		Status:            models.StatusPaid,
		CreatedAt:         time.Now(),
		KitchenLengthInch: app.formFloat32Parse(r, "kitchenLength"),
		KitchenWidthInch:  app.formFloat32Parse(r, "kitchenWidth"),
		KitchenHeightInch: app.formFloat32Parse(r, "kitchenHeight"),
		DoorWidthInch:     app.formFloat32Parse(r, "doorwayWidth"),
		DoorHeightInch:    app.formFloat32Parse(r, "doorwayHeight"),
		Street:            r.PostForm.Get("street"),
		City:              r.PostForm.Get("city"),
		State:             r.PostForm.Get("state"),
		Zip:               r.PostForm.Get("zip"),
	}

	err := app.estimates.Update(&estimate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info(fmt.Sprintf("The estimate with id %v has been updated", 2))

}

func (app *application) estimateDelete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
	}

	err = app.estimates.Delete(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info(fmt.Sprintf("The estimate with id %v has been deleted", id))

}

func (app *application) productCreate(w http.ResponseWriter, r *http.Request) {

	p := models.Product{
		Name:        "32 Inch LG Fridge",
		Description: "A fridge from lg that is 32 inches",
		UnitPrice:   120099,
		CreatedBy:   2,
	}

	err := app.products.Insert(&p)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	fmt.Printf("product inserted, new ID is : %v", p.ProductID)

}

func (app *application) productUpdate(w http.ResponseWriter, r *http.Request) {

	p := models.Product{
		ProductID:   1,
		Name:        "28 Inch GE Fridge",
		Description: "A fridge from GE that is 28 inches",
		UnitPrice:   50099,
		CreatedBy:   2,
	}

	err := app.products.Update(&p)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	fmt.Printf("product updated, new ID is : %v", p.ProductID)

}

func (app *application) productGet(w http.ResponseWriter, r *http.Request) {

	id := 2
	var p models.Product
	p, err := app.products.Get(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	fmt.Printf("GOT PRODUCT ! : %+v", p)
}

func (app *application) fetchProductsByFilters(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	category := queryParams.Get("category")
	subcategory := queryParams.Get("subcategory")
	color := queryParams.Get("color")

	products, err := app.products.GetByProductFilter(category, subcategory, color)
	if err != nil {
		app.serverError(w, r, err)
	}
	var buf bytes.Buffer

	ts, ok := app.templateCache["modals/addLineItemModal.tmpl"]
	if !ok {
		app.logger.Error("the template addLineItemModal.tmpl does not exist")
		http.Error(w, `{"status": "error", "message": "template not found"}`, http.StatusInternalServerError)
	}

	ts.ExecuteTemplate(&buf, "addLineItemModal", products)

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())

}

func (app *application) estimateAddItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	var req productRequestBody
	json.NewDecoder(r.Body).Decode(&req)

	item := &models.EstimateItem{
		EstimateID: id,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
	}

	estimate, err := app.estimates.Get(item.EstimateID)
	if err != nil {
		app.serverError(w, r, err)

	}

	product, err := app.products.Get(item.ProductID)
	if err != nil {
		app.serverError(w, r, err)
	}

	req.CheckField(validator.GreaterThanN(item.Quantity, 0), "quantity", "The quantity must be at least 1")

	req.CheckField(
		(product.Width <= estimate.DoorWidthInch+1 && product.Height <= estimate.DoorHeightInch+1) ||
			(product.Length <= estimate.DoorWidthInch+1 && product.Height <= estimate.DoorHeightInch+1) ||
			(product.Width <= estimate.DoorHeightInch+1 && product.Length <= estimate.DoorWidthInch+1),
		"product",
		"Product must have at least an one inch clearance of doorway width and height to fit through the doorway.",
	)

	if !req.Valid() {
		app.failedValidationJSON(w, req.FieldErrors)
		return
	}

	err = app.estimateItems.Insert(item)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (app *application) estimateUpdateItem(w http.ResponseWriter, r *http.Request) {
	lineItemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid line item id", http.StatusBadRequest)
		return
	}
	var req productRequestBody
	json.NewDecoder(r.Body).Decode(&req)

	estimateItem := models.EstimateItem{
		LineItemID: lineItemID,
		Quantity:   req.Quantity,
	}

	req.CheckField(validator.GreaterThanN(estimateItem.Quantity, 0), "quantity", "The quantity must be at least 1")

	if !req.Valid() {
		app.failedValidationJSON(w, req.FieldErrors)
		return
	}

	err = app.estimateItems.Update(estimateItem)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) estimateDeleteItem(w http.ResponseWriter, r *http.Request) {
	lineItemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || lineItemID < 1 {
		http.Error(w, "invalid line item id", http.StatusBadRequest)
		return
	}

	err = app.estimateItems.Delete(lineItemID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (app *application) productDelete(w http.ResponseWriter, r *http.Request) {

	id := 2

	err := app.products.Delete(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	fmt.Printf("PRODUCT WITH ID %v HAS BEEN DELETED", id)
}

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

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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
