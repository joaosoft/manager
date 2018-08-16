package main

import (
	"manager"

	"github.com/joaosoft/logger"
)

var log = logger.NewLogDefault("manager", logger.InfoLevel)

func main() {
	//
	// manager
	m := manager.NewManager(manager.WithRunInBackground(true))

	postgresConfig := manager.NewDBConfig("postgres", "postgres://postgres:postgres@localhost:5432?sslmode=disable")
	postgresConn := manager.NewSimpleDB(postgresConfig)
	m.AddDB("postgres", postgresConn)
	m.Start()
}
