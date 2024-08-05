package main

import (
	"time"
)

type Role string

const (
	Unauthenticated Role = "NK"
	Host            Role = "H"
	Guest           Role = "G"
)

type User struct {
	ID           string    `json:"id" bson:"_id"`
	Username     string    `json:"username" bson:"username"`
	PasswordHash string    `json:"passwordHash" bson:"passwordHash"`
	Role         Role      `json:"role" bson:"role"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt"`
}
