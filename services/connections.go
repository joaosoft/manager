package gomanager

import (
	"database/sql"
	"time"

	redis "github.com/mediocregopher/radix.v3"
	nsqlib "github.com/nsqio/go-nsq"
)

// createConnection ...
func (config *RedisConfig) connect() (*redis.Pool, error) {
	log.Infof("connecting with protocol [ %s ], address [ %s ] and size [ %d ]", config.Protocol, config.Address, config.Size)
	return redis.NewPool(config.Protocol, config.Address, config.Size, nil)
}

// createConnection ...
func (config *SqlConfig) connect() (*sql.DB, error) {
	log.Infof("connecting with driver [ %s ] and data source [ %s ]", config.Driver, config.DataSource)
	return sql.Open(config.Driver, config.DataSource)
}

// createConnection ...
func (config *NsqConfig) connect() (*nsqlib.Producer, error) {
	nsqConfig := nsqlib.NewConfig()
	nsqConfig.MaxAttempts = config.MaxAttempts
	nsqConfig.DefaultRequeueDelay = time.Duration(config.RequeueDelay) * time.Second
	nsqConfig.MaxInFlight = config.MaxInFlight
	nsqConfig.ReadTimeout = 120 * time.Second

	log.Infof("connecting with max attempts [ %d ]", config.MaxAttempts)

	return nsqlib.NewProducer(config.Lookupd, nsqConfig)
}
