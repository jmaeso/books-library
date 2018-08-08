package api

import (
	"net/http"

	"github.com/goincremental/negroni-sessions"
	"github.com/jmaeso/books-library/books-library"
	"github.com/yosssi/ace"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gorp.v2"
)

type LoginPage struct {
	Error string
}

func LoginHandler(dbmap *gorp.DbMap) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var p LoginPage
		if r.FormValue("register") != "" {
			secret, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			user := library.User{Username: r.FormValue("username"), Secret: secret}
			if err := dbmap.Insert(&user); err != nil {
				p.Error = err.Error()
			} else {
				sessions.GetSession(r).Set("User", user.Username)

				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		} else if r.FormValue("login") != "" {
			user, err := dbmap.Get(library.User{}, r.FormValue("username"))
			if err != nil {
				p.Error = err.Error()
			} else if user == nil {
				p.Error = "No such user with Username" + r.FormValue("username")
			} else {
				u := user.(*library.User)
				if err := bcrypt.CompareHashAndPassword(u.Secret, []byte(r.FormValue("password"))); err != nil {
					p.Error = err.Error()
				} else {
					sessions.GetSession(r).Set("User", u.Username)

					http.Redirect(w, r, "/", http.StatusFound)
					return
				}
			}
		}

		template, err := ace.Load("templates/login", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = template.Execute(w, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
