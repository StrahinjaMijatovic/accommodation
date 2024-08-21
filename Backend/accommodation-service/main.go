package main

import (
	"github.com/gocql/gocql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// Retry logic function to connect to Cassandra
func connectToCassandra() *gocql.Session {
	var session *gocql.Session
	var err error

	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "accommodations"
	cluster.Consistency = gocql.Quorum

	for i := 0; i < 10; i++ { // Retry 10 times
		session, err = cluster.CreateSession()
		if err == nil {
			log.Println("Connected to Cassandra")
			return session
		}

		log.Printf("Failed to connect to Cassandra (attempt %d): %v\n", i+1, err)
		time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
	}

	log.Fatalf("Failed to connect to Cassandra after multiple attempts: %v", err)
	return nil
}

func main() {
	// Inicijalizacija konekcije ka Cassandra bazi sa retry logikom
	session := connectToCassandra()
	defer session.Close()

	// Postavljanje routera
	router := mux.NewRouter()

	// Definisanje ruta
	router.HandleFunc("/accommodations", CreateAccommodationHandler(session)).Methods("POST")
	router.HandleFunc("/accommodations/{id}", GetAccommodationHandler(session)).Methods("GET")
	router.HandleFunc("/accommodations/{id}", UpdateAccommodationHandler(session)).Methods("PUT")
	router.HandleFunc("/accommodations/{id}", DeleteAccommodationHandler(session)).Methods("DELETE")
	router.HandleFunc("/accommodations", GetAllAccommodationsHandler(session)).Methods("GET")
	router.HandleFunc("/search", SearchAccommodationsHandler(session)).Methods("GET")
	router.HandleFunc("/accommodations/{id}", GetAccommodationByID(session)).Methods("GET")
	router.HandleFunc("/accommodations/{id}/availability", UpdateAvailabilityHandler(session)).Methods("POST", "PUT")
	router.HandleFunc("/accommodations/{id}/price", UpdatePriceHandler(session)).Methods("POST", "PUT")
	router.HandleFunc("/accommodations/{id}/availability-and-price", UpdateAvailabilityAndPriceHandler(session)).Methods("PUT")
	router.HandleFunc("/accommodations/{id}/availability", GetAvailabilityHandler(session)).Methods("GET")
	router.HandleFunc("/accommodations/{id}/availability", GetAvailabilityByAccommodationIDHandler(session)).Methods("GET")
	// Omogućavanje CORS-a za sve rute
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"}) // Dozvoli sve origin-e; za specifične origin-e, koristi {"http://localhost:4200"}

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(router)))
}
