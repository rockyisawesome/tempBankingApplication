package models

import (
	"time"
)

// Transaction represents a financial transaction between accounts.
// swagger:model Transaction
type Transaction struct {
	// The ID of the account from which the transaction originates.
	// Required: true
	// swagger:example "acc123"
	FromAccountID string `json:"from_account_id"` // Foreign key referencing the sender's Account.ID

	// The ID of the account to which the transaction is directed.
	// Required: true
	// swagger:example "acc456"
	ToAccountID string `json:"to_account_id"` // Foreign key referencing the recipient's Account.ID

	// The amount of money involved in the transaction.
	// Required: true
	// swagger:example 250.75
	Amount float64 `json:"amount"` // Transaction amount

	// The type of transaction (e.g., "transfer", "deposit", "withdrawal").
	// Required: true
	// swagger:example "transfer"
	TransactionType string `json:"transaction_type"` // Type of transaction (e.g., "transfer", "deposit", "withdrawal")

	// A brief description of the transaction.
	// swagger:example "Monthly rent payment"
	Description string `json:"description"` // Optional description of the transaction

	// The timestamp when the transaction was created.
	// swagger:example "2025-02-25T14:30:00Z"
	CreatedAt time.Time `json:"created_at"` // Timestamp of transaction creation

	// The current status of the transaction (e.g., "pending", "completed", "failed").
	// Required: true
	// swagger:example "pending"
	Status string `json:"status"` // Transaction status (e.g., "pending", "completed", "failed")
}
