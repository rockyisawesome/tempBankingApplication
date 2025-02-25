package database

import (
	"context"
	"ledgerservice/models"
)

type Database interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	InsertTransaction(ctx context.Context, ledger models.TransactionLedger) (string, error)
}
