package main

import (
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
)

func CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservation Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session := Neo4jDriver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	// Provjera da li postoji preklapanje datuma za isti smještaj
	result, err := session.Run(
		`MATCH (r:Reservation {accommodation_id: $accommodation_id})
         WHERE r.status = 'active' AND r.start_date < $end_date AND r.end_date > $start_date
         RETURN count(r) AS count`,
		map[string]interface{}{
			"accommodation_id": reservation.AccommodationID,
			"start_date":       reservation.StartDate,
			"end_date":         reservation.EndDate,
		},
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var count int64
	if result.Next() {
		count = result.Record().Values[0].(int64)
	}

	if count > 0 {
		http.Error(w, "Accommodation already booked for these dates", http.StatusConflict)
		return
	}

	// Kreiranje nove rezervacije
	_, err = session.Run(
		`CREATE (r:Reservation {accommodation_id: $accommodation_id, guest_id: $guest_id, start_date: $start_date, end_date: $end_date, status: 'active'}) RETURN r`,
		map[string]interface{}{
			"accommodation_id": reservation.AccommodationID,
			"guest_id":         reservation.GuestID,
			"start_date":       reservation.StartDate,
			"end_date":         reservation.EndDate,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reservation)
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
