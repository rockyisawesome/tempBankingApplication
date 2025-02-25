package repositories

import (
	"accountProducer/models"
	"context"
)

type Repository interface {
	FindTransactionByAccountNumber(ctx context.Context, accountNumber string) (*[]models.TransactionLedger, error)
}
