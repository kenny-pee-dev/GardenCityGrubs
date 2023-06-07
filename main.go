package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

var TIH_URL = "https://api.stb.gov.sg/content/food-beverages/v2/search?searchType=keyword&limit=50"

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

	r.Get("/.well-known/ai-plugin.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ai-plugin.json")
	})

	r.Get("/open-api.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./open-api.yaml")
	})

	r.Get("/foodPlaces", func(w http.ResponseWriter, r *http.Request) {
		searchValue := r.URL.Query().Get("searchValues")
		client := &http.Client{}

		reqURL, err := url.Parse(TIH_URL)
		if err != nil {
			log.Fatal(err)
		}
		values := reqURL.Query()
		values.Add("searchValues", searchValue)
		reqURL.RawQuery = values.Encode()

		req, err := http.NewRequest("GET", reqURL.String(), nil)
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
