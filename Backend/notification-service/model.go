package main

import (
	"time"
)

type Notification struct {
	ID        string    `json:"id"`
	HostID    string    `json:"host_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
