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
	ID           string    `json:"id" bson:"_id,omitempty"`
	Username     string    `json:"username" bson:"username"`
	PasswordHash string    `json:"passwordHash" bson:"passwordHash"`
	Email        string    `json:"email" bson:"email"`
	Age          int       `json:"age" bson:"age"`
	Country      string    `json:"country" bson:"country"`
	Role         Role      `json:"role" bson:"role"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt"`
}

type RegisterRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Email           string `json:"email"`
	Age             int    `json:"age"`
	Country         string `json:"country"`
	Role            string `json:"role"` // NK, H, G
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
