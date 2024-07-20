package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/JoelVCrasta/go-http/pkg/database"
	"github.com/gorilla/mux"
)

var index string = `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			h1, p {
				text-align: center;
			}
		</style>
	</head>
	<body>
		<h1>User Database</h1>
		<p>Welcome to the DB</p>
	</body>
	</html>
`

// userInfo is a struct that holds the user information
/* type userInfo struct {
	email string
	age   uint8
} */

// Server is a struct that holds the server configuration
type Server struct {
	ctx context.Context
	db  database.Database
}

func New(ctx context.Context, db database.Database) *Server {
	return &Server{
		ctx: ctx,
		db:  db,
	}
}

// handleIndex handles the "/" route
func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(index))
}

// HandleUsers handles the "/users/create" route
func (s *Server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodPut:
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Unmarshal the JSON
		var u database.User
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("Create User: %v", u.Name)

		// Validate the user
		if u.Name == "" {
			log.Print("No name provided")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check if the user already exists
		got := s.db.Get(s.ctx, u.Name)
		if got != nil {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("User already exists: %v", u.Name)))
			log.Printf("User already exists: %v", u.Name)
			return
		}

		// Write to the database
		err = s.db.Create(s.ctx, u)
		if err != nil {
			log.Printf("Failed to create user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) HandleUser(w http.ResponseWriter, r *http.Request) {

	// 	// fetch the user from query params
	// 	// name := r.URL.Query().Get("name")

	params := mux.Vars(r) // get the URL params using mux
	name := params["name"]

	switch r.Method {
	case http.MethodGet:
		log.Printf("Get User: %s", name)

		// Get the user from the database
		ret := s.db.Get(s.ctx, name)
		if ret == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Add the name to the user
		ret = &database.User{
			Name:  name,
			Email: ret.Email,
			Age:   ret.Age,
		}

		// Marshal the JSON
		msg, err := json.Marshal(ret)
		if err != nil {
			log.Printf("Failed to marshal JSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(msg)

	case http.MethodPatch:
		// check if its json
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Unmarshal the JSON
		var u database.User
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate the user
		if u.Name == "" {
			log.Print("No name provided")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check if the user exists
		got := s.db.Get(s.ctx, u.Name)
		if got == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("User does not exist: %v", u.Name)))
			log.Printf("User does not exist: %v", u.Name)
			return
		}

		log.Printf("Update User: %s", name)

		user, err := s.db.Update(s.ctx, u)
		if err != nil {
			log.Printf("Failed to update user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Return the updated user
		msg, err := json.Marshal(user)
		if err != nil {
			log.Printf("Failed to marshal JSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(msg)

	case http.MethodDelete:
		log.Printf("Delete User: %s", name)

		// Check if the user exists
		got := s.db.Get(s.ctx, name)
		if got == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("User does not exist: %v", name)))
			log.Printf("User does not exist: %v", name)
			return
		}

		// Delete the user
		err := s.db.Delete(s.ctx, name)
		if err != nil {
			log.Printf("Failed to delete user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
