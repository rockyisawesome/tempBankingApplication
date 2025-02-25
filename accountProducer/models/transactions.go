package models

import (
	"time"
)

type Transaction struct {
	FromAccountID   string    `json:"from_account_id"`  // Foreign key referencing the sender's Account.ID
	ToAccountID     string    `json:"to_account_id"`    // Foreign key referencing the recipient's Account.ID
	Amount          float64   `json:"amount"`           // Transaction amount
	TransactionType string    `json:"transaction_type"` // Type of transaction (e.g., "transfer", "deposit", "withdrawal")
	Description     string    `json:"description"`      // Optional description of the transaction
	CreatedAt       time.Time `json:"created_at"`       // Timestamp of transaction creation
	Status          string    `json:"status"`           // Transaction status (e.g., "pending", "completed", "failed")
}
