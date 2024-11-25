package main

import (
	"github.com/clevrf0x/go-multitenancy/api"
	"github.com/clevrf0x/go-multitenancy/db"
	"github.com/clevrf0x/go-multitenancy/server"
)

func main() {
	db := db.New()
	server := server.NewServer(db)
	api.RegisterMiddlewares(server.Router)
	api.RegisterRoutes(server)
	server.Start()
}
