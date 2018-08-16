package manager

import (
	"database/sql"

	"sync"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	_ "github.com/lib/pq"              // postgres driver
)

// SimpleDB ...
type SimpleDB struct {
	*sql.DB
	config  *DBConfig
	started bool
}

// NewSimpleDB ...
func NewSimpleDB(config *DBConfig) IDB {
	return &SimpleDB{
		config: config,
	}
}

// Get ...
func (db *SimpleDB) Get() *sql.DB {
	return db.DB
}

// Start ...
func (db *SimpleDB) Start(wg *sync.WaitGroup) error {
	defer wg.Done()

	if !db.started {
		if conn, err := db.config.Connect(); err != nil {
			return err
		} else {
			db.DB = conn
		}
		db.started = true
	}

	return nil
}

// Stop ...
func (db *SimpleDB) Stop(wg *sync.WaitGroup) error {
	defer wg.Done()

	if db.started {
		if err := db.Close(); err != nil {
			return err
		}
		db.started = false
	}

	return nil
}

// Started ...
func (db *SimpleDB) Started() bool {
	return db.started
}
