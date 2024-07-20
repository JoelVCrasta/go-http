package database

import "context"

// Database is an interface that represents the database.
// It abstracts the underlying implementation.
type Database interface {
	Create(ctx context.Context, data []byte) error
}
