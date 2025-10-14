package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"ezkitchen/internal/models"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type productRequestBody struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
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

	app.render(w, r, http.StatusOK, "createEstimate.tmpl", templateData{})
}

func (app *application) estimateCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
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

	customer, err := app.users.Get(estimate.CustomerID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, http.StatusOK, "editEstimate.tmpl", templateData{
		Estimate: estimate,
		Customer: customer,
		Products: estimateProducts,
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
		UnitPrice:   1200.99,
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
		UnitPrice:   500.99,
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
