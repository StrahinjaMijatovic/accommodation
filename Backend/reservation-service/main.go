package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	// Initialize Neo4j database connection
	if err := InitNeo4j(); err != nil {
		log.Fatal("Could not connect to Neo4j:", err)
	}
	defer CloseNeo4j()

	router := mux.NewRouter()

	session := Neo4jDriver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	router.HandleFunc("/reservations", CreateReservationHandler(session)).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"}) // Dozvoli sve origin-e; za specifične origin-e, koristi {"http://localhost:4200"}

	// Start the server
	log.Printf("Server is running on port %s\n", port)
	//log.Fatal(http.ListenAndServe(":"+port, router))
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)))
}
