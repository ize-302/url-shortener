// Package route
package route

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"ize-302/url-shortener/model"
	"ize-302/url-shortener/util"
)

func RegisterHandlers(store *util.Store) {
	http.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		row, err := store.FetchURLByCode(r.PathValue("id"))
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"message": "invalid short code"}`, http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, row.URL, http.StatusSeeOther)
	})

	http.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		defer r.Body.Close()

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var url model.URL
		err = json.Unmarshal(reqBody, &url)
		if err != nil {
			http.Error(w, "internal server error", http.StatusBadRequest)
			return
		}

		result, err := store.SaveURL(&url)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		row, err := store.FetchURLByID(int(lastInsertID))
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"message": "invalid id"}`, http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(row)
	})

	http.HandleFunc("GET /urls", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		urls, err := store.FetchURLs()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(urls)
	})
}
