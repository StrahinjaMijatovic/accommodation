package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// CreateAccommodationHandler kreira novi smeštaj
func CreateAccommodationHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var acc Accommodation
		if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		acc.ID = gocql.TimeUUID()

		query := `INSERT INTO accommodations (id, name, location, guests, price, amenities, images) VALUES (?, ?, ?, ?, ?, ?, ?)`
		if err := session.Query(query, acc.ID, acc.Name, acc.Location, acc.Guests, acc.Price, acc.Amenities, acc.Images).Exec(); err != nil {
			log.Printf("Failed to create accommodation: %v", err)
			http.Error(w, "Failed to create accommodation", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(acc)
	}
}

// GetAccommodationHandler vraća informacije o smeštaju na osnovu ID-a
//func GetAccommodationHandler(session *gocql.Session) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		vars := mux.Vars(r)
//		id, err := gocql.ParseUUID(vars["id"])
//		if err != nil {
//			http.Error(w, "Invalid ID format", http.StatusBadRequest)
//			return
//		}
//
//		var acc Accommodation
//		query := `SELECT id, name, location, guests, price FROM accommodations WHERE id = ? LIMIT 1`
//		if err := session.Query(query, id).Scan(&acc.ID, &acc.Name, &acc.Location, &acc.Guests, &acc.Price); err != nil {
//			log.Printf("Failed to retrieve accommodation: %v", err)
//			http.Error(w, "Accommodation not found", http.StatusNotFound)
//			return
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode(acc)
//	}
//}

func GetAccommodationHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := gocql.ParseUUID(vars["id"])
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		var acc Accommodation
		query := `SELECT id, name, location, guests, price, amenities, images FROM accommodations WHERE id = ? LIMIT 1`
		if err := session.Query(query, id).Scan(
			&acc.ID,
			&acc.Name,
			&acc.Location,
			&acc.Guests,
			&acc.Price,
			&acc.Amenities,
			&acc.Images); err != nil {
			log.Printf("Failed to retrieve accommodation: %v", err)
			http.Error(w, "Accommodation not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(acc)
	}
}

func GetAllAccommodationsHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accommodations []Accommodation

		query := `SELECT id, name, location, guests, price, amenities, images FROM accommodations`
		iter := session.Query(query).Iter()
		var acc Accommodation
		for iter.Scan(&acc.ID, &acc.Name, &acc.Location, &acc.Guests, &acc.Price, &acc.Amenities, &acc.Images) {
			accommodations = append(accommodations, acc)
		}
		if err := iter.Close(); err != nil {
			log.Printf("Failed to fetch accommodations: %v", err)
			http.Error(w, "Failed to fetch accommodations", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(accommodations)
	}
}

// UpdateAccommodationHandler ažurira postojeći smeštaj
func UpdateAccommodationHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := gocql.ParseUUID(vars["id"])
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		var acc Accommodation
		if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		query := `UPDATE accommodations SET name = ?, location = ?, guests = ?, price = ? WHERE id = ?`
		if err := session.Query(query, acc.Name, acc.Location, acc.Guests, acc.Price, id).Exec(); err != nil {
			log.Printf("Failed to update accommodation: %v", err)
			http.Error(w, "Failed to update accommodation", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(acc)
	}
}

// DeleteAccommodationHandler briše postojeći smeštaj
func DeleteAccommodationHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := gocql.ParseUUID(vars["id"])
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		query := `DELETE FROM accommodations WHERE id = ?`
		if err := session.Query(query, id).Exec(); err != nil {
			log.Printf("Failed to delete accommodation: %v", err)
			http.Error(w, "Failed to delete accommodation", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetAccommodationByID(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Preuzimanje ID-ja iz URL parametara
		vars := mux.Vars(r)
		id := vars["id"]

		// Konverzija ID-ja u gocql.UUID
		uuid, err := gocql.ParseUUID(id)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		// Preuzimanje smeštaja iz baze podataka
		var accommodation Accommodation
		query := "SELECT id, name, location, guests, price, amenities, images FROM accommodations WHERE id = ? LIMIT 1"
		if err := session.Query(query, uuid).Scan(
			&accommodation.ID,
			&accommodation.Name,
			&accommodation.Location,
			&accommodation.Guests,
			&accommodation.Price,
			&accommodation.Amenities,
			&accommodation.Images,
		); err != nil {
			if err == gocql.ErrNotFound {
				http.Error(w, "Accommodation not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to fetch accommodation: %v", err), http.StatusInternalServerError)
			}
			return
		}

		// Slanje odgovora u JSON formatu
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(accommodation); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		}
	}
}
func SearchAccommodationsHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		location := r.URL.Query().Get("location")
		guestsStr := r.URL.Query().Get("guests")

		var accommodations []Accommodation
		var query string
		var params []interface{}

		// Initialize query with base SELECT statement
		query = "SELECT id, name, location, guests, price, amenities, images FROM accommodations"
		allowFiltering := false

		if location != "" {
			query += " WHERE location = ?"
			params = append(params, location)
			allowFiltering = true
		}

		if guestsStr != "" {
			if allowFiltering {
				query += " AND guests >= ?"
			} else {
				query += " WHERE guests >= ?"
				allowFiltering = true
			}
			guests, err := strconv.Atoi(guestsStr)
			if err != nil {
				http.Error(w, "Invalid number of guests", http.StatusBadRequest)
				return
			}
			params = append(params, guests)
		}

		// Add ALLOW FILTERING to the query if needed
		if allowFiltering {
			query += " ALLOW FILTERING"
		}

		// Execute query
		iter := session.Query(query, params...).Iter()
		var acc Accommodation
		for iter.Scan(&acc.ID, &acc.Name, &acc.Location, &acc.Guests, &acc.Price, &acc.Amenities, &acc.Images) {
			accommodations = append(accommodations, acc)
		}
		if err := iter.Close(); err != nil {
			log.Printf("Failed to fetch accommodations: %v", err)
			http.Error(w, "Failed to fetch accommodations", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(accommodations)
	}

}

