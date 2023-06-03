package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

var TIH_URL = "https://api.stb.gov.sg/content/food-beverages/v2/search?searchType=keyword"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	TIH_API_KEY := os.Getenv("TIH_API_KEY")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		searchValue := r.URL.Query().Get("searchValue")
		client := &http.Client{}

		reqURL := fmt.Sprintf("%s%s", TIH_URL, searchValue)
		req, err := http.NewRequest("GET", reqURL, nil)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}

		req.Header.Add("X-API-Key", TIH_API_KEY)
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Error making request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response", http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})

	log.Fatal(http.ListenAndServe(":3333", r))
}
