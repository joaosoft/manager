package go_manager

import (
	"database/sql"
	"fmt"
)

// SQLConController ... creates a new sql pool controller
type SQLConController struct {
	Connection *sql.DB
	Config     *Config
}

// NewSQLPool ... create a new sql pool
func NewSQLConnection(config *Config) (*SQLConController, error) {
	conn, _ := newConnection(config)

	manager := &SQLConController{
		Connection: conn,
		Config:     config,
	}

	return manager, nil
}

// GetConnection ... get connection
func (manager *SQLConController) GetConnection() (*sql.DB, error) {
	conn := manager.Connection

	if conn == nil {
		return nil, fmt.Errorf("could not get a connection")
	}

	return conn, nil
}

// AddConnection ... set connection
func (manager *SQLConController) SetConnection(config *Config) (*sql.DB, error) {
	var err error

	manager.Connection, err = newConnection(config)

	return manager.Connection, err
}

func newConnection(config *Config) (*sql.DB, error) {
	var conn *sql.DB
	var err error

	if conn, err = sql.Open(config.Driver, config.Endpoint); err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(config.MaxIdleConnections)
	conn.SetMaxOpenConns(config.MaxOpenConnections)

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
