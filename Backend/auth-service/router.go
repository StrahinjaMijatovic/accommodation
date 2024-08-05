package main

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	return r
}
