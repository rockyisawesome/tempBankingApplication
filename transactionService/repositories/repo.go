package repositories

import (
	"context"
	"transactionService/models"
)

type Repository interface {
	TransactionRouter(ctx context.Context, transmodel *models.Transaction) error
	CheckAccountExists(ctx context.Context, accountNumber string) (bool, error)
	UpdateBalance(ctx context.Context, accountNumber string, amount float64, isCredit bool) error
	Debit(ctx context.Context, accountNumber string, amount float64) error
	Credit(ctx context.Context, accountNumber string, amount float64) error
}
