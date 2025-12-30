package main

import (
	"ezkitchen/internal/models"
	"fmt"
	"net/http"
)

func (app *application) productCreate(w http.ResponseWriter, r *http.Request) {

	currUser := app.currentUser(r)

	if currUser.Role != models.RoleAdmin {
		app.clientError(w, r, http.StatusNotFound)
		return
	}

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

	currUser := app.currentUser(r)

	if currUser.Role != models.RoleAdmin {
		app.clientError(w, r, http.StatusNotFound)
		return
	}

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

	currUser := app.currentUser(r)

	if currUser.Role != models.RoleAdmin {
		app.clientError(w, r, http.StatusNotFound)
		return
	}

	id := 2

	err := app.products.Delete(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	fmt.Printf("PRODUCT WITH ID %v HAS BEEN DELETED", id)
}
