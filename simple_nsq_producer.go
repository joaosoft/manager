package manager

import (
	"fmt"
	"time"

	"sync"

	"github.com/nsqio/go-nsq"
)

// Producer ...
type SimpleNSQProducer struct {
	client  *nsq.Producer
	config  *NSQConfig
	started bool
}

// NewSimpleNSQProducer ...
func NewSimpleNSQProducer(config *NSQConfig) (INSQProducer, error) {
	var addr string

	// nsq configuration
	nsqConfig := nsq.NewConfig()
	nsqConfig.MaxAttempts = config.MaxAttempts
	nsqConfig.DefaultRequeueDelay = time.Duration(config.RequeueDelay) * time.Second
	nsqConfig.MaxInFlight = config.MaxInFlight
	nsqConfig.ReadTimeout = 120 * time.Second

	if config.Lookupd != nil && len(config.Lookupd) > 0 {
		addr = config.Lookupd[0]
	} else {
		return nil, fmt.Errorf("nsq producer hasn't the address to Connect")
	}

	log.Infof("connecting nsq producer to %s", addr)
	nsqProducer, err := nsq.NewProducer(addr, nsqConfig)
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
func (producer *SimpleNSQProducer) Start(wg *sync.WaitGroup) error {
	defer wg.Done()

	if !producer.started {
		producer.started = true
	}
	return nil
}

// Stop ...
func (producer *SimpleNSQProducer) Stop(wg *sync.WaitGroup) error {
	defer wg.Done()

	if producer.started {
		producer.client.Stop()
		producer.started = false
	}
	return nil
}

// Start ...
func (producer *SimpleNSQProducer) Started() bool {
	return true
}

// Ping ...
func (producer *SimpleNSQProducer) Ping() error {
	return producer.client.Ping()
}
