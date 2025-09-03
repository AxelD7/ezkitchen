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
		Name:         r.PostForm.Get("name"),
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
		CustomerID:       customer.UserID,
		CreatedBy:        2,
		Status:           models.StatusDraft,
		CreatedAt:        time.Now(),
		KitchenLengthFt:  app.formFloat32Parse(r, "kitchenLength"),
		KitchenWidthFt:   app.formFloat32Parse(r, "kitchenWidth"),
		KitchenHeightFt:  app.formFloat32Parse(r, "kitchenHeight"),
		DoorWidthInches:  app.formFloat32Parse(r, "doorwayWidth"),
		DoorHeightInches: app.formFloat32Parse(r, "doorwayHeight"),
		HasIsland:        false,
		FlooringType:     r.PostForm.Get("flooringType"),
		Street:           r.PostForm.Get("street"),
		City:             r.PostForm.Get("city"),
		State:            r.PostForm.Get("state"),
		Zip:              r.PostForm.Get("zip"),
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
		CustomerID:       2,
		CreatedBy:        2,
		Status:           models.StatusDraft,
		CreatedAt:        time.Now(),
		KitchenLengthFt:  app.formFloat32Parse(r, "kitchenLength"),
		KitchenWidthFt:   app.formFloat32Parse(r, "kitchenWidth"),
		KitchenHeightFt:  app.formFloat32Parse(r, "kitchenHeight"),
		DoorWidthInches:  app.formFloat32Parse(r, "doorwayWidth"),
		DoorHeightInches: app.formFloat32Parse(r, "doorwayHeight"),
		HasIsland:        false,
		FlooringType:     r.PostForm.Get("flooringType"),
		Street:           r.PostForm.Get("street"),
		City:             r.PostForm.Get("city"),
		State:            r.PostForm.Get("state"),
		Zip:              r.PostForm.Get("zip"),
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
