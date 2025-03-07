package repositories

import (
	"accountservice/models"
	"context"
)

// Repository defines common database operations for a generic entity
// type Repository interface {
// 	Create(ctx context.Context, entity *T) error
// 	GetByID(ctx context.Context, id int) (*T, error)
// Update(ctx context.Context, entity *T) error
// Delete(ctx context.Context, id int) error
// List(ctx context.Context, limit, offset int) ([]*T, error)
// Count(ctx context.Context) (int64, error)
// }

type Repository interface {
	Create(ctx context.Context, user *models.Account) error
	GetByID(ctx context.Context, id string) (*models.Account, error)
	CheckAccountExists(ctx context.Context, accountNumber string) (bool, error)
	UpdateBalance(ctx context.Context, accountNumber string, amount float64, isCredit bool) error
	Debit(ctx context.Context, accountNumber string, amount float64) error
	Credit(ctx context.Context, accountNumber string, amount float64) error
}
