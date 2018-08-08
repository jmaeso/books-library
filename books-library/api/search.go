package api

import (
	"encoding/json"
	"net/http"
)

func PostSearchHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := search(r.FormValue("search"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
