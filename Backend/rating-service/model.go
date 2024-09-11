package main

import (
	"time"

	"github.com/gocql/gocql"
)

type Rating struct {
	ID       gocql.UUID `json:"id"`
	UserID   string     `json:"user_id"`
	TargetID gocql.UUID `json:"target_id"`
	Rating   int        `json:"rating"`
	Comment  string     `json:"comment"`
	RatedAt  time.Time  `json:"rated_at"`
}

type AverageRating struct {
	TargetID      gocql.UUID `json:"target_id"`
	AverageRating float64    `json:"average_rating"`
	TotalRatings  int        `json:"total_ratings"`
}
