// model.go
package main

import (
	"time"

	"github.com/gocql/gocql"
)

type Reservation struct {
	ID              gocql.UUID `json:"id"`
	AccommodationID gocql.UUID `json:"accommodation_id"`
	GuestID         string     `json:"guest_id"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
}
