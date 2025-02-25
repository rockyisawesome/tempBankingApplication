package repositories

import (
	// "accountservice/database"
	"accountservice/database"
	"accountservice/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// UserRepository implements Repository for User entities
type UserRepository struct {
	db *database.PostgresPoolDB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *database.PostgresPoolDB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, account *models.Account) error {
	query := "INSERT INTO  usersschema.accounts (account_number, username, email, balance, created_at, updated_at, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	conn, err := r.db.Pool().Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, query, account.AccountNumber, account.Username, account.Email, account.Balance, account.CreatedAt, account.UpdatedAt, account.IsActive).Scan(&account.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.Account, error) {
	query := "SELECT account_number, username, email FROM  usersschema.accounts WHERE account_number = $1"
	conn, err := r.db.Pool().Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	account := &models.Account{}
	err = conn.QueryRow(ctx, query, id).Scan(&account.AccountNumber, &account.Username, &account.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // or a custom "not found" error
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return account, nil
}

// CheckAccountExists checks if an account exists based on account_number
func (r *UserRepository) CheckAccountExists(ctx context.Context, accountNumber string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM usersschema.accounts WHERE account_number = $1)"
	conn, err := r.db.Pool().Acquire(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	var exists bool
	err = conn.QueryRow(ctx, query, accountNumber).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check account existence: %w", err)
	}
	return exists, nil
}

// UpdateBalance updates account balance with transaction support for ACID compliance
func (r *UserRepository) UpdateBalance(ctx context.Context, accountNumber string, amount float64, isCredit bool) error {
	// Get a connection and start a transaction
	conn, err := r.db.Pool().Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	// Begin transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Lock the row for update to ensure consistency
	query := `
        SELECT balance 
        FROM usersschema.accounts 
        WHERE account_number = $1 
        FOR UPDATE`

	var currentBalance float64
	err = tx.QueryRow(ctx, query, accountNumber).Scan(&currentBalance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("account not found: %s", accountNumber)
		}
		return fmt.Errorf("failed to get current balance: %w", err)
	}

	// Calculate new balance
	newBalance := currentBalance
	if isCredit {
		newBalance += amount
	} else {
		newBalance -= amount
		if newBalance < 0 {
			return fmt.Errorf("insufficient funds: current balance %.2f, attempted debit %.2f",
				currentBalance, amount)
		}
	}

	// Update balance and updated_at timestamp
	updateQuery := `
        UPDATE usersschema.accounts 
        SET balance = $1, 
            updated_at = NOW() 
        WHERE account_number = $2`

	result, err := tx.Exec(ctx, updateQuery, newBalance, accountNumber)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("account not found during update: %s", accountNumber)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Debit
func (r *UserRepository) Debit(ctx context.Context, accountNumber string, amount float64) error {
	if amount < 0 {
		return fmt.Errorf("debit amount cannot be negative: %.2f", amount)
	}
	return r.UpdateBalance(ctx, accountNumber, amount, false)
}

// credit
func (r *UserRepository) Credit(ctx context.Context, accountNumber string, amount float64) error {
	if amount < 0 {
		return fmt.Errorf("credit amount cannot be negative: %.2f", amount)
	}
	return r.UpdateBalance(ctx, accountNumber, amount, true)
}
