package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

//	func CreateAccommodationHandler(session *gocql.Session) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			tokenString := r.Header.Get("Authorization")
//			userID, err := ExtractUserIDFromToken(tokenString)
//			if err != nil {
//				http.Error(w, "Unauthorized", http.StatusUnauthorized)
//				return
//			}
//
//			var acc Accommodation
//			if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
//				http.Error(w, "Invalid input", http.StatusBadRequest)
//				return
//			}
//			acc.ID = gocql.TimeUUID()
//			acc.UserID = userID
//
//			query := `INSERT INTO accommodations (id, user_id, name, location, guests, price, amenities, images) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
//			if err := session.Query(query, acc.ID, acc.UserID, acc.Name, acc.Location, acc.Guests, acc.Price, acc.Amenities, acc.Images).Exec(); err != nil {
//				log.Printf("Failed to create accommodation: %v", err)
//				http.Error(w, "Failed to create accommodation", http.StatusInternalServerError)
//				return
//			}
//
//			w.Header().Set("Content-Type", "application/json")
//			json.NewEncoder(w).Encode(acc)
//		}
//	}
func CreateAccommodationHandler(session *gocql.Session, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		userID, err := ExtractUserIDFromToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var acc Accommodation
		if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		acc.ID = gocql.TimeUUID()
		acc.UserID = userID

		// Keširanje slika u Redis-u
		for i, image := range acc.Images {
			redisKey := "accommodation:" + acc.ID.String() + ":image:" + strconv.Itoa(i)
			err := rdb.Set(ctx, redisKey, image, 24*time.Hour).Err() // Keširaj sliku na 24 sata
			if err != nil {
				log.Printf("Failed to cache image: %v", err)
				http.Error(w, "Failed to cache image", http.StatusInternalServerError)
				return
			}
		}

		// Kreiranje smeštaja u Cassandri
		query := `INSERT INTO accommodations (id, user_id, name, location, guests, price, amenities, images) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		if err := session.Query(query, acc.ID, acc.UserID, acc.Name, acc.Location, acc.Guests, acc.Price, acc.Amenities, acc.Images).Exec(); err != nil {
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
		query := `SELECT id, name, location, guests, price, amenities, images, user_id FROM accommodations WHERE id = ? LIMIT 1`
		if err := session.Query(query, id).Scan(
			&acc.ID,
			&acc.Name,
			&acc.Location,
			&acc.Guests,
			&acc.Price,
			&acc.Amenities,
			&acc.Images,
			&acc.UserID); err != nil {
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

		query := `SELECT id, name, location, guests, price, amenities, images, user_id FROM accommodations`
		iter := session.Query(query).Iter()
		var acc Accommodation
		for iter.Scan(&acc.ID, &acc.Name, &acc.Location, &acc.Guests, &acc.Price, &acc.Amenities, &acc.Images, &acc.UserID) {
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
		query := "SELECT id, name, location, guests, price, amenities, images, user_id FROM accommodations WHERE id = ? LIMIT 1"
		if err := session.Query(query, uuid).Scan(
			&accommodation.ID,
			&accommodation.Name,
			&accommodation.Location,
			&accommodation.Guests,
			&accommodation.Price,
			&accommodation.Amenities,
			&accommodation.Images,
			&accommodation.UserID,
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
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")

		var accommodations []Accommodation
		var query string
		var params []interface{}

		// Proveri dostupnost za zadate datume
		if startDateStr != "" && endDateStr != "" {
			startDate, err := time.Parse("2006-01-02", startDateStr)
			if err != nil {
				http.Error(w, "Invalid start date", http.StatusBadRequest)
				return
			}
			endDate, err := time.Parse("2006-01-02", endDateStr)
			if err != nil {
				http.Error(w, "Invalid end date", http.StatusBadRequest)
				return
			}

			// Filtriraj akomodacije po dostupnosti koristeći ALLOW FILTERING
			availabilityQuery := `SELECT accommodation_id FROM availability WHERE start_date <= ? AND end_date >= ? ALLOW FILTERING`
			iter := session.Query(availabilityQuery, endDate, startDate).Iter()
			var availableAccommodationIDs []string
			var accID string
			for iter.Scan(&accID) {
				availableAccommodationIDs = append(availableAccommodationIDs, accID)
			}
			if err := iter.Close(); err != nil {
				log.Printf("Failed to fetch available accommodations: %v", err)
				http.Error(w, "Failed to fetch available accommodations", http.StatusInternalServerError)
				return
			}

			// Ako nema dostupnih akomodacija, vrati prazan rezultat
			if len(availableAccommodationIDs) == 0 {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(accommodations)
				return
			}

			// Pripremi upit za pretragu akomodacija na osnovu dobijenih ID-eva
			query = "SELECT id, name, location, guests, price, amenities, images FROM accommodations WHERE id IN ("
			for i, id := range availableAccommodationIDs {
				if i > 0 {
					query += ", "
				}
				query += "?"
				params = append(params, id)
			}
			query += ")"

			// Dodaj filtere za lokaciju i broj gostiju
			if location != "" {
				query += " AND location = ?"
				params = append(params, location)
			}

			if guestsStr != "" {
				guests, err := strconv.Atoi(guestsStr)
				if err != nil {
					http.Error(w, "Invalid number of guests", http.StatusBadRequest)
					return
				}
				query += " AND guests >= ?"
				params = append(params, guests)
			}
		} else {
			// Ako nema zadatih datuma, vrati sve akomodacije koje zadovoljavaju ostale filtere
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

			if allowFiltering {
				query += " ALLOW FILTERING"
			}
		}

		// Izvrši upit za akomodacije
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
		// Define a struct to parse the input JSON data
		var data struct {
			StartDate string  `json:"startDate"`
			EndDate   string  `json:"endDate"`
			Amount    float64 `json:"amount"`
			Strategy  string  `json:"strategy"`
		}

		// Parse the JSON input from the request
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Parse dates from strings to time.Time
		startDate, err := time.Parse("2006-01-02", data.StartDate) // Format: 'YYYY-MM-DD'
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", data.EndDate) // Format: 'YYYY-MM-DD'
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}

		// Extract the accommodation ID from the URL parameters
		vars := mux.Vars(r)
		idStr := vars["id"]

		accommodationID, err := gocql.ParseUUID(idStr)
		if err != nil {
			http.Error(w, "Invalid accommodation ID format", http.StatusBadRequest)
			return
		}

		// Log the received data for debugging
		log.Println("Received start date:", startDate)
		log.Println("Received end date:", endDate)
		log.Println("Amount:", data.Amount)
		log.Println("Strategy:", data.Strategy)

		// Create and insert new availability record
		availabilityID := gocql.TimeUUID()
		availabilityQuery := `INSERT INTO availability (id, accommodation_id, start_date, end_date) VALUES (?, ?, ?, ?)`
		if err := session.Query(availabilityQuery, availabilityID, accommodationID, startDate, endDate).Exec(); err != nil {
			log.Printf("Failed to update availability: %v", err)
			http.Error(w, "Failed to update availability", http.StatusInternalServerError)
			return
		}

		// Create and insert new price record
		priceID := gocql.TimeUUID()
		priceQuery := `INSERT INTO prices (id, accommodation_id, start_date, end_date, amount, strategy) VALUES (?, ?, ?, ?, ?, ?)`
		if err := session.Query(priceQuery, priceID, accommodationID, startDate, endDate, data.Amount, data.Strategy).Exec(); err != nil {
			log.Printf("Failed to update price: %v", err)
			http.Error(w, "Failed to update price", http.StatusInternalServerError)
			return
		}

		// Send success response
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
func GetAvailabilityByAccommodationIDHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accommodationID := vars["id"]

		query := `SELECT id, start_date, end_date FROM availability WHERE accommodation_id = ?`
		iter := session.Query(query, accommodationID).Iter()

		var availabilityList []Availability
		var availability Availability
		for iter.Scan(&availability.ID, &availability.StartDate, &availability.EndDate) {
			availabilityList = append(availabilityList, availability)
		}
		if err := iter.Close(); err != nil {
			log.Printf("Failed to fetch availability: %v", err)
			http.Error(w, "Failed to fetch availability", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(availabilityList)
	}
}

// GetPricesHandler vraća sve cene za određeni smeštaj na osnovu ID-a smeštaja
func GetPricesHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accommodationID := vars["id"]

		var prices []Price
		query := `SELECT id, accommodation_id, start_date, end_date, amount, strategy FROM prices WHERE accommodation_id = ? ALLOW FILTERING`
		iter := session.Query(query, accommodationID).Iter()

		var price Price
		for iter.Scan(&price.ID, &price.AccommodationID, &price.StartDate, &price.EndDate, &price.Amount, &price.Strategy) {
			prices = append(prices, price)
		}

		if err := iter.Close(); err != nil {
			log.Printf("Failed to retrieve prices: %v", err)
			http.Error(w, "Failed to retrieve prices", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prices)
	}
}

func GetMyAccommodationsHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Caoooo")

		vars := mux.Vars(r)
		id := vars["id"]

		log.Println(id)

		log.Println("Received request:", r.Method, r.URL)

		log.Println("Handler GetMyAccommodationsHandler started")

		log.Println("UUID format is valid")

		// Priprema upita za bazu podataka
		query := `SELECT id, name, location, guests, base_price, price_strategy, amenities, images FROM accommodations WHERE user_id = ? ALLOW FILTERING`
		iter := session.Query(query, id).Iter()

		log.Println("Query executed")

		// Priprema slice-a za čuvanje rezultata upita
		var accommodations []Accommodation

		// Iteracija kroz rezultate i popunjavanje slice-a
		var acc Accommodation
		for iter.Scan(&acc.ID, &acc.Name, &acc.Location, &acc.Guests, &acc.BasePrice, &acc.PriceStrategy, &acc.Amenities, &acc.Images) {
			acc.Price = acc.BasePrice
			accommodations = append(accommodations, acc)
		}

		log.Println("Scan completed")

		// Provera za greške tokom iteracije
		if err := iter.Close(); err != nil {
			log.Printf("Failed to fetch accommodations: %v", err)
			http.Error(w, "Failed to fetch accommodations", http.StatusInternalServerError)
			return
		}

		log.Println("Iteration closed successfully")

		// Postavljanje Content-Type header-a i slanje JSON odgovora
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(accommodations); err != nil {
			log.Printf("Failed to encode response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Println("Response sent successfully")
	}
}

func DeleteAccommodationsByUserIDHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"]

		// Pronađi sve akomodacije korisnika
		query := `SELECT id FROM accommodations WHERE user_id = ? ALLOW FILTERING`
		iter := session.Query(query, userID).Iter()
		var accID gocql.UUID
		var ids []gocql.UUID

		for iter.Scan(&accID) {
			ids = append(ids, accID)
		}
		if err := iter.Close(); err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch accommodations: %v", err), http.StatusInternalServerError)
			return
		}

		// Obriši sve akomodacije
		for _, id := range ids {
			query = `DELETE FROM accommodations WHERE id = ?`
			if err := session.Query(query, id).Exec(); err != nil {
				log.Printf("Failed to delete accommodation: %v", err)
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
func HasAccommodationsHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"]

		// Proverava da li korisnik ima akomodacije u bazi
		query := `SELECT COUNT(*) FROM accommodations WHERE user_id = ? ALLOW FILTERING`
		var count int
		if err := session.Query(query, userID).Scan(&count); err != nil {
			http.Error(w, "Failed to check accommodations", http.StatusInternalServerError)
			return
		}

		hasAccommodations := count > 0
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(hasAccommodations)
	}
}
func DeleteAccommodationsByUserHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"]

		query := `DELETE FROM accommodations WHERE user_id = ? ALLOW FILTERING`
		if err := session.Query(query, userID).Exec(); err != nil {
			http.Error(w, "Failed to delete accommodations", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Accommodations deleted successfully")
	}
}
