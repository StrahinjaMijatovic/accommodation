package main

import (
	"github.com/gocql/gocql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func connectToCassandra() *gocql.Session {
	var session *gocql.Session
	var err error

	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "accommodations"
	cluster.Consistency = gocql.Quorum

	for i := 0; i < 10; i++ {
		session, err = cluster.CreateSession()
		if err == nil {
			log.Println("Connected to Cassandra")
			return session
		}

		log.Printf("Failed to connect to Cassandra (attempt %d): %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}

	log.Fatalf("Failed to connect to Cassandra after multiple attempts: %v", err)
	return nil
}

func main() {
	session := connectToCassandra()
	defer session.Close()

	router := mux.NewRouter()

	router.HandleFunc("/hosts/{id}/rate", RateHostHandler(session)).Methods("POST", "PUT", "DELETE")
	router.HandleFunc("/hosts/{id}/ratings", GetHostRatingsHandler(session)).Methods("GET")
	router.HandleFunc("/accommodations/{id}/rate", RateAccommodationHandler(session)).Methods("POST", "PUT", "DELETE")
	router.HandleFunc("/accommodations/{id}/ratings", GetAccommodationRatingsHandler(session)).Methods("GET")
	router.HandleFunc("/hosts/{id}/ratings", DeleteHostRatingHandler(session)).Methods("DELETE")
	router.HandleFunc("/accommodations/{id}/ratings", DeleteAccommodationRatingHandler(session)).Methods("DELETE")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Println("Server started on :8082")
	log.Fatal(http.ListenAndServe(":8082", handlers.CORS(headers, methods, origins)(router)))
}
