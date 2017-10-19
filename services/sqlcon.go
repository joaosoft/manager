package mgr

import (
	"database/sql"
	"fmt"

	"github.com/joaosoft/go-manager/services/sqlcon"
	"github.com/labstack/gommon/log"
)

// -------------- SQL POOLS --------------
// NewSQLPool ... creates a new sql connection pool
func (instance *Manager) NewSQLConnection(config *sqlcon.Config) (*sqlcon.SQLConController, error) {
	log.Infof(fmt.Sprintf("sqlcon, connection created"))
	return sqlcon.NewSQLConnection(config)
}

// -------------- METHODS --------------
// GetConnection ... get a sql connection with key
func (instance *Manager) GetConnection(key string) (*sql.DB, error) {
	connection, err := instance.sqlConController[key].GetConnection()
	return connection, err
}

// AddConnection ... add a connection with key
func (instance *Manager) AddConnection(key string, sqlConController *sqlcon.SQLConController) error {
	instance.sqlConController[key] = sqlConController
	log.Infof(fmt.Sprintf("sqlcon, connection '%s' added", key))

	return nil
}

// RemConnection ... remove the connection by bey
func (instance *Manager) RemConnection(key string) (*sql.DB, error) {
	// get connection
	controller := instance.sqlConController[key]

	// delete connection
	delete(instance.sqlConController, key)
	log.Infof(fmt.Sprintf("sqlcon, connection '%s' removed", key))

	return controller.Connection, nil
}
