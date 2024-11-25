package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB, sharedDir string, tenantDir string) error {
	return runForAllSchemas(db, sharedDir, tenantDir, func(db *sql.DB, dir string) error {
		return goose.Up(db, dir)
	})
}

func DowngradeLastMigration(db *sql.DB, sharedDir string, tenantDir string) error {
	return runForAllSchemas(db, sharedDir, tenantDir, func(db *sql.DB, dir string) error {
		return goose.Down(db, dir)
	})
}

func MigrationStatus(db *sql.DB, sharedDir string, tenantDir string) error {
	return runForAllSchemas(db, sharedDir, tenantDir, func(db *sql.DB, dir string) error {
		return goose.Status(db, dir)
	})
}

func RedoLastMigration(db *sql.DB, sharedDir string, tenantDir string) error {
	return runForAllSchemas(db, sharedDir, tenantDir, func(db *sql.DB, dir string) error {
		if err := goose.Down(db, dir); err != nil {
			return err
		}
		return goose.Up(db, dir)
	})
}

func runForAllSchemas(db *sql.DB, sharedDir string, tenantDir string, migrationFunc func(*sql.DB, string) error) error {
	// Run for public schema
	if err := runForSchema(db, "public", sharedDir, migrationFunc); err != nil {
		return err
	}

	// Run for tenant schemas
	tenants, err := GetTenantsID(db)
	if err != nil {
		return err
	}

	for _, tenant := range tenants {
		if err := runForSchema(db, tenant, tenantDir, migrationFunc); err != nil {
			return err
		}
	}

	return nil
}

func runForSchema(db *sql.DB, schema, dir string, migrationFunc func(*sql.DB, string) error) error {
	// Set schema
	if _, err := db.Exec(fmt.Sprintf("SET search_path TO \"%s\"", schema)); err != nil {
		return err
	}

	// Ensure we switch back to public schema
	defer db.Exec("SET search_path TO public")

	log.Printf("Running migration for schema \"%s\"\n", schema)
	return migrationFunc(db, dir)
}
