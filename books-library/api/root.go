package api

import (
	"net/http"

	"github.com/jmaeso/books-library/books-library"
	"github.com/yosssi/ace"
	"gopkg.in/gorp.v2"
)

type Page struct {
	Books  []library.Book
	Filter string
	User   string
}

func GetRootHandler(dbmap *gorp.DbMap) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		template, err := ace.Load("templates/index", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		p := Page{
			Books:  []library.Book{},
			Filter: getStringFromSession(r, "Filter"),
			User:   getStringFromSession(r, "User"),
		}

		if !getBookCollection(dbmap, &p.Books, getStringFromSession(r, "SortBy"), getStringFromSession(r, "Filter"),
			p.User, w) {
			return
		}

		if err := template.Execute(w, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
