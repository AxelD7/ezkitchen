package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	isAuth := app.isAuthenticated(r)
	if isAuth {
		http.Redirect(w, r, "/estimate/list", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}

}
