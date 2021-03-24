package models

import "time"

// User contains user properties
type User struct {
	ID        int64     `json:"user_id"`
	Name      string    `json:"user_name"`
	Email     string    `json:"user_mail"`
	CreatedAt time.Time `json:"createdAt"`
	Password  string    `json:"password,omitempty"`
}
