// helper.go contains our error handling to distinguish between our client side or server side errors.
// this also contains any of our data type parsers to convert our html form data.

package main

import (
	"log"
	"net/http"
	"strconv"
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
