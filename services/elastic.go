package mgr

import (
	"fmt"

	"github.com/joaosoft/go-manager/services/elastic"
	"github.com/labstack/gommon/log"
)

// -------------- ELASTIC --------------
// NewElasticClient ... creates a new elastic client
func (instance *Manager) NewElasticClient(config *elastic.Config) *elastic.ElasticController {
	log.Infof(fmt.Sprintf("elastic, creating elastic"))
	return elastic.NewElastic(config)
}

// -------------- METHODS --------------
// GetElasticClient ... get elastic client by key
func (instance *Manager) GetElasticClient(key string) (*elastic.ElasticController, error) {
	return instance.ElasticController[key], nil
}

// AddElasticClient, add a new elastic client
func (instance *Manager) AddElasticClient(key string, elasticClient *elastic.ElasticController) error {
	log.Infof(fmt.Sprintf("elastic, add a new elastic client '%s'", key))
	instance.ElasticController[key] = elasticClient
	return nil
}

// RemElasticClient, remove the elastic client by key
func (instance *Manager) RemElasticClient(key string) (*elastic.ElasticController, error) {
	log.Infof(fmt.Sprintf("elastic, remove the elastic client '%s'", key))

	// get gateway
	controller := instance.ElasticController[key]

	// delete gateway
	delete(instance.ElasticController, key)
	log.Infof(fmt.Sprintf("elastic, elastic client '%s' removed", key))

	return controller, nil
}