// Dodavanje ili ažuriranje dostupnosti
func UpdateAvailabilityHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var availability Availability
		if err := json.NewDecoder(r.Body).Decode(&availability); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Proveri da li postoji rezervacija u tom terminu pre nego što ažuriraš dostupnost

		query := "INSERT INTO availability (id, accommodation_id, start_date, end_date) VALUES (?, ?, ?, ?)"
		if err := session.Query(query, gocql.TimeUUID(), availability.AccommodationID, availability.StartDate, availability.EndDate).Exec(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Dodavanje ili ažuriranje cena
func UpdatePriceHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var price Price
		if err := json.NewDecoder(r.Body).Decode(&price); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Proveri da li postoji rezervacija u tom terminu pre nego što ažuriraš cenu

		query := "INSERT INTO prices (id, accommodation_id, start_date, end_date, amount, strategy) VALUES (?, ?, ?, ?, ?, ?)"
		if err := session.Query(query, gocql.TimeUUID(), price.AccommodationID, price.StartDate, price.EndDate, price.Amount, price.Strategy).Exec(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

//func UpdateAvailabilityAndPriceHandler(session *gocql.Session) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var data struct {
//			StartDate time.Time `json:"start_date"`
//			EndDate   time.Time `json:"end_date"`
//			Amount    float64   `json:"amount"`
//			Strategy  string    `json:"strategy"`
//		}
//
//		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
//			http.Error(w, "Invalid input", http.StatusBadRequest)
//			return
//		}
//
//		id := mux.Vars(r)["id"]
//
//		log.Println("Received start date:", data.StartDate)
//		log.Println("Received end date:", data.EndDate)
//
//		//// Konvertujte datume iz string formata u time.Time
//		//startDate, err := time.Parse("2006-01-02", data.StartDate) // Očekivani format: "yyyy-mm-dd"
//		//if err != nil {
//		//	http.Error(w, "Invalid start date format", http.StatusBadRequest)
//		//	return
//		//}
//		//
//		//endDate, err := time.Parse("2006-01-02", data.EndDate)
//		//if err != nil {
//		//	http.Error(w, "Invalid end date format", http.StatusBadRequest)
//		//	return
//		//}
//
//		// Ažuriranje dostupnosti
//		availabilityID := gocql.TimeUUID()
//		availabilityQuery := `INSERT INTO availability (id, accommodation_id, start_date, end_date) VALUES (?, ?, ?, ?)`
//		if err := session.Query(availabilityQuery, availabilityID, id, startDate, endDate).Exec(); err != nil {
//			log.Printf("Failed to update availability: %v", err)
//			http.Error(w, "Failed to update availability", http.StatusInternalServerError)
//			return
//		}
//
//		// Ažuriranje cene
//		priceID := gocql.TimeUUID()
//		priceQuery := `INSERT INTO prices (id, accommodation_id, start_date, end_date, amount, strategy) VALUES (?, ?, ?, ?, ?, ?)`
//		if err := session.Query(priceQuery, priceID, id, startDate, endDate, data.Amount, data.Strategy).Exec(); err != nil {
//			log.Printf("Failed to update price: %v", err)
//			http.Error(w, "Failed to update price", http.StatusInternalServerError)
//			return
//		}
//
//		w.WriteHeader(http.StatusOK)
//		json.NewEncoder(w).Encode("Availability and price updated successfully")
//	}
//}

func UpdateAvailabilityAndPriceHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
			Amount    float64   `json:"amount"`
			Strategy  string    `json:"strategy"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		id := mux.Vars(r)["id"]

		log.Println("Received start date:", data.StartDate)
		log.Println("Received end date:", data.EndDate)

		// Konvertovanje time.Time u Unix timestamp
		startDate := data.StartDate.Unix()
		endDate := data.EndDate.Unix()

		// Ažuriranje dostupnosti
		availabilityID := gocql.TimeUUID()
		availabilityQuery := `INSERT INTO availability (id, accommodation_id, start_date, end_date) VALUES (?, ?, ?, ?)`
		if err := session.Query(availabilityQuery, availabilityID, id, startDate, endDate).Exec(); err != nil {
			log.Printf("Failed to update availability: %v", err)
			http.Error(w, "Failed to update availability", http.StatusInternalServerError)
			return
		}

		// Ažuriranje cene
		priceID := gocql.TimeUUID()
		priceQuery := `INSERT INTO prices (id, accommodation_id, start_date, end_date, amount, strategy) VALUES (?, ?, ?, ?, ?, ?)`
		if err := session.Query(priceQuery, priceID, id, startDate, endDate, data.Amount, data.Strategy).Exec(); err != nil {
			log.Printf("Failed to update price: %v", err)
			http.Error(w, "Failed to update price", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Availability and price updated successfully")
	}
}

func GetAvailabilityHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accommodationID, err := gocql.ParseUUID(vars["id"])
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		var availabilities []Availability
		query := `SELECT id, accommodation_id, start_date, end_date FROM availability WHERE accommodation_id = ? ALLOW FILTERING`
		iter := session.Query(query, accommodationID).Iter()
		var availability Availability
		for iter.Scan(&availability.ID, &availability.AccommodationID, &availability.StartDate, &availability.EndDate) {
			availabilities = append(availabilities, availability)
		}

		if err := iter.Close(); err != nil {
			log.Printf("Failed to retrieve availabilities: %v", err)
			http.Error(w, "Failed to retrieve availabilities", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(availabilities)
	}
}
