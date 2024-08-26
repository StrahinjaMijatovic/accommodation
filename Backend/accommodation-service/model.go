package main

import (
	"time"

	"github.com/gocql/gocql"
)

// Accommodation predstavlja osnovne informacije o smeštaju
type Accommodation struct {
	ID            gocql.UUID     `json:"id"`
	UserID        string         `json:"user_id"`
	Name          string         `json:"name"`
	Location      string         `json:"location"`
	Guests        int            `json:"guests"`
	BasePrice     float64        `json:"base_price"`     // Osnovna cena, može se koristiti kao default
	PriceStrategy string         `json:"price_strategy"` // Strategija cene (po gostu ili po jedinici)
	Amenities     string         `json:"amenities"`      // Dodatni sadržaji
	Images        []string       `json:"images"`         // Slike smeštaja
	Availability  []Availability `json:"availability"`   // Dostupni termini
	BookedPeriods []Availability `json:"booked_periods"` // Zauzeti termini
	Price         float64        `json:"prices"`         // Lista promenljivih cena
}

// Availability predstavlja period kada je smeštaj dostupan
type Availability struct {
	ID              gocql.UUID `json:"id"`
	AccommodationID gocql.UUID `json:"accommodation_id"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
}

// Price predstavlja promenljivu cenu za određeni period
type Price struct {
	ID              gocql.UUID `json:"id"`
	AccommodationID gocql.UUID `json:"accommodation_id"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	Amount          float64    `json:"amount"`
	Strategy        string     `json:"strategy"` // "per_guest" ili "per_unit"
}
