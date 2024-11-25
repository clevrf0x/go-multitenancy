package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	db_sqlc "github.com/clevrf0x/go-multitenancy/db/sqlc"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)

func CreateTenant(db *sql.DB, name, workspace string) (db_sqlc.CreateTenantRow, error) {
	tenant := db_sqlc.CreateTenantRow{}

	queries := db_sqlc.New(db)
	tenantID, err := uuid.NewRandom()
	if err != nil {
		return tenant, fmt.Errorf("failed to generate uuid: %w", err)
	}
	// Create Schema for tenant
	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA \"%s\"", tenantID.String()))
	if err != nil {
		return tenant, fmt.Errorf("failed to create schema: %w", err)
	}

	// Apply migrations
	if err = ApplyTenantMigrations(db, tenantID.String(), queries); err != nil {
		return tenant, err
	}

	// Create tenant record inside shared Schema
	tenant, err = queries.CreateTenant(context.Background(), db_sqlc.CreateTenantParams{
		ID:        tenantID,
		Name:      name,
		Workspace: workspace,
	})
	if err != nil {
		return tenant, fmt.Errorf(
			"failed to insert tenant record, manually remove %v schema: %w",
			tenantID,
			err,
		)
	}
	return tenant, nil
}

func GetTenantsID(db *sql.DB) ([]string, error) {
	// Implement logic to retrieve list of tenants
	tenantIDs := []string{}
	queries := db_sqlc.New(db)

	tenants, err := queries.ListAllTenants(context.Background())
	if err != nil {
		return tenantIDs, err
	}

	for _, tenant := range tenants {
		tenantIDs = append(tenantIDs, tenant.ID.String())
	}

	return tenantIDs, nil
}

func ApplyTenantMigrations(db *sql.DB, tenantID string, queries *db_sqlc.Queries) error {
	// Set schema
	if _, err := db.Exec(fmt.Sprintf("SET search_path TO \"%s\"", tenantID)); err != nil {
		return err
	}

	// Ensure we switch back to public schema
	defer db.Exec("SET search_path TO public")
	log.Printf("Running migration for schema \"%s\"\n", tenantID)
	if err := goose.Up(db, "./db/migrations/tenant"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Seed default data
	SeedTenantDefaultData(db, queries, tenantID)

	return nil
}

func SeedTenantDefaultData(db *sql.DB, queries *db_sqlc.Queries, tenantID string) error {
	// Set schema
	if _, err := db.Exec(fmt.Sprintf("SET search_path TO \"%s\"", tenantID)); err != nil {
		return err
	}
	// Execute seed queries here if needed

	return nil
}
