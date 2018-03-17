package gomanager

import (
	"fmt"
	"time"

	nsqlib "github.com/nsqio/go-nsq"
)

// Producer ...
type SimpleNSQProducer struct {
	client *nsqlib.Producer
	config *NSQConfig
}

// NewSimpleNSQProducer ...
func NewSimpleNSQProducer(config *NSQConfig) (INSQProducer, error) {
	var addr string

	// nsq configuration
	nsqConfig := nsqlib.NewConfig()
	nsqConfig.MaxAttempts = config.MaxAttempts
	nsqConfig.DefaultRequeueDelay = time.Duration(config.RequeueDelay) * time.Second
	nsqConfig.MaxInFlight = config.MaxInFlight
	nsqConfig.ReadTimeout = 120 * time.Second

	if config.Lookupd != nil && len(config.Lookupd) > 0 {
		addr = config.Lookupd[0]
	} else {
		return nil, fmt.Errorf("nsq producer hasn't the address to connect")
	}

	log.Infof("connecting nsq producer to %s", addr)
	nsqProducer, err := nsqlib.NewProducer(addr, nsqConfig)
	if err != nil {
		panic(err)
	}

	producer := &SimpleNSQProducer{
		client: nsqProducer,
		config: config,
	}

	return producer, nil
}

// Publish ...
func (producer *SimpleNSQProducer) Publish(topic string, body []byte, maxRetries int) error {
	var err error

	for count := 0; count < maxRetries; count++ {
		if err = producer.client.Publish(topic, body); err == nil {
			return nil
		}
	}

	return err
}

// Start ...
func (producer *SimpleNSQProducer) Start() error {
	return nil
}

// Stop ...
func (producer *SimpleNSQProducer) Stop() error {
	log.Infof("producer producer, stopping")
	producer.client.Stop()
	log.Infof("producer producer, stopped")
	return nil
}

// Ping ...
func (producer *SimpleNSQProducer) Ping() error {
	return producer.client.Ping()
}
