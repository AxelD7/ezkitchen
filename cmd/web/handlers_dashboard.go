package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	id := app.sessionManager.Get(r.Context(), "authenticatedUserID")

	app.logger.Info(fmt.Sprintf("CURRENT USER ID %v", id))

	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
