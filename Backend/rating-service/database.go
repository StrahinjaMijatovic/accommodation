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

	createKeyspace(session)
	createTables(session)

	return session
}

func createKeyspace(session *gocql.Session) {
	query := `CREATE KEYSPACE IF NOT EXISTS accommodations WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};`
	if err := session.Query(query).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace: %v", err)
	}
}

func createTables(session *gocql.Session) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS host_ratings (id UUID PRIMARY KEY, user_id text, target_id UUID, rating int, comment text, rated_at timestamp)`,
		`CREATE TABLE IF NOT EXISTS accommodation_ratings (id UUID PRIMARY KEY, user_id text, target_id UUID, rating int, comment text, rated_at timestamp)`,
	}

	for _, query := range queries {
		if err := session.Query(query).Exec(); err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}
}
