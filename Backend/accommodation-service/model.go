package main

import (
	"time"

	"github.com/gocql/gocql"
)

type Accommodation struct {
	ID            gocql.UUID     `json:"id"`
	UserID        string         `json:"user_id"`
	Name          string         `json:"name"`
	Location      string         `json:"location"`
	Guests        int            `json:"guests"`
	BasePrice     float64        `json:"base_price"`
	PriceStrategy string         `json:"price_strategy"`
	Amenities     string         `json:"amenities"`
	Images        []string       `json:"images"`
	Availability  []Availability `json:"availability"`
	BookedPeriods []Availability `json:"booked_periods"`
	Price         float64        `json:"prices"`
}

type Availability struct {
	ID              gocql.UUID `json:"id"`
	AccommodationID gocql.UUID `json:"accommodation_id"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
}

type Price struct {
	ID              gocql.UUID `json:"id"`
	AccommodationID gocql.UUID `json:"accommodation_id"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	Amount          float64    `json:"amount"`
	Strategy        string     `json:"strategy"`
}
