package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func main() {
	// Povezivanje sa MongoDB bazom
	ConnectDatabase()

	// Inicijalizacija router-a
	r := mux.NewRouter()
	r.HandleFunc("/profile", GetProfileHandler).Methods("GET")
	r.HandleFunc("/profile", UpdateProfileHandler).Methods("PUT")

	// Pokretanje HTTP servera
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
func ConnectDatabase() {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database("authdb") // Koristite tačno ime baze
}
