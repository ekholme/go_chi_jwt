package gochijwt

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// define my types
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Secret   string `json:"secret"`
}

type Store struct {
	Users []*User
}

// add authservice to this
type Server struct {
	Store  *Store
	Router *chi.Mux
	Srvr   *http.Server
	//maybe some other stuff?
}

// constructor for a new server
func NewServer(r *chi.Mux) Server {

	var us []*User

	st := &Store{
		Users: us,
	}

	return Server{
		Store:  st,
		Router: r,
		Srvr: &http.Server{
			Addr: ":8080",
		},
	}
}

// run method for server
func (s Server) Run() {
	s.RegisterRoutes()

	fmt.Println("Example chi router running on port 8080")

	s.Srvr.Handler = s.Router

	s.Srvr.ListenAndServe()
}

// handlers
func (s Server) handleCreateNewUser(w http.ResponseWriter, r *http.Request) {

	var user *User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	s.Store.Users = append(s.Store.Users, user)

	msg := make(map[string]*User, 1)

	msg["user created"] = user

	writeJSON(w, http.StatusOK, msg)
}

func (s Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	us := s.Store.Users

	writeJSON(w, http.StatusOK, us)
}

func (s Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World"))
}

// register handlers
func (s Server) RegisterRoutes() {
	s.Router.Get("/", s.handleIndex)
	s.Router.Post("/register", s.handleCreateNewUser)
	s.Router.Get("/users", s.handleGetUsers)
}

// writeJSON helper
func writeJSON(w http.ResponseWriter, statusCode int, v any) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(v)
}
