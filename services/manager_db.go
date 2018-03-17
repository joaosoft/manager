package gomanager

import "database/sql"

// DB ...
type DB struct {
	*sql.DB
}

// AddWeb ...
func (manager *GoManager) AddDB(key string, db *DB) error {
	manager.dbs[key] = db
	log.Infof("database %s added", key)

	return nil
}

// RemoveWeb ...
func (manager *GoManager) RemoveDB(key string) (*DB, error) {
	db := manager.dbs[key]

	delete(manager.dbs, key)
	log.Infof("database %s removed", key)

	return db, nil
}

// GetDB ...
func (manager *GoManager) GetDB(key string) *DB {
	if db, exists := manager.dbs[key]; exists {
		return db
	}
	log.Infof("database %s doesn't exist", key)
	return nil
}
