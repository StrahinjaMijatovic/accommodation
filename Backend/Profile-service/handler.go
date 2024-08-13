package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	collection := db.Collection("profiles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var profile Profile
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Profil nije pronađen", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var req Profile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email je obavezan", http.StatusBadRequest)
		return
	}

	collection := db.Collection("profiles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateFields := bson.M{
		"firstName": req.FirstName,
		"lastName":  req.LastName,
		"gender":    req.Gender,
		"age":       req.Age,
		"location":  req.Location,
		"updatedAt": time.Now(),
	}

	update := bson.M{
		"$set": updateFields,
	}

	_, err := collection.UpdateOne(ctx, bson.M{"email": req.Email}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Profil uspešno ažuriran")
}
