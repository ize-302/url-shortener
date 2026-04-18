package routes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"ize-302/url-shortener/utils"
)

type URL struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Code      string `json:"code"`
	CreatedAt string `json:"createdAt"`
}

func gotoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		db := utils.ReadDB()
		defer db.Close()

		rows, err := db.Query(`SELECT * FROM urls WHERE code = ?`, r.URL.Path[1:])
		if err != nil {
			log.Fatalln(err)
		}
		defer rows.Close()

		var urls []URL
		for rows.Next() {
			var id int
			var url string
			var code string
			var createdAt string
			err = rows.Scan(&id, &url, &code, &createdAt)
			if err != nil {
				log.Fatal(err)
			}
			urlObj := URL{ID: id, Code: code, URL: url, CreatedAt: createdAt}
			urls = append(urls, urlObj)

		}
		if len(urls) == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "invalid short code"}`))
		} else {
			w.WriteHeader(http.StatusFound)
			url := urls[0]
			http.Redirect(w, r, url.URL, http.StatusSeeOther)
		}

		if err = rows.Err(); err != nil {
			log.Fatal(err)
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

		db := utils.ReadDB()
		defer db.Close()

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
		code := utils.GenerateRandomCode()
		createdAt := time.Now().UTC().Format(time.RFC3339)
		_, err = db.Exec(`INSERT into urls(url, code, createdAt) VALUES(?, ?, ?)`, url.URL, code, createdAt)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusCreated)
		urlData, err := db.Query(`SELECT * FROM urls WHERE code = ?`, code)
		if err != nil {
			log.Fatal(err)
		}
		for urlData.Next() {
			var id int
			var url string
			var code string
			var createdAt string
			err = urlData.Scan(&id, &url, &code, &createdAt)
			if err != nil {
				log.Fatal(err)
			}
			urlObj := URL{ID: id, Code: code, URL: url, CreatedAt: createdAt}
			json.NewEncoder(w).Encode(urlObj)
		}
	case "GET":
		db := utils.ReadDB()
		defer db.Close()

		rows, err := db.Query("SELECT * FROM urls")
		if err != nil {
			log.Fatalln(err)
		}

		var urls []URL
		for rows.Next() {
			var id int
			var url string
			var code string
			var createdAt string
			err = rows.Scan(&id, &url, &code, &createdAt)
			if err != nil {
				log.Fatal(err)
			}
			urlObj := URL{ID: id, Code: code, URL: url, CreatedAt: createdAt}
			urls = append(urls, urlObj)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(urls)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}
