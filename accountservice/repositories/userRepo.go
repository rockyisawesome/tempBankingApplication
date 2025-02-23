package repositories

import (
	// "accountservice/database"
	"accountservice/database"
	"accountservice/models"
	"context"
	"fmt"

	"github.com/google/uuid"
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
	query := "INSERT INTO  usersschema.accounts (username, email, balance, created_at, updated_at, is_active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	conn, err := r.db.Pool().Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, query, account.Username, account.Email, account.Balance, account.CreatedAt, account.UpdatedAt, account.IsActive).Scan(&account.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	query := "SELECT id, username, email FROM  usersschema.accounts WHERE id = $1"
	conn, err := r.db.Pool().Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	account := &models.Account{}
	err = conn.QueryRow(ctx, query, id).Scan(&account.ID, &account.Username, &account.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // or a custom "not found" error
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return account, nil
}

// func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
// 	query := "UPDATE users SET name = $1, email = $2 WHERE id = $3"
// 	conn, err := r.db.Pool().Acquire(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to acquire connection: %w", err)
// 	}
// 	defer conn.Release()

// 	result, err := conn.Exec(ctx, query, user.Name, user.Email, user.ID)
// 	if err != nil {
// 		return fmt.Errorf("failed to update user: %w", err)
// 	}
// 	if result.RowsAffected() == 0 {
// 		return fmt.Errorf("no user found with id %d", user.ID)
// 	}
// 	return nil
// }

// func (r *UserRepository) Delete(ctx context.Context, id int) error {
// 	query := "DELETE FROM users WHERE id = $1"
// 	conn, err := r.db.Pool().Acquire(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to acquire connection: %w", err)
// 	}
// 	defer conn.Release()

// 	result, err := conn.Exec(ctx, query, id)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete user: %w", err)
// 	}
// 	if result.RowsAffected() == 0 {
// 		return fmt.Errorf("no user found with id %d", id)
// 	}
// 	return nil
// }

// func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
// 	query := "SELECT id, name, email FROM users ORDER BY id LIMIT $1 OFFSET $2"
// 	conn, err := r.db.Pool().Acquire(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to acquire connection: %w", err)
// 	}
// 	defer conn.Release()

// 	rows, err := conn.Query(ctx, query, limit, offset)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list users: %w", err)
// 	}
// 	defer rows.Close()

// 	var users []*models.User
// 	for rows.Next() {
// 		user := &models.User{}
// 		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
// 			return nil, fmt.Errorf("failed to scan user: %w", err)
// 		}
// 		users = append(users, user)
// 	}
// 	return users, nil
// }

// func (r *UserRepository) Count(ctx context.Context) (int64, error) {
// 	query := "SELECT COUNT(*) FROM users"
// 	conn, err := r.db.Pool().Acquire(ctx)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to acquire connection: %w", err)
// 	}
// 	defer conn.Release()

// 	var count int64
// 	err = conn.QueryRow(ctx, query).Scan(&count)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to count users: %w", err)
// 	}
// 	return count, nil
// }
