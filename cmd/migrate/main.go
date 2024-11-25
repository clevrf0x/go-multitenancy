package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/clevrf0x/go-multitenancy/db"
)

func main() {
	command := flag.String("command", "up", "Migration command (up, down, status, redo)")
	flag.Parse()

	dbInstance := db.New()
	err := runMigration(*command, dbInstance.DB)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}

func runMigration(command string, dbInstance *sql.DB) error {
	sharedDir := "./db/migrations/public"
	tenantDir := "./db/migrations/tenant"

	switch command {
	case "up":
		return db.RunMigrations(dbInstance, sharedDir, tenantDir)
	case "down":
		return db.DowngradeLastMigration(dbInstance, sharedDir, tenantDir)
	case "status":
		return db.MigrationStatus(dbInstance, sharedDir, tenantDir)
	case "redo":
		return db.RedoLastMigration(dbInstance, sharedDir, tenantDir)
	default:
		log.Printf("unknown command: %s", command)
		return fmt.Errorf("unknown command: %s", command)
	}
}
