package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	err := InitNeo4j()
	if err != nil {
		log.Fatal("Could not connect to Neo4j", err)
	}
	defer CloseNeo4j()

	router := mux.NewRouter()
	router.HandleFunc("/notifications", CreateNotificationHandler).Methods("POST")
	router.HandleFunc("/notifications/{hostId}", GetNotificationsHandler).Methods("GET")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Println("Server started on :8083")
	log.Fatal(http.ListenAndServe(":8083", handlers.CORS(headers, methods, origins)(router)))

	//// Get the port from environment or set to default
	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = "8083"
	//}
	//
	//// Start the server
	//log.Printf("Server is running on port %s\n", port)
	//log.Fatal(http.ListenAndServe(":"+port, router))
}
