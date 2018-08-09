package api

import (
	"encoding/json"
	"net/http"

	"github.com/jmaeso/books-library/pkg/classify"
)

func PostSearchHandler(classifyService classify.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := classifyService.SearchByTitle(r.FormValue("search"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
