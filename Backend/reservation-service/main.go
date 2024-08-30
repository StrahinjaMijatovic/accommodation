// main.go
package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	session := InitCassandra()
	defer session.Close()

	router := mux.NewRouter()
	router.HandleFunc("/reservations", CreateReservationHandler(session)).Methods("POST")
	//router.HandleFunc("/reservations/{id}", CancelReservationHandler(session)).Methods("DELETE")
	router.HandleFunc("/guests/{userID}/reservations", GetReservationsByUserHandler(session)).Methods("GET")
	router.HandleFunc("/reservations/{reservationID}", CancelReservationHandler(session)).Methods("DELETE")
	router.HandleFunc("/reservations/active/{userID}", HasActiveReservationsHandler(session)).Methods("GET")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"}) // Dozvoli sve origin-e; za specifične origin-e, koristi {"http://localhost:4200"}

	log.Println("Reservation service running on :8081")
	//log.Fatal(http.ListenAndServe(":8081", router))
	log.Fatal(http.ListenAndServe(":8081", handlers.CORS(headers, methods, origins)(router)))
}
