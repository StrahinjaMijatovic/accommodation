package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
	"time"
)

func CreateNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var notification Notification

	// Parse the request body into the notification struct
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set the creation time
	notification.CreatedAt = time.Now()

	session := Neo4jDriver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	// Save the notification in the Neo4j database
	_, err = session.Run(
		"CREATE (n:Notification {id: randomUUID(), host_id: $host_id, message: $message, created_at: $created_at})",
		map[string]interface{}{
			"host_id":    notification.HostID,
			"message":    notification.Message,
			"created_at": notification.CreatedAt,
		},
	)
	if err != nil {
		log.Printf("Failed to create notification: %v", err)
		http.Error(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}

//func GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	hostID := vars["hostId"]
//
//	log.Printf("Received request for notifications of host ID: %s", hostID)
//
//	session := Neo4jDriver.NewSession(neo4j.SessionConfig{})
//	defer session.Close()
//
//	// Retrieve notifications for the specified host ID
//	result, err := session.Run(
//		"MATCH (n:Notification {host_id: $host_id}) RETURN n.id AS id, n.message AS message, n.created_at AS created_at ORDER BY n.created_at DESC",
//		map[string]interface{}{
//			"host_id": hostID,
//		},
//	)
//	if err != nil {
//		log.Printf("Error retrieving notifications from the database: %v", err)
//		http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
//		return
//	}
//
//	var notifications []Notification
//	for result.Next() {
//		record := result.Record()
//
//		id, _ := record.Get("id")
//		message, _ := record.Get("message")
//		createdAt, _ := record.Get("created_at")
//
//		log.Printf("Fetched notification from database: ID=%d, Message=%s, CreatedAt=%s", id, message, createdAt)
//
//		notification := Notification{
//			ID:        id.(int64),
//			HostID:    hostID,
//			Message:   message.(string),
//			CreatedAt: createdAt.(time.Time),
//		}
//		notifications = append(notifications, notification)
//	}
//
//	if err = result.Err(); err != nil {
//		log.Printf("Error after fetching notifications: %v", err)
//		http.Error(w, "Error processing notifications", http.StatusInternalServerError)
//		return
//	}
//
//	log.Printf("Returning %d notifications to the client", len(notifications))
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(notifications)
//}

func GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostID := vars["hostId"]

	log.Printf("Received request for notifications of host ID: %s", hostID)

	session := Neo4jDriver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(
		"MATCH (n:Notification {host_id: $host_id}) RETURN n.id AS id, n.message AS message, n.created_at AS created_at ORDER BY n.created_at DESC",
		map[string]interface{}{
			"host_id": hostID,
		},
	)
	if err != nil {
		log.Printf("Error retrieving notifications from the database: %v", err)
		http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
		return
	}

	var notifications []Notification
	for result.Next() {
		record := result.Record()

		id, ok := record.Get("id")
		if !ok {
			log.Printf("Error fetching id from record: %v", result.Err())
			continue
		}

		message, ok := record.Get("message")
		if !ok {
			log.Printf("Error fetching message from record: %v", result.Err())
			continue
		}

		createdAt, ok := record.Get("created_at")
		if !ok {
			log.Printf("Error fetching created_at from record: %v", result.Err())
			continue
		}

		notification := Notification{
			ID:        id.(string), // Očekuje se string
			HostID:    hostID,
			Message:   message.(string),
			CreatedAt: createdAt.(time.Time),
		}
		notifications = append(notifications, notification)
	}

	if err = result.Err(); err != nil {
		log.Printf("Error after fetching notifications: %v", err)
		http.Error(w, "Error processing notifications", http.StatusInternalServerError)
		return
	}

	log.Printf("Returning %d notifications to the client", len(notifications))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}
