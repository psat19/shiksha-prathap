package main

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
	forms "github.com/psat/shiksha-prathap/pkg/forms"
	models "github.com/psat/shiksha-prathap/pkg/models"
)

type templateData struct {
	AuthenticatedUser *models.User
	CSRFToken         string
	Flash             string
	Form              *forms.Form
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		layouts, _ := filepath.Glob(filepath.Join(dir, "*.layout.tmpl"))
		if len(layouts) > 0 {
			ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
			if err != nil {
				return nil, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}

func (app application) render(w http.ResponseWriter, r http.Request, name string, data *templateData) {
	ts, ok := app.TemplateCache[name]
	if !ok {
		log.Printf("the template %s does not exist", name)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(data, &r, w))
	if err != nil {
		log.Printf("error executing the template set")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request, w http.ResponseWriter) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CSRFToken = nosurf.Token(r)
	td.Flash = app.getFlashMessages(r, w)

	return td
}

func (app *application) addFlashMessage(r *http.Request, w http.ResponseWriter, msg string) {
	session, _ := app.Session.Get(r, "session-name")
	session.AddFlash(msg)
	err := session.Save(r, w)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) getFlashMessages(r *http.Request, w http.ResponseWriter) string {
	session, _ := app.Session.Get(r, "session-name")
	flash := session.Flashes()
	if len(flash) > 0 {
		err := session.Save(r, w)
		if err != nil {
			app.serverError(w, err)
			return ""
		}
		return flash[0].(string)
	}
	return ""
}
