package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile := Profile{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Age:       user.Age,
		Location:  user.Country,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var req Profile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email je obavezan", http.StatusBadRequest)
		return
	}

	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateFields := bson.M{
		"firstname": req.FirstName,
		"lastname":  req.LastName,

		"age":       req.Age,
		"country":   req.Location,
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
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func validatePassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var req struct {
		Email              string `json:"email"`
		CurrentPassword    string `json:"currentPassword"`
		NewPassword        string `json:"newPassword"`
		ConfirmNewPassword string `json:"confirmNewPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.NewPassword != req.ConfirmNewPassword {
		http.Error(w, "New passwords do not match", http.StatusBadRequest)
		return
	}

	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !validatePassword(req.CurrentPassword, user.PasswordHash) {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	newPasswordHash, err := hashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updateResult, err := collection.UpdateOne(ctx, bson.M{"email": req.Email}, bson.M{"$set": bson.M{"passwordHash": newPasswordHash}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if updateResult.MatchedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Password updated successfully")
}

func DeleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	vars := mux.Vars(r)
	userID := vars["userID"]

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	hasActiveReservations, err := checkActiveReservations(userID)
	if err != nil {
		http.Error(w, "Failed to check reservations", http.StatusInternalServerError)
		return
	}
	if hasActiveReservations {
		http.Error(w, "Cannot delete profile with active reservations", http.StatusForbidden)
		return
	}

	hasAccommodations, err := checkAccommodations(userID)
	if err != nil {
		http.Error(w, "Failed to check accommodations", http.StatusInternalServerError)
		return
	}
	if hasAccommodations {
		err = deleteAccommodations(userID)
		if err != nil {
			http.Error(w, "Failed to delete accommodations", http.StatusInternalServerError)
			return
		}
	}

	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Profile and accommodations deleted successfully")
}

func checkAccommodations(userID string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("http://accommodation-service:8080/accommodations/exists/%s", userID))
	if err != nil {
		log.Printf("Error checking accommodations: %v", err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to check accommodations, status code: %d", resp.StatusCode)
		return false, fmt.Errorf("failed to check accommodations, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return false, err
	}

	var hasAccommodations bool
	if err := json.Unmarshal(body, &hasAccommodations); err != nil {
		log.Printf("Error unmarshaling response body: %v", err)
		return false, err
	}

	return hasAccommodations, nil
}

func deleteAccommodations(userID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://accommodation-service:8080/accommodations/user/%s", userID), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete accommodations, status code: %d", resp.StatusCode)
	}

	return nil
}

func checkActiveReservations(userID string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("http://reservation-service:8081/reservations/active/%s", userID))
	if err != nil {
		log.Printf("Error checking active reservations: %v", err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to check reservations, status code: %d", resp.StatusCode)
		return false, fmt.Errorf("failed to check reservations, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return false, err
	}

	var hasActiveReservations bool
	if err := json.Unmarshal(body, &hasActiveReservations); err != nil {
		log.Printf("Error unmarshaling response body: %v", err)
		return false, err
	}

	return hasActiveReservations, nil
}
