package repositories

import (
	"context"
	"ledgerservice/models"
)

type Repository interface {
	InsertTransaction(ctx context.Context, ledger models.TransactionLedger) error
}
