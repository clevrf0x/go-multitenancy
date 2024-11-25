package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	database "github.com/clevrf0x/go-multitenancy/db"
	db_sqlc "github.com/clevrf0x/go-multitenancy/db/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	DB        *database.Service
	Router    *chi.Mux
	Queries   *db_sqlc.Queries
	Validator *validator.Validate
}

func NewServer(db *database.Service) *Server {
	router := chi.NewRouter()
	queries := db_sqlc.New(db.DB)

	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())
	return &Server{
		DB:        db,
		Router:    router,
		Queries:   queries,
		Validator: validate,
	}
}

func (s *Server) Start() {
	port := os.Getenv("PORT")
	log.Printf("Server starting on http://0.0.0.0:%s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), s.Router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
