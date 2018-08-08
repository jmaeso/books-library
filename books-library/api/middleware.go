package api

import (
	"database/sql"
	"net/http"

	"github.com/jmaeso/books-library/books-library"
	"gopkg.in/gorp.v2"
)

func VerifyUser(dbmap *gorp.DbMap) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.URL.Path == "/login" {
			next(w, r)
			return
		}

		if username := getStringFromSession(r, "User"); username != "" {
			if user, _ := dbmap.Get(library.User{}, username); user != nil {
				next(w, r)
				return
			}
		}

		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}
}

func VerifyDB(db *sql.DB) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if err := db.Ping(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next(w, r)
	}
}
