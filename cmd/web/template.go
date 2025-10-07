// template.go is a file containing any functions related to template data and their files or directories.
package main

import (
	"ezkitchen/internal/models"
	"path/filepath"
	"text/template"
)

// templateData contains keys and types for models of Estimate and Customers(users)
// in the future when multiple tables are required to load an estimate
// (ie. Surveyor(user), Estimate, and Customer(user)) be sure to update this.
type templateData struct {
	Estimate models.Estimate
	Customer models.User
	Product  models.Product
	Products []models.Product
}

// newTemplateCache is a function that runs on server start. This function parses any pages/partial/modals
// templates used in our server to prevent repetitive code and frequent file parsing on each page load.
func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/pages/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts

	}

	modals, err := filepath.Glob("./ui/html/modals/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, modal := range modals {
		name := filepath.Base(modal)

		ts, err := template.ParseFiles(modal)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
