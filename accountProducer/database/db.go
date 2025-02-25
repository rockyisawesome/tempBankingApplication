package database

import (
	"accountProducer/models"
	"context"
)

type Database interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	GetTransactionsByAccountNumber(ctx context.Context, accountNumber string) (*[]models.TransactionLedger, error)
	// InsertTransaction(ctx context.Context, ledger models.TransactionLedger) (string, error)
	// we need to find transactions
}
