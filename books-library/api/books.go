package api

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/goincremental/negroni-sessions"
	"github.com/gorilla/mux"
	"github.com/jmaeso/books-library/books-library"
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

func CreateBooksHandler(dbmap *gorp.DbMap) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		book, err := find(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b := library.Book{
			PK:             -1,
			Title:          book.BookData.Title,
			Author:         book.BookData.Author,
			Classification: book.Classification.MostPopular,
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

type ClassifyBookResponse struct {
	BookData struct {
		Title  string `xml:"title,attr"`
		Author string `xml:"author,attr"`
		ID     string `xml:"owi,attr"`
	} `xml:"work"`
	Classification struct {
		MostPopular string `xml:"sfa,attr"`
	} `xml:"recommendations>ddc>mostPopular"`
}

type SearchResult struct {
	Title  string `xml:"title,attr"`
	Author string `xml:"author,attr"`
	Year   string `xml:"hyr,attr"`
	ID     string `xml:"owi,attr"`
}

type ClassifyResponse struct {
	Results []SearchResult `xml:"works>work"`
}

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

func find(id string) (ClassifyBookResponse, error) {
	body, err := classifyAPI("http://classify.oclc.org/classify2/Classify?&summary=true&owi=" + url.QueryEscape(id))
	if err != nil {
		return ClassifyBookResponse{}, err
	}

	var c ClassifyBookResponse
	if err := xml.Unmarshal(body, &c); err != nil {
		return ClassifyBookResponse{}, err
	}

	return c, nil
}

func search(query string) ([]SearchResult, error) {
	body, err := classifyAPI("http://classify.oclc.org/classify2/Classify?&summary=true&title=" + url.QueryEscape(query))
	if err != nil {
		return []SearchResult{}, err
	}

	var c ClassifyResponse
	err = xml.Unmarshal(body, &c)

	return c.Results, err
}

func classifyAPI(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
