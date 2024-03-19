package main

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

var log = logger.NewLogDefault("manager", logger.InfoLevel)

func main() {
	//
	// manager
	m := manager.NewManager(manager.WithRunInBackground(false))

	postgresConfig := manager.NewDBConfig("postgres", "postgres://postgres:postgres@localhost:5432?sslmode=disable")
	postgresConn := m.NewSimpleDB(postgresConfig)
	m.AddDB("postgres", postgresConn)
	m.Start()
}
