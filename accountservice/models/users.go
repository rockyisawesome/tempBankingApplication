package models

import "github.com/google/uuid"

// User represents the users table in the database
type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
