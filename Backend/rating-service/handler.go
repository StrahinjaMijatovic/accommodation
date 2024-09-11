package main

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func RateHostHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rating Rating

		log.Println("Received request to rate host")

		if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
			log.Printf("Failed to decode rating data: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		rating.ID = gocql.TimeUUID()
		rating.RatedAt = time.Now()

		vars := mux.Vars(r)
		targetID, err := gocql.ParseUUID(vars["id"])
		if err != nil {
			log.Printf("Invalid UUID format for TargetID: %v", err)
			http.Error(w, "Invalid UUID format", http.StatusBadRequest)
			return
		}
		rating.TargetID = targetID

		log.Printf("Inserting rating: %+v", rating)

		query := `INSERT INTO ratings (id, user_id, target_id, rating, comment, rated_at) VALUES (?, ?, ?, ?, ?, ?)`
		if err := session.Query(query, rating.ID, rating.UserID, rating.TargetID, rating.Rating, rating.Comment, rating.RatedAt).Exec(); err != nil {
			log.Printf("Failed to rate host: %v", err)
			http.Error(w, "Failed to rate host", http.StatusInternalServerError)
			return
		}

		log.Println("Rating successfully saved")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(rating)
	}
}

func DeleteHostRatingHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ratingID, err := gocql.ParseUUID(vars["id"])
		if err != nil {
			log.Printf("Invalid UUID format for rating ID: %v", err)
			http.Error(w, "Invalid UUID format", http.StatusBadRequest)
			return
		}

		query := `DELETE FROM ratings WHERE id = ?`
		if err := session.Query(query, ratingID).Exec(); err != nil {
			log.Printf("Failed to delete host rating: %v", err)
			http.Error(w, "Failed to delete host rating", http.StatusInternalServerError)
			return
		}

		log.Printf("Host rating with ID %s successfully deleted", ratingID)
		w.WriteHeader(http.StatusNoContent)
	}
}

func RateAccommodationHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rating Rating

		log.Println("Received request to rate accommodation")

		if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
			log.Printf("Failed to decode rating data: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		rating.ID = gocql.TimeUUID()
		rating.RatedAt = time.Now()

		vars := mux.Vars(r)
		targetID, err := gocql.ParseUUID(vars["id"])
		if err != nil {
			log.Printf("Invalid UUID format for TargetID: %v", err)
			http.Error(w, "Invalid UUID format", http.StatusBadRequest)
			return
		}
		rating.TargetID = targetID

		log.Printf("Inserting rating: %+v", rating)

		query := `INSERT INTO ratings2 (id, user_id, target_id, rating, comment, rated_at) VALUES (?, ?, ?, ?, ?, ?)`
		if err := session.Query(query, rating.ID, rating.UserID, rating.TargetID, rating.Rating, rating.Comment, rating.RatedAt).Exec(); err != nil {
			log.Printf("Failed to rate accommodation: %v", err)
			http.Error(w, "Failed to rate accommodation", http.StatusInternalServerError)
			return
		}

		log.Println("Rating successfully saved")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(rating)
	}
}

func DeleteAccommodationRatingHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ratingID, err := gocql.ParseUUID(vars["id"])
		if err != nil {
			log.Printf("Invalid UUID format for rating ID: %v", err)
			http.Error(w, "Invalid UUID format", http.StatusBadRequest)
			return
		}

		query := `DELETE FROM ratings2 WHERE id = ?`
		if err := session.Query(query, ratingID).Exec(); err != nil {
			log.Printf("Failed to delete accommodation rating: %v", err)
			http.Error(w, "Failed to delete accommodation rating", http.StatusInternalServerError)
			return
		}

		log.Printf("Accommodation rating with ID %s successfully deleted", ratingID)
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetHostRatingsHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ratings []Rating
		vars := mux.Vars(r)
		targetID, _ := gocql.ParseUUID(vars["id"])

		query := `SELECT id, user_id, target_id, rating, comment, rated_at FROM ratings WHERE target_id = ? ALLOW FILTERING`
		iter := session.Query(query, targetID).Iter()
		var rating Rating
		for iter.Scan(&rating.ID, &rating.UserID, &rating.TargetID, &rating.Rating, &rating.Comment, &rating.RatedAt) {
			ratings = append(ratings, rating)
		}
		if err := iter.Close(); err != nil {
			log.Printf("Failed to fetch ratings: %v", err)
			http.Error(w, "Failed to fetch ratings", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ratings)
	}
}

func GetAccommodationRatingsHandler(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ratings []Rating
		vars := mux.Vars(r)
		targetID, _ := gocql.ParseUUID(vars["id"])

		query := `SELECT id, user_id, target_id, rating, comment, rated_at FROM ratings2 WHERE target_id = ? ALLOW FILTERING`
		iter := session.Query(query, targetID).Iter()
		var rating Rating
		for iter.Scan(&rating.ID, &rating.UserID, &rating.TargetID, &rating.Rating, &rating.Comment, &rating.RatedAt) {
			ratings = append(ratings, rating)
		}
		if err := iter.Close(); err != nil {
			log.Printf("Failed to fetch ratings: %v", err)
			http.Error(w, "Failed to fetch ratings", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ratings)
	}
}
