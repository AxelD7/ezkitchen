package main

import (
	"errors"
	"ezkitchen/internal/models"
	"ezkitchen/internal/validator"
	"fmt"
	"net/http"
)

type userPostForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	validator.Validator
}

func (app *application) userLoginView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userPostForm{}

	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	var form userPostForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.IsValidEmail(form.Email), "email", "You must use a valid email!")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	userID, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {

			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
			return
		} else {
			app.serverError(w, r, err)

		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", userID)

	app.sessionManager.Put(r.Context(), "flash", FlashMessage{
		Type:    "success",
		Message: "You have logged in successfully",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {

	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", FlashMessage{
		Type:    "success",
		Message: "You've been logged out successfully",
	})

	fmt.Println("log out successful")

	http.Redirect(w, r, "/", http.StatusSeeOther)

}
