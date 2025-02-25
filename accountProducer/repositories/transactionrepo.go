package repositories

import (
	"accountProducer/database" // Importing database package for database operations
	"accountProducer/models"   // Importing models package for the TransactionLedger struct
	"context"                  // Importing context for handling request-scoped values and cancellation

	"github.com/hashicorp/go-hclog" // Importing hclog for structured logging
)

// TransactionRepo implements the Repository interface for transaction-related data operations.
// It uses a database instance for persistence and a logger for tracking operations and errors.
type TransactionRepo struct {
	mgdb  database.Database // mgdb is the database instance (e.g., MongoDB) for querying transactions
	loggs *hclog.Logger     // loggs is the logger instance for logging repository activities
}

// NewTransactionRepository creates a new TransactionRepo instance.
// It takes a database instance and a logger as dependencies, fulfilling the Repository interface.
// Returns a Repository interface type initialized with a TransactionRepo struct.
func NewTransactionRepository(mgdb database.Database, lobbs *hclog.Logger) Repository {
	return &TransactionRepo{
		loggs: lobbs, // Set the logger for logging operations
		mgdb:  mgdb,  // Set the database instance for data access
	}
}

// FindTransactionByAccountNumber retrieves all transactions for a given account number.
// Implements the Repository interface's method. It delegates the query to the underlying database
// instance and logs any errors. Takes a context for cancellation/timeout and an accountNumber
// to filter transactions. Returns a pointer to a slice of TransactionLedger structs or an error.
func (t *TransactionRepo) FindTransactionByAccountNumber(ctx context.Context, accountNumber string) (*[]models.TransactionLedger, error) {
	// Query the database for transactions using the provided account number
	list, err := t.mgdb.GetTransactionsByAccountNumber(ctx, accountNumber)
	if err != nil {
		(*t.loggs).Error("Error Fetching the transactions", "Error", err) // Log the error if retrieval fails
		return nil, err                                                   // Return the error to the caller
	}

	// Return the retrieved transaction list on success
	return list, nil
}
