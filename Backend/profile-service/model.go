package main

import (
	"time"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
	Other  Gender = "other"
)

type Profile struct {
	ID        string    `json:"id" bson:"_id"`
	FirstName string    `json:"firstName" bson:"firstName"`
	LastName  string    `json:"lastName" bson:"lastName"`
	Email     string    `json:"email" bson:"email"`
	Gender    Gender    `json:"gender" bson:"gender"`
	Age       int       `json:"age" bson:"age"`
	Location  string    `json:"location" bson:"location"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
