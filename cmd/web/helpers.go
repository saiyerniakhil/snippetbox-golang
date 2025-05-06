package main

import (
	"net/http"
	"runtime/debug"
	"strings"
	"time"
	"unicode/utf8"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	app.logger.Error("Error", http.StatusText(status), status)
	http.Error(w, http.StatusText(status), status)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (formData *snippetCreateForm) validateForm() {
	fieldErrors := make(map[string]string)
	//validate title
	if strings.TrimSpace(formData.Title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(formData.Title) > 100 {
		fieldErrors["title"] = "This field cannot be more more than 100 characters long"
	}

	// validate content
	if strings.TrimSpace(formData.Content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	//validate expires field
	if formData.Expires != 7 && formData.Expires != 1 && formData.Expires != 365 {
		fieldErrors["expires"] = "This field must be equal to either of 1, 7 or 365"
	}

	formData.FieldErrors = fieldErrors

}
