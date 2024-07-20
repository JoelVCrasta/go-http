package database

import "context"

// represents the JSON object that will be sent to the server
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   uint8  `json:"age"`
}

// Database is an interface that represents the database.
// It abstracts the underlying implementation.
type Database interface {
	Create(ctx context.Context, user User) error
	Get(ctx context.Context, name string) *User
	Update(ctx context.Context, name User) (*User, error)
	Delete(ctx context.Context, name string) error
}
