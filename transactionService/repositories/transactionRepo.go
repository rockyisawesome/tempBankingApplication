package repositories

import (
	// "accountservice/database"
	"context"
	"fmt"
	"transactionService/database"
	"transactionService/models"

	"github.com/jackc/pgx/v5"
)

// UserRepository implements Repository for User entities
type TransactionRepository struct {
	db *database.PostgresPoolDB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *database.PostgresPoolDB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) TransactionRouter(ctx context.Context, transmodel *models.Transaction) error {
	// check account exists or not
	_, err := r.CheckAccountExists(ctx, transmodel.FromAccountID)
	if err != nil {
		return err
	}

	if transmodel.TransactionType == "deposit" {
		err := r.Credit(ctx, transmodel.FromAccountID, transmodel.Amount)
		if err != nil {
			return err
		}
	} else if transmodel.TransactionType == "withdrawal" {
		err := r.Debit(ctx, transmodel.FromAccountID, transmodel.Amount)
		if err != nil {
			return err
		}
	} else if transmodel.TransactionType == "transfer" {
		err := r.TransferAmount(ctx, transmodel.FromAccountID, transmodel.ToAccountID, transmodel.Amount)
		if err != nil {
			return err
		}
	}
	return nil
}

// CheckAccountExists checks if an account exists based on account_number
func (r *TransactionRepository) CheckAccountExists(ctx context.Context, accountNumber string) (bool, error) {
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
func (r *TransactionRepository) UpdateBalance(ctx context.Context, accountNumber string, amount float64, isCredit bool) error {
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
func (r *TransactionRepository) Debit(ctx context.Context, accountNumber string, amount float64) error {
	if amount < 0 {
		return fmt.Errorf("debit amount cannot be negative: %.2f", amount)
	}
	return r.UpdateBalance(ctx, accountNumber, amount, false)
}

// credit
func (r *TransactionRepository) Credit(ctx context.Context, accountNumber string, amount float64) error {
	if amount < 0 {
		return fmt.Errorf("credit amount cannot be negative: %.2f", amount)
	}
	return r.UpdateBalance(ctx, accountNumber, amount, true)
}

// TransferAmount transfers money from one account to another with ACID compliance
func (r *TransactionRepository) TransferAmount(ctx context.Context, fromAccountNumber, toAccountNumber string, amount float64) error {
	// Validate input
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be positive: %.2f", amount)
	}
	if fromAccountNumber == toAccountNumber {
		return fmt.Errorf("cannot transfer to the same account: %s", fromAccountNumber)
	}

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

	// Lock both accounts for update to prevent race conditions
	// Order by account_number to avoid deadlocks (consistent locking order)
	query := "SELECT account_number, balance FROM usersschema.accounts WHERE account_number IN ($1, $2) FOR UPDATE ORDER BY account_number"

	rows, err := tx.Query(ctx, query, fromAccountNumber, toAccountNumber)
	if err != nil {
		return fmt.Errorf("failed to lock accounts: %w", err)
	}
	defer rows.Close()

	// Collect balances and verify both accounts exist
	balances := make(map[string]float64)
	for rows.Next() {
		var accountNumber string
		var balance float64
		if err := rows.Scan(&accountNumber, &balance); err != nil {
			return fmt.Errorf("failed to scan account data: %w", err)
		}
		balances[accountNumber] = balance
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating account rows: %w", err)
	}

	// Check if both accounts were found
	if _, fromExists := balances[fromAccountNumber]; !fromExists {
		return fmt.Errorf("source account not found: %s", fromAccountNumber)
	}
	if _, toExists := balances[toAccountNumber]; !toExists {
		return fmt.Errorf("destination account not found: %s", toAccountNumber)
	}

	// Verify sufficient funds
	fromBalance := balances[fromAccountNumber]
	if fromBalance < amount {
		return fmt.Errorf("insufficient funds in %s: current balance %.2f, transfer amount %.2f",
			fromAccountNumber, fromBalance, amount)
	}

	// Calculate new balances
	newFromBalance := fromBalance - amount
	newToBalance := balances[toAccountNumber] + amount

	// Update both accounts in a single transaction
	updateQuery := `
	    UPDATE usersschema.accounts
	    SET balance = CASE
	                    WHEN account_number = $1 THEN $2
	                    WHEN account_number = $3 THEN $4
	                  END,
	        updated_at = NOW()
	    WHERE account_number IN ($1, $3)`

	result, err := tx.Exec(ctx, updateQuery, fromAccountNumber, newFromBalance, toAccountNumber, newToBalance)
	if err != nil {
		return fmt.Errorf("failed to update balances: %w", err)
	}

	// Update balance and updated_at timestamp

	if rowsAffected := result.RowsAffected(); rowsAffected != 1 {
		return fmt.Errorf("expected 2 rows affected, got %d", rowsAffected)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
