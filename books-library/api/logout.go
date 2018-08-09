package api

import (
	"net/http"

	"github.com/goincremental/negroni-sessions"
)

func LogoutHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions.GetSession(r).Set("User", nil)
		sessions.GetSession(r).Set("Filter", nil)

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
