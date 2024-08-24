package main

import (
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
)

func CreateReservationHandler(session neo4j.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reservation Reservation

		if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Here you'd typically check if the accommodation is available and not already booked
		// Example: checkAvailabilityAndBook(reservation) - implement this based on your logic

		query := `CREATE (r:Reservation {id: $id, accommodation_id: $accommodation_id, guest_id: $guest_id, start_date: $start_date, end_date: $end_date, status: $status})`
		params := map[string]interface{}{
			"id":               reservation.ID,
			"accommodation_id": reservation.AccommodationID,
			"guest_id":         reservation.GuestID,
			"start_date":       reservation.StartDate.Format("2006-01-02"),
			"end_date":         reservation.EndDate.Format("2006-01-02"),
			"status":           reservation.Status,
		}

		_, err := session.Run(query, params)
		if err != nil {
			log.Printf("Failed to create reservation: %v", err)
			http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(reservation)
	}
}

func CancelReservation(w http.ResponseWriter, r *http.Request) {
	var reservationID struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reservationID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session := Neo4jDriver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	// Proveravanje da li rezervacija postoji i može se otkazati
	result, err := session.Run(
		`MATCH (r:Reservation {id: $id, status: 'active'})
         WHERE r.start_date > date()
         RETURN r`,
		map[string]interface{}{
			"id": reservationID.ID,
		},
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !result.Next() {
		http.Error(w, "Reservation not found or cannot be cancelled", http.StatusNotFound)
		return
	}

	// Otkazivanje rezervacije (status postaje 'cancelled')
	_, err = session.Run(
		`MATCH (r:Reservation {id: $id})
         SET r.status = 'cancelled'
         RETURN r`,
		map[string]interface{}{
			"id": reservationID.ID,
		},
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reservation cancelled successfully"))
}
