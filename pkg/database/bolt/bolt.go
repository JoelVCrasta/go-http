package bolt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BUCKETNAME))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(user.Name), v)
	})

	return err
}

func (b *Bolt) Get(ctx context.Context, name string) *database.User {
	var raw []byte

	// Get the user from the database
	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKETNAME))
		raw = b.Get([]byte(name))
		return nil
	})

	if len(raw) == 0 {
		return nil
	}

	// Unmarshal the JSON
	var u database.User
	err := json.Unmarshal(raw, &u)
	if err != nil {
		log.Fatalf("Database corrupted: %v", err)
	}
	return &u
}

// Update updates a record in the database
func (b *Bolt) Update(ctx context.Context, user database.User) (*database.User, error) {
	var raw []byte

	// Get the user from the database
	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKETNAME))
		raw = b.Get([]byte(user.Name))
		return nil
	})

	// Unmarshal the JSON
	var cur database.User
	err := json.Unmarshal(raw, &cur)
	if err != nil {
		return nil, err
	}

	// Validate and update the user
	cur.Name = user.Name
	if user.Email != "" {
		cur.Email = user.Email
	}
	if user.Age != 0 {
		cur.Age = user.Age
	}

	//	Marshal the JSON
	v, err := json.Marshal(cur)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal JSON: %v", err))
	}

	// Update the user in the database
	err = b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKETNAME))
		return b.Put([]byte(user.Name), v)
	})
	if err != nil {
		return nil, err
	}

	return &cur, nil

}

// Delete deletes a record from the database
func (b *Bolt) Delete(ctx context.Context, name string) error {
	// Delete the user from the database
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKETNAME))
		return b.Delete([]byte(name))
	})
	return err
}
