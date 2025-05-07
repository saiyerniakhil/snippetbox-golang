package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.saiyerniakhil.in/internal/models"
	"snippetbox.saiyerniakhil.in/internal/validator"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	buf.WriteTo(w)

	w.WriteHeader(status)

}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.SnippetList = snippets

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)

	// fmt.Fprintf(w, "Hello from Snippetbox!")
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)

	data.Snippet = s

	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// parse the form
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	formData := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	formData.CheckField(validator.NotBlank(formData.Title), "title", "This field cannot be blank")
	formData.CheckField(validator.NotBlank(formData.Content), "content", "This field cannot be blank")
	formData.CheckField(validator.MaxChars(formData.Title, 100), "title", "This field cannot more than 100 chars")
	formData.CheckField(validator.PermitteddValue(formData.Expires, 1, 7, 365), "content", "This field cannot anything other than 1, 7 and 365")

	if len(formData.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = formData
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
	}

	id, err := app.snippets.Insert(formData.Title, formData.Content, formData.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
