package models

import (
	"errors"
	"fmt"
	"time"
)

// lower_case private, upper_case public
// Uppercase variable is mandatory for exposing to json

// UserResponse model
type UserResponse struct {
	ID        int64     `json:"user_id"`
	Username  string    `json:"user_name"`
	Email     string    `json:"user_mail"`
	Mobile    string    `json:"user_mobile"`
	CreatedAt time.Time `json:"createdAt"`
}

// UserCreate model
type UserCreate struct {
	ID        int64     `json:"user_id"`
	Username  string    `json:"user_name"`
	Email     string    `json:"user_mail"`
	Mobile    string    `json:"user_mobile"`
	CreatedAt time.Time `json:"createdAt"`
	Password  string    `json:"password"`
	Hash      string
}

// SPResponse model
type SPResponse struct {
	Msg string `json:"msg"`
	Id  int64  `json:"id"`
}

// GetUsername returns the username of an user
func (u *UserCreate) GetUsername() string {
	return u.Username
}

// Print method of User
func (u *UserCreate) Print() string {
	return fmt.Sprintf("%v (%v) - %v", u.Username, u.ID, u.CreatedAt)
}

// Validate returns an error when the username or email is to short.
func (u *UserCreate) Validate() error {
	if len(u.Username) < 10 {
		return ErrUsernameTooShort
	}

	if len(u.Email) < 10 {
		return ErrEmailTooShort
	}

	return nil
}

// Validate returns an error when the username or email is to short.
func (u *UserResponse) Validate() error {
	if len(u.Username) < 10 {
		return ErrUsernameTooShort
	}

	if len(u.Email) < 10 {
		return ErrEmailTooShort
	}

	return nil
}

// ValidatePassword returns an error if the password is to small.
func (u *UserCreate) ValidatePassword() error {
	if len(u.Password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}

// ErrUsernameTooShort is a error and is used when username is too short.
var ErrUsernameTooShort = errors.New("username is too short")

// ErrEmailTooShort is an error and is used when email address is too short.
var ErrEmailTooShort = errors.New("email address is to short")

// ErrPasswordTooShort is an error and is used when a password is too short.
var ErrPasswordTooShort = errors.New("password is to short")
