package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TransactionLedger represents a record of a transaction in the ledger.
// swagger:model TransactionLedger
type TransactionLedger struct {
	// The unique identifier for the transaction ledger entry.
	// swagger:example 507f1f77bcf86cd799439011
	ID bson.ObjectID `bson:"_id"`

	// The ID of the account from which the transaction originates.
	// Required: true
	// swagger:example "acc123"
	FromAccountID string `bson:"from_account_id" json:"from_account_id"`

	// The ID of the account to which the transaction is directed.
	// Required: true
	// swagger:example "acc456"
	ToAccountID string `bson:"to_account_id" json:"to_account_id"`

	// The amount of money involved in the transaction.
	// Required: true
	// swagger:example 100.50
	Amount float64 `bson:"amount" json:"amount"`

	// The type of transaction (e.g., "transfer", "deposit", "withdrawal").
	// Required: true
	// swagger:example "transfer"
	TransactionType string `bson:"transaction_type" json:"transaction_type"`

	// A brief description of the transaction.
	// swagger:example "Payment for services"
	Description string `bson:"description" json:"description"`

	// The timestamp when the transaction was created.
	// swagger:example "2025-02-25T10:00:00Z"
	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	// The current status of the transaction (e.g., "pending", "completed", "failed").
	// Required: true
	// swagger:example "completed"
	Status string `bson:"status" json:"status"`
}
