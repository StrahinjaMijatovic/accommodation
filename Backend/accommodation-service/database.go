package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gocql/gocql"
	"log"
	"strings"
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

var jwtKey = []byte("my_secret_key")

type Claims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

func ExtractUserIDFromToken(tokenString string) (string, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", errors.New("invalid token")
}
