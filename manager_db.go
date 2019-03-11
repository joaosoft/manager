package manager

import (
	"database/sql"
	"sync"
)

type IDB interface {
	Get() *sql.DB
	Start(waitGroup ...*sync.WaitGroup) error
	Stop(waitGroup ...*sync.WaitGroup) error
	Started() bool
}

// DBConfig ...
type DBConfig struct {
	Driver     string `json:"driver"`
	DataSource string `json:"datasource"`
}

// NewDBConfig...
func NewDBConfig(driver, datasource string) *DBConfig {
	return &DBConfig{
		Driver:     driver,
		DataSource: datasource,
	}
}

// AddDB ...
func (manager *Manager) AddDB(key string, db IDB) error {
	manager.dbs[key] = db
	manager.logger.Infof("database %s added", key)

	return nil
}

// RemoveDB ...
func (manager *Manager) RemoveDB(key string) (IDB, error) {
	db := manager.dbs[key]

	delete(manager.dbs, key)
	manager.logger.Infof("database %s removed", key)

	return db, nil
}

// GetDB ...
func (manager *Manager) GetDB(key string) IDB {
	if db, exists := manager.dbs[key]; exists {
		return db
	}
	manager.logger.Infof("database %s doesn't exist", key)
	return nil
}
