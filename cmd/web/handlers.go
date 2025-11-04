package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"ezkitchen/internal/models"
	"ezkitchen/internal/validator"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type productRequestBody struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
	validator.Validator
}

type estimateCreateForm struct {
	Name          string
	StreetAddress string
	City          string
	State         string
	Zip           string
	Email         string
	Phone         string
	Length        float32
	Width         float32
	Height        float32
	DoorWidth     float32
	DoorHeight    float32
	validator.Validator
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

	fmt.Printf("ESTIMATE OBJECT: %+v\n", estimate)
}

func (app *application) estimateCreate(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, http.StatusOK, "createEstimate.tmpl", templateData{
		Form: estimateCreateForm{},
	})
}

func (app *application) estimateCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
	}

	form := estimateCreateForm{
		Name:          r.PostForm.Get("customerName"),
		StreetAddress: r.PostForm.Get("streetAddress"),
		City:          r.PostForm.Get("city"),
		State:         r.PostForm.Get("state"),
		Zip:           r.PostForm.Get("zip"),
		Email:         r.PostForm.Get("email"),
		Phone:         r.PostForm.Get("phone"),
		Length:        app.formFloat32Parse(r, "kitchenLength"),
		Width:         app.formFloat32Parse(r, "kitchenWidth"),
		Height:        app.formFloat32Parse(r, "kitchenHeight"),
		DoorWidth:     app.formFloat32Parse(r, "doorwayWidth"),
		DoorHeight:    app.formFloat32Parse(r, "doorwayHeight"),
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
		data := templateData{
			Form: form,
		}
		app.render(w, r, http.StatusUnprocessableEntity, "createEstimate.tmpl", data)
		return
	}

	customer := models.User{
		Name:         r.PostForm.Get("customerName"),
		Email:        r.PostForm.Get("email"),
		PasswordHash: "",
		Phone:        r.PostForm.Get("phone"),
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

	err = app.estimates.Insert(&estimate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

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

	app.render(w, r, http.StatusOK, "editEstimate.tmpl", templateData{
		Estimate:       estimate,
		Customer:       customer,
		Products:       estimateProducts,
		EstimateTotals: estimateTotals,
	})

	app.logger.Info(fmt.Sprintf("Viewing and editting the estimate with id %v", estimate.EstimateID))

}

func (app *application) estimateUpdate(w http.ResponseWriter, r *http.Request) {

	estimate := models.Estimate{
		CustomerID:        2,
		CreatedBy:         2,
		Status:            models.StatusDraft,
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

	ts, ok := app.templateCache["addLineItemModal.tmpl"]
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
