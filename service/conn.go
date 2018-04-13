package gomanager

import (
	"database/sql"
	"time"

	redis "github.com/alphazero/Go-Redis"
	nsqlib "github.com/nsqio/go-nsq"
)

// createConnection ...
func (config *RedisConfig) connect() (redis.Client, error) {
	log.Infof("connecting redis with [ host: %s, port: %d ]", config.Host, config.Port)
	spec := redis.DefaultSpec().Host(config.Host).Port(config.Port).Password(config.Password).Db(config.Database)
	return redis.NewSynchClientWithSpec(spec)
}

// createConnection ...
func (config *DBConfig) connect() (*sql.DB, error) {
	log.Infof("connecting database with driver [ %s ] and data source [ %s ]", config.Driver, config.DataSource)
	return sql.Open(config.Driver, config.DataSource)
}

// createConnection ...
func (config *NSQConfig) connect() (*nsqlib.Producer, error) {
	nsqConfig := nsqlib.NewConfig()
	nsqConfig.MaxAttempts = config.MaxAttempts
	nsqConfig.DefaultRequeueDelay = time.Duration(config.RequeueDelay) * time.Second
	nsqConfig.MaxInFlight = config.MaxInFlight
	nsqConfig.ReadTimeout = 120 * time.Second

	log.Infof("connecting nsq with max attempts [ %d ]", config.MaxAttempts)
	return nsqlib.NewProducer(config.Lookupd[0], nsqConfig)
}
