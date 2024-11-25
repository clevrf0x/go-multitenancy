package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Service struct {
	*sql.DB
}

var dbInstance *Service

func New() *Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	connStr := os.Getenv("DATABASE_URI")
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &Service{db}
	return dbInstance
}

// SwitchSchema switches the current database schema to the provided schema name
func (s *Service) SwitchSchema(ctx context.Context, schemaName string) error {
	_, err := s.ExecContext(ctx, fmt.Sprintf("SET search_path TO \"%s\"", schemaName))
	if err != nil {
		return fmt.Errorf("failed to switch schema to %s: %w", schemaName, err)
	}
	return nil
}

// SwitchPublicSchema switches the current database schema back to the public schema
func (s *Service) SwitchPublicSchema(ctx context.Context) error {
	return s.SwitchSchema(ctx, "public")
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *Service) ConnClose() error {
	log.Printf("Disconnected from database: %s", os.Getenv("DB_DATABASE"))
	return s.Close()
}
