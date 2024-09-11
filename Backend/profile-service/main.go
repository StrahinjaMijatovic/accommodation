package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func main() {
	ConnectDatabase()

	r := mux.NewRouter()
	r.HandleFunc("/profile", GetProfileHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/profile", UpdateProfileHandler).Methods("PUT", "OPTIONS")
	r.HandleFunc("/change-password", ChangePasswordHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/profile/{userID}", DeleteProfileHandler).Methods("DELETE", "OPTIONS")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(r)

	http.Handle("/", corsHandler)
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

	db = client.Database("authdb")
}
