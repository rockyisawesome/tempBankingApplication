package database

import "context"

type DB interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Ping(ctx context.Context) error
}
