package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	ConnectDatabase() // Povezivanje sa MongoDB bazom

	r := mux.NewRouter()
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/update-profile", UpdateProfileHandler).Methods("PUT")

	// Apply CORS middleware
	handler := cors.Default().Handler(r)

	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
