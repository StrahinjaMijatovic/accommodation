package main

import (
	"github.com/gocql/gocql"
	"log"
)

// InitCassandra inicijalizuje konekciju sa Cassandra bazom
func InitCassandra() *gocql.Session {
	cluster := gocql.NewCluster("cassandra")
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}

	// Kreiraj keyspace ako ne postoji
	query := `CREATE KEYSPACE IF NOT EXISTS accommodations WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};`
	if err := session.Query(query).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace: %v", err)
	}

	// Prebaci na accommodations keyspace
	cluster.Keyspace = "accommodations"
	session.Close()

	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra keyspace: %v", err)
	}

	return session
}
