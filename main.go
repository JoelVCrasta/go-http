package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"JoelVCrasta/go-http/pkg/database/bolt"
	"JoelVCrasta/go-http/pkg/server"

	"github.com/gorilla/mux"
)

func main() {
	// mux := http.NewServeMux() // This is the http mux
	address := ":8080"
	mux := mux.NewRouter()

	// Create a new context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// Create a new instance of the BoltDB
	b, err := bolt.New(ctx, "./data")
	if err != nil {
		log.Fatalf("Failed to start the database %v", err)
	}
	defer b.Close(ctx)

	srv := server.New(b)

	mux.HandleFunc("/", srv.HandleIndex)
	mux.HandleFunc("/users/create", srv.HandleCreateUser)
	mux.HandleFunc("/users/{name}", srv.HandleUser)

	s := &http.Server{
		Addr:           address,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Start server: %v", address)

	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
