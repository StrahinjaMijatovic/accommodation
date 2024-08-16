package main

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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
func GetAccommodationHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := gocql.ParseUUID(vars["id"])
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		var acc Accommodation
		query := `SELECT id, name, location, guests, price FROM accommodations WHERE id = ? LIMIT 1`
		if err := session.Query(query, id).Scan(&acc.ID, &acc.Name, &acc.Location, &acc.Guests, &acc.Price); err != nil {
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
