// template.go is a file containing any functions related to template data and their files or directories.
package main

import (
	"ezkitchen/internal/models"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

// templateData contains keys and types for models of Estimate, users, products, estimate totals, and forms
// in the future when multiple tables are required to load an estimate
// (ie. Surveyor(user), Estimate, and Customer(user)) be sure to update this.
type templateData struct {
	Estimate       models.Estimate
	Customer       models.User
	Products       []models.EstimateProduct
	EstimateTotals models.EstimateTotals
	Form           any
	Flash          FlashMessage
}

type FlashMessage struct {
	Type    string
	Message string
}

// newTemplateCache is a function that runs on server start. This function parses any pages/partial/modals
// templates as well as any functions used within tmpl to prevent repetitive code and frequent file parsing on each page load.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	funcs := template.FuncMap{
		"centsToDollars": func(centValue int, quantity int) string {
			return fmt.Sprintf("%.2f", float64(centValue*quantity)/100)
		},
		"list": func(vals ...string) []string {
			return vals
		},
	}

	pages, err := filepath.Glob("./ui/html/pages/**/*.tmpl")
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir("./ui/html/pages", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".tmpl") {
			pages = append(pages, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	layouts, err := filepath.Glob("./ui/html/layouts/*.tmpl")
	if err != nil {
		return nil, err
	}

	partials, err := filepath.Glob("./ui/html/partials/**/*.tmpl")
	if err != nil {
		return nil, err
	}

	modals, err := filepath.Glob("./ui/html/modals/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts := template.New(name).Funcs(funcs)

		ts, err = ts.ParseFiles(layouts...)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(partials...)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(modals...)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
