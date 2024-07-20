package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/JoelVCrasta/go-http/pkg/database"
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
type userInfo struct {
	email string
	age   uint8
}

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

		log.Printf("User: %v", u.Name)

		// Validate the user
		if u.Name == "" {
			log.Print("No name provided")
			w.WriteHeader(http.StatusBadRequest)
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

	// 	params := mux.Vars(r) // get the URL params using mux
	// 	name := params["name"]
	// 	u, ok := s.users[name] // check if the user exists
	// 	if !ok {
	// 		w.WriteHeader(http.StatusNotFound)
	// 		return
	// 	}

	// 	switch r.Method {
	// 	case http.MethodGet:
	// 		ret := user{
	// 			Name:  name,
	// 			Email: u.email,
	// 			Age:   u.age,
	// 		}

	// 		// Marshal the JSON
	// 		msg, err := json.Marshal(ret)
	// 		if err != nil {
	// 			log.Printf("Failed to marshal JSON: %v", err)
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}

	// 		log.Printf("Get User: %s", name)

	// 		w.Header().Add("Content-Type", "application/json")
	// 		w.Write(msg)

	// 	case http.MethodPatch:
	// 		// check if its json
	// 		contentType := r.Header.Get("Content-Type")
	// 		if contentType != "application/json" {
	// 			w.WriteHeader(http.StatusUnsupportedMediaType)
	// 			return
	// 		}

	// 		// Read the request body
	// 		body, err := io.ReadAll(r.Body)
	// 		if err != nil {
	// 			log.Printf("Failed to read request body: %v", err)
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}
	// 		defer r.Body.Close()

	// 		// Unmarshal the JSON
	// 		var u user
	// 		err = json.Unmarshal(body, &u)
	// 		if err != nil {
	// 			log.Printf("Failed to unmarshal JSON: %v", err)
	// 			w.WriteHeader(http.StatusBadRequest)
	// 			return
	// 		}

	// 		log.Printf("Update User: %s", name)

	// 		userInfo := s.users[name] // Get the user
	// 		if u.Age != 0 {
	// 			userInfo.age = u.Age
	// 		}
	// 		if u.Email != "" {
	// 			userInfo.email = u.Email
	// 		}

	// 		s.users[name] = userInfo

	// 	case http.MethodDelete:
	// 		log.Printf("Delete User: %s", name)

	// 		delete(s.users, name)

	// 	default:
	// 		w.WriteHeader(http.StatusMethodNotAllowed)
	// 	}

}
