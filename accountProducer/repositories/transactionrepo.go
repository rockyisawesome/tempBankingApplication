package repositories

import (
	"accountProducer/database"
	"accountProducer/models"
	"context"

	"github.com/hashicorp/go-hclog"
)

type TransactionRepo struct {
	mgdb  database.Database
	loggs *hclog.Logger
}

func NewTransactionRepository(mgdb database.Database, lobbs *hclog.Logger) Repository {
	return &TransactionRepo{
		loggs: lobbs,
		mgdb:  mgdb,
	}
}

func (t *TransactionRepo) FindTransactionByAccountNumber(ctx context.Context, accountNumber string) (*[]models.TransactionLedger, error) {

	// Use MongoDB
	list, err := t.mgdb.GetTransactionsByAccountNumber(ctx, accountNumber)
	if err != nil {
		(*t.loggs).Error("Error Fetching the transactions", "Error", err)
		return nil, err
	}
	return list, nil
}
