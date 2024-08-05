package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	ConnectDatabase() // Povezivanje sa MongoDB bazom

	r := mux.NewRouter()
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
