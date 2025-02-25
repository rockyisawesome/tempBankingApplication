package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID            uuid.UUID `json:"id"`             // Unique identifier for the account
	AccountNumber string    `json:"account_number"` // Unique account number
	Username      string    `json:"username"`       // Account username
	Email         string    `json:"email"`          // Account email address
	Balance       float64   `json:"balance"`        // Account balance (e.g., for financial apps)
	CreatedAt     time.Time `json:"created_at"`     // Timestamp of account creation
	UpdatedAt     time.Time `json:"updated_at"`     // Timestamp of last update
	IsActive      bool      `json:"is_active"`      // Account status (active or inactive)
}

// NewAccount creates a new Account instance with default values
func NewAccount(username, email, password string) *Account {
	now := time.Now()
	return &Account{
		Username:  username,
		Email:     email,
		Balance:   0.0, // Default balance
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  true, // Default to active
	}
}
