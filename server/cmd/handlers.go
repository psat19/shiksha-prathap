package main

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"

	"github.com/psat/shiksha-prathap/pkg/forms"
	models "github.com/psat/shiksha-prathap/pkg/models"
)

type userStore interface {
	Insert(string, string) (int, error)
}

func (app application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.InfoLog.Printf("Invalid Page: %s", r.URL.Path)
		app.notFound(w)
		return
	}

	session, err := app.Session.Get(r, "session-name")
	if err != nil {
		app.serverError(w, err)
	}

	if session.IsNew {
		app.render(w, *r, "signup.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	} else {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func processSignup(post url.Values, userStore userStore) (*forms.Form, int, error) {
	form := forms.New(post)
	form.Required("email", "password")
	form.MatchesPattern("email", forms.EmailRX)

	if !form.Valid() {
		return form, 0, nil
	}

	newID, err := userStore.Insert(form.Get("email"), form.Get("password"))
	if err != nil {
		return form, 0, err
	}

	return form, newID, nil
}

func (app application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form, newId, err := processSignup(r.PostForm, app.User)

	if !form.Valid() {
		app.render(w, *r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	app.InfoLog.Printf("newId is %d \n", newId)

	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Email already in use")
		app.render(w, *r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	session, _ := app.Session.Get(r, "session-name")
	session.Values["userID"] = newId
	session.Save(r, w)

	app.addFlashMessage(r, w, "Your signup was successful.")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app application) showDashboard(w http.ResponseWriter, r *http.Request) {
	session, err := app.Session.Get(r, "session-name")
	if err != nil {
		app.serverError(w, err)
	}

	if session.IsNew {

		app.render(w, *r, "login.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	}

	if userID, ok := session.Values["userID"].(int); ok {
		user := &models.User{}

		user, err := app.User.Get(userID)
		if err != nil {
			app.render(w, *r, "login.page.tmpl", &templateData{
				Form: forms.New(nil),
			})
		} else {
			uid := strconv.Itoa(user.ID)

			formData := url.Values{}
			formData.Add("email", user.Email)
			formData.Add("username", user.Name.String)
			formData.Add("phone", user.Phone.String)
			formData.Add("id", uid)

			if age := fmt.Sprintf("%d", user.Age.Int32); age != "" {
				formData.Add("age", age)
			} else {
				formData.Add("age", "")
			}

			app.render(w, *r, "dashboard.page.tmpl", &templateData{
				Form: forms.New(formData),
			})
		}

	}
}

func (app application) updateUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	userID, _ := strconv.Atoi(form.Get("id"))

	fmt.Println(r.PostForm)

	if form.Has("username") {
		form.Required("username")
		form.MaxLength("username", 255)
	}

	if form.Has("phone") {
		form.MatchesPattern("phone", forms.PhoneRX)
	}

	age := form.IsAgeValid()

	if !form.Valid() {
		app.render(w, *r, "dashboard.page.tmpl", &templateData{Form: form})
		return
	}

	app.InfoLog.Printf("Updating user %d: Name=%s, Phone=%s, Age=%d", userID, form.Get("username"), form.Get("phone"), age)

	err = app.User.Update(form.Get("username"), form.Get("phone"), int32(age), userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.addFlashMessage(r, w, "User information updated successfully.")
	app.render(w, *r, "dashboard.page.tmpl", &templateData{Form: form})
}

func (app application) showLogin(w http.ResponseWriter, r *http.Request) {
	app.render(w, *r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)

	userID, err := app.User.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, *r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	session, _ := app.Session.Get(r, "session-name")
	session.Values["userID"] = userID
	session.Save(r, w)

	app.addFlashMessage(r, w, "Logged in successfully.")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	session, _ := app.Session.Get(r, "session-name")
	delete(session.Values, "userID")
	session.Save(r, w)

	app.addFlashMessage(r, w, "You've been logged out successfully!")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
