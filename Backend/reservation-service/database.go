// database.go
package main

import (
	"github.com/gocql/gocql"
	"log"
)

func InitCassandra() *gocql.Session {
	cluster := gocql.NewCluster("cassandra")
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}

	query := `CREATE KEYSPACE IF NOT EXISTS reservations WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};`
	if err := session.Query(query).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace: %v", err)
	}

	cluster.Keyspace = "reservations"
	session.Close()

	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra keyspace: %v", err)
	}

	query = `
		CREATE TABLE IF NOT EXISTS reservations (
			id UUID PRIMARY KEY,
			accommodation_id UUID,
			user_id UUID,
			start_date date,
			end_date date
		);
	`
	if err := session.Query(query).Exec(); err != nil {
		log.Fatalf("Failed to create reservations table: %v", err)
	}

	return session
}
