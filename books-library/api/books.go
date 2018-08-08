package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/goincremental/negroni-sessions"
	"github.com/gorilla/mux"
	"github.com/jmaeso/books-library/books-library"
	"github.com/jmaeso/books-library/books-library/storage"
	"github.com/jmaeso/books-library/pkg/classify"
)

func GetFilteredBooksHandler(bs storage.BooksStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := bs.GetAllSortedAndFilteredForUser(getStringFromSession(r, "SortBy"), r.FormValue("filter"), getStringFromSession(r, "User"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessions.GetSession(r).Set("Filter", r.FormValue("filter"))

		if err := json.NewEncoder(w).Encode(books); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetSortedBooksHandler(bs storage.BooksStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := bs.GetAllSortedAndFilteredForUser(r.FormValue("sortBy"), getStringFromSession(r, "Filter"), getStringFromSession(r, "User"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessions.GetSession(r).Set("SortBy", r.FormValue("sortBy"))

		if err := json.NewEncoder(w).Encode(books); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CreateBooksHandler(bs storage.BooksStore, classifyService classify.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bookResponse, err := classifyService.FindByID(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b := library.Book{
			PK:             -1,
			Title:          bookResponse.BookData.Title,
			Author:         bookResponse.BookData.Author,
			Classification: bookResponse.Classification.MostPopular,
			ID:             r.FormValue("id"),
			Username:       getStringFromSession(r, "User"),
		}

		if err := bs.Insert(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func DeleteBooksHandler(bs storage.BooksStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pk, err := strconv.ParseInt(mux.Vars(r)["pk"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := bs.DeleteForUser(pk, getStringFromSession(r, "User")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
