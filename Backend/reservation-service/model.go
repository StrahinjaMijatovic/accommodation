package main

import "time"

type Reservation struct {
	ID              int64     `json:"id"`
	AccommodationID int64     `json:"accommodation_id"`
	GuestID         int64     `json:"guest_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Status          string    `json:"status"`
}

func (r *Reservation) IsOverlapping(startDate, endDate time.Time) bool {
	return r.StartDate.Before(endDate) && r.EndDate.After(startDate)
}
