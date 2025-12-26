// helper.go contains our error handling to distinguish between our client side or server side errors.
// this also contains any of our data type parsers to convert our html form data.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"ezkitchen/internal/models"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

// ---ERROR HANDLING HELPERS---

// serverError takes in the Response Writer, Request and Error to return an Internal Server Error Status to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError takes in the Response Writer, Request and Status to return any distinct error caused by the client.
func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) failedValidationJSON(w http.ResponseWriter, errors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]any{
		"errors": errors,
	})
}

// ---FORM PARSING HELPERS---

// formFloat32Parse takes in our request and a field string for the html form when a post/update is called.
// This will convert the value into a type of float32
// This will log an error when the field contains a unparsable value from the field.
func (app *application) formFloat32Parse(r *http.Request, field string) float32 {
	val := r.PostForm.Get(field)
	if val == "" {
		return 0
	}

	valF64, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Printf("Invalid float for field %s: %v", field, err)
		return 0
	}

	return float32(valF64)

}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecodeError *form.InvalidDecoderError

		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}

		return err
	}

	return nil
}

// render is a function to render our templates(html) and any templateData when called in a handler function.
// this function does a test render before actually writing to the response writer in the event of any errors they get
// handled.
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

// isAuthenticated is a function to check if the session from the request has the
// "authenticatedUserID"
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}

// currentUser is a helper function utilizing request context to fetch the
// logged in users' record. This utilizes the "user" context key
func (app *application) currentUser(r *http.Request) models.User {
	currUser, ok := r.Context().Value(contextUserKey).(models.User)
	if !ok {
		return models.User{}
	}

	return currUser

}

// newTemplateData generates a struct of templateData. This should be used in all renders.
// This function also instantiates the session flash, IsAuthenticated, and CSRFToken properties to template data.
func (app *application) newTemplateData(r *http.Request) templateData {
	var flash FlashMessage
	val := app.sessionManager.Pop(r.Context(), "flash")
	if val != nil {
		flash = val.(FlashMessage)
	}

	return templateData{
		Flash:           flash,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}
