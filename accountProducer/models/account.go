package models

import (
	"time"

	"github.com/google/uuid"
)

// Account represents a user account in the system.
// swagger:model Account
type Account struct {
	// The unique identifier for the account.
	// swagger:example "550e8400-e29b-41d4-a716-446655440000"
	ID uuid.UUID `json:"id"` // Unique identifier for the account

	// The unique account number assigned to the account.
	// Required: true
	// swagger:example "ACC123456789"
	AccountNumber string `json:"account_number"` // Unique account number

	// The username chosen by the account holder.
	// Required: true
	// swagger:example "johndoe"
	Username string `json:"username"` // Account username

	// The email address associated with the account.
	// Required: true
	// swagger:example "john.doe@example.com"
	Email string `json:"email"` // Account email address

	// The current balance of the account.
	// Required: true
	// swagger:example 1000.50
	Balance float64 `json:"balance"` // Account balance (e.g., for financial apps)

	// The timestamp when the account was created.
	// swagger:example "2025-02-25T09:00:00Z"
	CreatedAt time.Time `json:"created_at"` // Timestamp of account creation

	// The timestamp when the account was last updated.
	// swagger:example "2025-02-25T10:00:00Z"
	UpdatedAt time.Time `json:"updated_at"` // Timestamp of last update

	// Indicates whether the account is active or inactive.
	// Required: true
	// swagger:example true
	IsActive bool `json:"is_active"` // Account status (active or inactive)
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
