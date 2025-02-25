package repositories

import (
	"accountProducer/models" // Importing models package for the TransactionLedger struct
	"context"                // Importing context for handling request-scoped values and cancellation
)

// Repository defines the interface for data access operations related to transactions.
// This interface abstracts the underlying data storage mechanism (e.g., database, in-memory store),
// enabling dependency injection and facilitating unit testing with mock implementations.
type Repository interface {
	// FindTransactionByAccountNumber retrieves all transactions associated with a given account number.
	// It accepts a context for cancellation and timeout support, and an accountNumber to filter transactions.
	// Returns a pointer to a slice of TransactionLedger structs, representing the transaction records,
	// or an error if the retrieval fails (e.g., due to storage unavailability or invalid account number).
	FindTransactionByAccountNumber(ctx context.Context, accountNumber string) (*[]models.TransactionLedger, error)
}
