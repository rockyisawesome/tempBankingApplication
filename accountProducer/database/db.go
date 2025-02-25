package database

import (
	"accountProducer/models" // Importing the models package to use TransactionLedger struct
	"context"                // Importing context for handling request-scoped values and cancellation
)

// Database defines the interface for database operations related to account transactions.
// This interface abstracts the underlying database implementation, allowing for flexibility
// in choosing the storage backend (e.g., SQL, NoSQL) and facilitating unit testing with mocks.
type Database interface {
	// Connect establishes a connection to the database.
	// It takes a context to allow for cancellation and timeouts, ensuring the operation
	// can be gracefully interrupted if needed.
	// Returns an error if the connection fails (e.g., due to invalid credentials or network issues).
	Connect(ctx context.Context) error

	// Disconnect closes the connection to the database.
	// It uses a context to manage the disconnection process, allowing for timeouts or cancellation.
	// Returns an error if the disconnection fails (e.g., due to resource cleanup issues).
	Disconnect(ctx context.Context) error

	// GetTransactionsByAccountNumber retrieves all transactions associated with a given account number.
	// The method accepts a context for cancellation and timeout support, and an accountNumber to filter
	// the transactions. It returns a pointer to a slice of TransactionLedger structs from the models
	// package, which represent the transaction records, or an error if the query fails (e.g., due to
	// database unavailability or invalid account number).
	GetTransactionsByAccountNumber(ctx context.Context, accountNumber string) (*[]models.TransactionLedger, error)
}
