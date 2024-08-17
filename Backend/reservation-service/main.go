package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Inicijalizacija Neo4j baze podataka
	err := InitNeo4j()
	if err != nil {
		log.Fatal("Could not connect to Neo4j", err)
	}
	defer CloseNeo4j()

	// Postavljanje ruta za servis rezervacija
	http.HandleFunc("/reservations", CreateReservation)
	http.HandleFunc("/reservations/cancel", CancelReservation)

	// Dobijanje porta iz okruženja ili default na 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Pokretanje servera
	log.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
