package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func gotoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		fmt.Println(r.URL.Path)
		urls := readDB()
		isValidID := URL{ID: "", OriginalURL: ""}
		for i, url := range urls {
			_ = i
			if url.ID == r.URL.Path[1:] {
				isValidID = url
			}
		}
		if isValidID.ID == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "invalid short code"}`))
			return
		} else {
			w.WriteHeader(http.StatusFound)
			http.Redirect(w, r, isValidID.OriginalURL, http.StatusSeeOther)
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		defer r.Body.Close()
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		var url URL
		err = json.Unmarshal(reqBody, &url)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		code := generateRandomCode()
		url.ID = code
		urls := readDB()
		urls = append(urls, url)

		dataBytes, err := json.Marshal(urls)
		if err != nil {
			fmt.Println("Error marshalling data", err)
			return
		}

		err = os.WriteFile("db.json", dataBytes, 0o644)
		if err != nil {
			fmt.Println("Error writing to file", err)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(url)
	case "GET":
		urls := readDB()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(urls)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

func handleRequests() {
	http.HandleFunc("/urls", urlHandler)
	http.HandleFunc("/{id}", gotoHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
