package go_manager

import (
	"database/sql"
	"fmt"

	"go-manager/services/sqlcon"
	"github.com/labstack/gommon/log"
)

// -------------- SQL POOLS --------------
// NewSQLPool ... creates a new sql connection pool
func (manager *Manager) NewSQLConnection(config *sqlcon.Config) (*sqlcon.SQLConController, error) {
	log.Infof(fmt.Sprintf("sqlcon, connection created"))
	return sqlcon.NewSQLConnection(config)
}

// -------------- METHODS --------------
// GetConnection ... get a sql connection with key
func (manager *Manager) GetConnection(key string) (*sql.DB, error) {
	connection, err := manager.SqlConController[key].GetConnection()
	return connection, err
}

// AddConnection ... add a connection with key
func (manager *Manager) AddConnection(key string, SqlConController *sqlcon.SQLConController) error {
	manager.SqlConController[key] = SqlConController
	log.Infof(fmt.Sprintf("sqlcon, connection '%s' added", key))

	return nil
}

// RemConnection ... remove the connection by bey
func (manager *Manager) RemConnection(key string) (*sql.DB, error) {
	// get connection
	controller := manager.SqlConController[key]

	// delete connection
	delete(manager.SqlConController, key)
	log.Infof(fmt.Sprintf("sqlcon, connection '%s' removed", key))

	return controller.Connection, nil
}
