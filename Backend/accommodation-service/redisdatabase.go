package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var ctx = context.Background()

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis adresa iz docker-compose.yml
		Password: "",           // Bez lozinke
		DB:       0,            // Koristi default bazu
	})

	// Test konekcije
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	return rdb
}
