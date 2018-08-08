package api

import (
	"net/http"

	"github.com/jmaeso/books-library/books-library"
	"github.com/jmaeso/books-library/books-library/storage"
	"github.com/yosssi/ace"
)

type Page struct {
	Books  *[]library.Book
	Filter string
	User   string
}

func GetRootHandler(bs storage.BooksStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		template, err := ace.Load("templates/index", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		p := Page{
			Books:  &[]library.Book{},
			Filter: getStringFromSession(r, "Filter"),
			User:   getStringFromSession(r, "User"),
		}

		p.Books, err = bs.GetAllSortedAndFilteredForUser(getStringFromSession(r, "SortBy"), p.Filter, p.User)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := template.Execute(w, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
