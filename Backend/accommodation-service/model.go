package main

import "github.com/gocql/gocql"

type Accommodation struct {
	ID        gocql.UUID `json:"id"`
	Name      string     `json:"name"`
	Location  string     `json:"location"`
	Guests    int        `json:"guests"`
	Price     float64    `json:"price"`
	Amenities []string   `json:"amenities"`
	Images    []string   `json:"images"` // Dodaj polje za slike
}
