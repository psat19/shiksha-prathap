package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func (app application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("POST /signup", app.signupUser)
	mux.HandleFunc("/dashboard", app.showDashboard)
	mux.HandleFunc("POST /dashboard", app.updateUser)
	mux.HandleFunc("/login", app.showLogin)
	mux.HandleFunc("POST /login", app.loginUser)
	mux.HandleFunc("POST /logout", app.logoutUser)

	cs := nosurf.New(mux)

	cs.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid CSRF token", http.StatusBadRequest)
	}))

	cs.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return cs
}
