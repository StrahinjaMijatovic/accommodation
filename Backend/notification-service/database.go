package main

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
)

var Neo4jDriver neo4j.Driver

func InitNeo4j() error {
	uri := os.Getenv("NEO4J_URI")
	username := os.Getenv("NEO4J_USERNAME")
	password := os.Getenv("NEO4J_PASSWORD")

	var err error
	Neo4jDriver, err = neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return err
	}

	session := Neo4jDriver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	if _, err = session.Run("RETURN 1", nil); err != nil {
		return err
	}

	log.Println("Connected to Neo4j!")
	return nil
}

func CloseNeo4j() {
	if Neo4jDriver != nil {
		Neo4jDriver.Close()
	}
}
