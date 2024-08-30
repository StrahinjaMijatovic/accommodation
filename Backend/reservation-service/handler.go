// handler.go
package main

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//func CreateReservationHandler(session *gocql.Session) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var res Reservation
//		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
//			fmt.Println("Failed to decode reservation data:", err)
//			http.Error(w, "Invalid input", http.StatusBadRequest)
//			return
//		}
//
//		fmt.Printf("Received Reservation: %+v\n", res)
//
//		// Provera dostupnosti
//		query := `SELECT id FROM reservations WHERE accommodation_id = ? AND (start_date <= ? AND end_date >= ?) ALLOW FILTERING`
//		iter := session.Query(query, res.AccommodationID, res.EndDate, res.StartDate).Iter()
//
//		if iter.NumRows() > 0 {
//			http.Error(w, "Accommodation is not available for the selected dates", http.StatusConflict)
//			return
//		}
//
//		// Kreiranje rezervacije
//		res.ID = gocql.TimeUUID()
//		query = `INSERT INTO reservations (id, accommodation_id, guest_id, start_date, end_date) VALUES (?, ?, ?, ?, ?)`
//		if err := session.Query(query, res.ID, res.AccommodationID, res.GuestID, res.StartDate, res.EndDate).Exec(); err != nil {
//			http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
//			return
//		}
//
//		w.WriteHeader(http.StatusCreated)
//		json.NewEncoder(w).Encode(res)
//	}
//}

func CreateReservationHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res Reservation
		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			log.Printf("Failed to decode reservation data: %v\n", err)
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Provera da li su datumi već rezervisani
		query := `SELECT COUNT(*) FROM reservations WHERE accommodation_id = ? AND start_date <= ? AND end_date >= ? ALLOW FILTERING`
		var count int
		err := session.Query(query, res.AccommodationID, res.EndDate, res.StartDate).Scan(&count)
		if err != nil {
			log.Printf("Database error during availability check: %v\n", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if count > 0 {
			log.Printf("Dates are already reserved for accommodation ID: %v\n", res.AccommodationID)
			http.Error(w, "The selected dates are already reserved", http.StatusConflict)
			return
		}

		// Ako su datumi slobodni, kreirajte rezervaciju
		res.ID = gocql.TimeUUID()
		query = `INSERT INTO reservations (id, accommodation_id, guest_id, start_date, end_date) VALUES (?, ?, ?, ?, ?)`
		err = session.Query(query, res.ID, res.AccommodationID, res.GuestID, res.StartDate, res.EndDate).Exec()
		if err != nil {
			log.Printf("Failed to insert reservation into the database: %v\n", err)
			http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
			return
		}

		log.Printf("Reservation created successfully: %+v\n", res)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	}
}

func CancelReservationHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		reservationIDStr := vars["reservationID"]
		log.Printf("Received request to cancel reservation with ID: '%s'", reservationIDStr)

		if reservationIDStr == "" {
			log.Printf("Empty reservation ID received")
			http.Error(w, "Invalid reservation ID format", http.StatusBadRequest)
			return
		}

		reservationID, err := gocql.ParseUUID(reservationIDStr)
		if err != nil {
			log.Printf("Invalid reservation ID format: %v", err)
			http.Error(w, "Invalid reservation ID format", http.StatusBadRequest)
			return
		}

		var count int
		log.Printf("Checking if reservation exists with ID: %s", reservationID)
		queryCheck := `SELECT COUNT(*) FROM reservations WHERE id = ?`
		if err := session.Query(queryCheck, reservationID).Scan(&count); err != nil || count == 0 {
			log.Printf("Reservation not found or error checking reservation: %v", err)
			http.Error(w, "Reservation not found", http.StatusNotFound)
			return
		}

		log.Printf("Reservation found with ID: %s, proceeding to delete", reservationID)
		queryDelete := `DELETE FROM reservations WHERE id = ?`
		if err := session.Query(queryDelete, reservationID).Exec(); err != nil {
			log.Printf("Failed to cancel reservation with ID %s: %v", reservationID, err)
			http.Error(w, "Failed to cancel reservation", http.StatusInternalServerError)
			return
		}

		log.Printf("Reservation with ID %s successfully canceled", reservationID)
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetReservationsByUserHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"]

		var reservations []Reservation
		query := `SELECT id, accommodation_id, guest_id, start_date, end_date FROM reservations WHERE guest_id = ? ALLOW FILTERING`
		iter := session.Query(query, userID).Iter()

		var reservation Reservation
		for iter.Scan(&reservation.ID, &reservation.AccommodationID, &reservation.GuestID, &reservation.StartDate, &reservation.EndDate) {
			reservations = append(reservations, reservation)
		}

		if err := iter.Close(); err != nil {
			log.Printf("Failed to fetch reservations: %v", err)
			http.Error(w, "Failed to fetch reservations", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(reservations); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func HasActiveReservationsHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"]

		// Proverava da li korisnik ima rezervacije u bazi
		query := `SELECT COUNT(*) FROM reservations WHERE guest_id = ? ALLOW FILTERING`
		var count int
		if err := session.Query(query, userID).Scan(&count); err != nil {
			http.Error(w, "Failed to check reservations", http.StatusInternalServerError)
			return
		}

		hasReservations := count > 0
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(hasReservations)
	}
}
