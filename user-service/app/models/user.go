package models

// User is the model for Users on DB
type User struct {
	ID     int64  `json:"user_id"`
	Name   string `json:"user_name"`
	Email  string `json:"user_mail"`
	Mobile string `json:"user_mobile"`
}
