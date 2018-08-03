package manager

import (
	"database/sql"
	"time"

	"fmt"

	"github.com/alphazero/Go-Redis"
	"github.com/nsqio/go-nsq"
	"github.com/streadway/amqp"
)

// Connect ...
func (config *RedisConfig) Connect() (redis.Client, error) {
	log.Infof("connecting redis with [ host: %s, port: %d ]", config.Host, config.Port)
	spec := redis.DefaultSpec().Host(config.Host).Port(config.Port).Password(config.Password).Db(config.Database)
	return redis.NewSynchClientWithSpec(spec)
}

// Connect ...
func (config *DBConfig) Connect() (*sql.DB, error) {
	log.Infof("connecting database with driver [ %s ] and data source [ %s ]", config.Driver, config.DataSource)
	return sql.Open(config.Driver, config.DataSource)
}

// Connect ...
func (config *NSQConfig) Connect() (*nsq.Producer, error) {
	nsqConfig := nsq.NewConfig()
	nsqConfig.MaxAttempts = config.MaxAttempts
	nsqConfig.DefaultRequeueDelay = time.Duration(config.RequeueDelay) * time.Second
	nsqConfig.MaxInFlight = config.MaxInFlight
	nsqConfig.ReadTimeout = 120 * time.Second

	log.Infof("connecting nsq with max attempts [ %d ]", config.MaxAttempts)
	return nsq.NewProducer(config.Lookupd[0], nsqConfig)
}

// Connect ...
func (config *RabbitmqConfig) Connect() (*amqp.Connection, error) {
	log.Infof("dialing %s", config.Uri)
	connection, err := amqp.Dial(config.Uri)
	if err != nil {
		return nil, fmt.Errorf("dial: %s", err)
	}

	return connection, nil
}
