package main

import (
	"errors"
	"ezkitchen/internal/models"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

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

	customer, err := app.users.Get(estimate.CustomerID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, http.StatusOK, "editEstimate.tmpl", templateData{
		Estimate: estimate,
		Customer: customer,
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

func (app *application) productDelete(w http.ResponseWriter, r *http.Request) {

	id := 2

	err := app.products.Delete(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	fmt.Printf("PRODUCT WITH ID %v HAS BEEN DELETED", id)
}
