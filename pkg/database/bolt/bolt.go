package bolt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/JoelVCrasta/go-http/pkg/database"
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

type userInfo struct {
	Email string `json:"email"`
	Age   uint8  `json:"age"`
}

// Close closes the database connection
func (b *Bolt) Close(ctx context.Context) {
	b.db.Close()
}

// Create creates a new record in the database
func (b *Bolt) Create(ctx context.Context, user database.User) error {

	// Create a new user
	userInfo := userInfo{
		Email: user.Email,
		Age:   user.Age,
	}

	v, err := json.Marshal(userInfo)
	if err != nil {
		return err
	}

	// Insert the user into the database
	b.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(BUCKETNAME))
		err = b.Put([]byte(user.Name), v)
		return err
	})

	return nil
}
