package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresPoolDB implements the DB interface with connection pooling
type PostgresPoolDB struct {
	pool     *pgxpool.Pool
	dsn      string
	maxConns int32 // Maximum number of connections
	minConns int32 // Minimum number of connections
}

// NewPostgresPoolDB creates a new PostgresPoolDB instance with configurable max and min connections
func NewPostgresPoolDB(dsn string, maxConns, minConns int32) *PostgresPoolDB {
	return &PostgresPoolDB{
		dsn:      dsn,
		maxConns: maxConns,
		minConns: minConns,
	}
}

// Connect establishes a connection pool to the PostgreSQL database
func (p *PostgresPoolDB) Connect(ctx context.Context) error {
	// Parse the connection string into a config
	config, err := pgxpool.ParseConfig(p.dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Set pool configuration
	config.MaxConns = p.maxConns // Maximum number of connections in the pool
	config.MinConns = p.minConns // Minimum number of idle connections to maintain
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	// Create the connection pool with the configured settings
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	p.pool = pool
	return nil
}

// Close terminates all connections in the pool
func (p *PostgresPoolDB) Close(ctx context.Context) error {
	if p.pool == nil {
		return nil
	}
	p.pool.Close()
	p.pool = nil
	return nil
}

// Ping verifies the connection pool is still alive
func (p *PostgresPoolDB) Ping(ctx context.Context) error {
	if p.pool == nil {
		return fmt.Errorf("no active connection pool")
	}
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection for ping: %w", err)
	}
	defer conn.Release()

	err = conn.Ping(ctx)
	if err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	return nil
}

// Pool returns the underlying pgxpool.Pool for use in repositories
func (p *PostgresPoolDB) Pool() *pgxpool.Pool {
	return p.pool
}
