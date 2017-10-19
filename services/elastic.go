package mgr

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/joaosoft/go-Manager/elastic"
)

// -------------- ELASTIC --------------
// NewElasticClient ... creates a new elastic client
func (instance *Manager) NewElasticClient(config *elastic.Config) *elastic.elasticController {
	log.Infof(fmt.Sprintf("elastic, creating elastic"))
	return elastic.NewElastic(config)
}

// -------------- METHODS --------------
// GetElasticClient ... get elastic client by key
func (instance *Manager) GetElasticClient(key string) (*elastic.elasticController, error) {
	return instance.elasticController[key], nil
}

// AddElasticClient, add a new elastic client
func (instance *Manager) AddElasticClient(key string, elasticClient *elastic.elasticController) error {
	log.Infof(fmt.Sprintf("elastic, add a new elastic client '%s'", key))
	instance.elasticController[key] = elasticClient
	return nil
}

// RemElasticClient, remove the elastic client by key
func (instance *Manager) RemElasticClient(key string) (*elastic.elasticController, error) {
	log.Infof(fmt.Sprintf("elastic, remove the elastic client '%s'", key))

	// get gateway
	controller := instance.elasticController[key]

	// delete gateway
	delete(instance.elasticController, key)
	log.Infof(fmt.Sprintf("elastic, elastic client '%s' removed", key))

	return controller, nil
}
