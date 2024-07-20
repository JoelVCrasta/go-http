package bolt

import (
	"context"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// Bolt is a struct that holds the configuration for the BoltDB
type Bolt struct {
	db *bolt.DB
}

const (
	DBNAME     string = "users.db"
	BUCKETNAME string = "users"
)

// New creates a new instance of the BoltDB
func New(ctx context.Context, dir string) (*Bolt, error) {
	db, err := bolt.Open(fmt.Sprintf("%s/%s", dir, DBNAME), 0600, nil)
	if err != nil {
		return nil, err
	}

	// Create a bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Bolt{
		db: db,
	}, nil
}

// Close closes the database connection
func (b *Bolt) Close(ctx context.Context) {
	b.db.Close()
}

// Create creates a new record in the database
func (b *Bolt) Create(ctx context.Context, data []byte) error {

	fmt.Println(string(data))

	return nil
}
