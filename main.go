package main

import (
	"log"
	"fmt"
	"net/http"

	// "os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello Worwwwld!")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// TIH_API_KEY := os.Getenv("TIH_API_KEY")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte(TIH_API_KEY))
		w.Write([]byte("something"))
	})

	http.ListenAndServe(":3333", r)
}
