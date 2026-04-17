package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"ize-302/url-shortener/utils"
)

type URL struct {
	ID   int    `json:"id"`
	URL  string `json:"url"`
	Code string `json:"code"`
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
			err = rows.Scan(&id, &url, &code)
			if err != nil {
				log.Fatal(err)
			}
			urlObj := URL{ID: id, Code: code, URL: url}
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
		_, err = db.Exec(`INSERT into urls(url, code) VALUES(?, ?)`, url.URL, code)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("New url successfully created")
		w.WriteHeader(http.StatusCreated)
		urlData, err := db.Exec(`SELECT * FROM urls WHERE code = ?`, code)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(urlData)
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
			err = rows.Scan(&id, &url, &code)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, url, code)
			urlObj := URL{ID: id, Code: code, URL: url}
			urls = append(urls, urlObj)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(urls)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}
