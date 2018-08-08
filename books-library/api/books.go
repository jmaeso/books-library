package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/goincremental/negroni-sessions"
	"github.com/gorilla/mux"
	"github.com/jmaeso/books-library/books-library"
	"github.com/jmaeso/books-library/pkg/classify"
	"gopkg.in/gorp.v2"
)

func GetFilteredBooksHandler(dbmap *gorp.DbMap) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var b []library.Book
		if !getBookCollection(dbmap, &b, getStringFromSession(r, "SortBy"), r.FormValue("filter"),
			getStringFromSession(r, "User"), w) {
			return
		}

		sessions.GetSession(r).Set("Filter", r.FormValue("filter"))

		if err := json.NewEncoder(w).Encode(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetSortedBooksHandler(dbmap *gorp.DbMap) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var b []library.Book
		if !getBookCollection(dbmap, &b, r.FormValue("sortBy"), getStringFromSession(r, "Filter"),
			getStringFromSession(r, "User"), w) {
			return
		}

		sessions.GetSession(r).Set("SortBy", r.FormValue("sortBy"))

		if err := json.NewEncoder(w).Encode(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CreateBooksHandler(dbmap *gorp.DbMap, classifyService classify.Service) func(w http.ResponseWriter, r *http.Request) {
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

		if err := dbmap.Insert(&b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func DeleteBooksHandler(dbmap *gorp.DbMap) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pk, err := strconv.ParseInt(mux.Vars(r)["pk"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var b library.Book
		if err := dbmap.SelectOne(&b, "select * from books where pk=? and user=?", pk, getStringFromSession(r, "User")); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := dbmap.Delete(&b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// TODO: All the underlying code needs to be refactored and restructured.

func getBookCollection(dbmap *gorp.DbMap, books *[]library.Book, sortCol, filterByClass, username string, w http.ResponseWriter) bool {
	if sortCol == "" {
		sortCol = "pk"
	}

	where := "where user = ?"
	if filterByClass == "fiction" {
		where += "and classification between '800' and '900'"
	} else if filterByClass == "nonfiction" {
		where += "and classification not between '800' and '900'"
	}

	if _, err := dbmap.Select(books, "select * from books "+where+" order by "+sortCol, username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
}
