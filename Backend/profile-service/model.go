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

type User struct {
	ID           string    `json:"id" bson:"_id"`
	FirstName    string    `json:"firstName" bson:"firstname"`
	LastName     string    `json:"lastName" bson:"lastname"`
	Email        string    `json:"email" bson:"email"`
	PasswordHash string    `json:"-" bson:"passwordHash"`
	Age          int       `json:"age" bson:"age"`
	Country      string    `json:"country" bson:"country"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt"`
	Role         string    `json:"role" bson:"role"`
}

type Profile struct {
	ID        string    `json:"id" bson:"_id"`
	FirstName string    `json:"firstName" bson:"firstname"`
	LastName  string    `json:"lastName" bson:"lastname"`
	Email     string    `json:"email" bson:"email"`
	Age       int       `json:"age" bson:"age"`
	Location  string    `json:"location" bson:"country"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
