package manager

import (
	"fmt"
	"github.com/joaosoft/logger"
	"time"

	"sync"

	"github.com/nsqio/go-nsq"
)

// SimpleNSQConsumer ...
type SimpleNSQConsumer struct {
	client  *nsq.Consumer
	handler INSQHandler
	logger logger.ILogger
	config  *NSQConfig
	started bool
}

// NewSimpleNSQConsumer ...
func (manager *Manager) NewSimpleNSQConsumer(config *NSQConfig, handler INSQHandler) (INSQConsumer, error) {
	manager.logger.Infof("nsq consumer, creating consumer [ topic: %s, channel: %s ]", config.Topic, config.Channel)

	// Creating nsq configuration
	nsqConfig := nsq.NewConfig()
	nsqConfig.MaxAttempts = config.MaxAttempts
	nsqConfig.DefaultRequeueDelay = time.Duration(config.RequeueDelay) * time.Second
	nsqConfig.MaxInFlight = config.MaxInFlight
	nsqConfig.ReadTimeout = 120 * time.Second

	nsqConsumer, err := nsq.NewConsumer(config.Topic, config.Channel, nsqConfig)
	nsqConsumer.AddHandler(handler)
	if err != nil {
		panic(err)
	}

	consumer := &SimpleNSQConsumer{
		client:  nsqConsumer,
		config:  config,
		handler: handler,
	}

	manager.logger.Infof("nsq consumer, consumer [ topic: %s, channel: %s ] created", config.Topic, config.Channel)

	return consumer, nil
}

// HandleMessage ...
func (consumer *SimpleNSQConsumer) HandleMessage(message *nsq.Message) error {
	message.DisableAutoResponse()

	if err := consumer.handler.HandleMessage(message); err != nil {
		return err
	}

	return nil
}

// Stop ...
func (consumer *SimpleNSQConsumer) Started() bool {
	return consumer.started
}

// Start ...
func (consumer *SimpleNSQConsumer) Start(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if consumer.started {
		return nil
	}

	if consumer.handler == nil {
		return fmt.Errorf("nsq consumer, no handler configured")
	}

	if consumer.config.Lookupd != nil && len(consumer.config.Lookupd) > 0 {
		for _, addr := range consumer.config.Lookupd {
			consumer.logger.Infof("nsq consumer, consumer connecting to %s", addr)
		}
		if err := consumer.client.ConnectToNSQLookupds(consumer.config.Lookupd); err != nil {
			consumer.logger.Infof("nsq consumer, error connecting to loookupd %s", consumer.config.Nsqd)
			consumer.logger.Error(err)
			return err
		}
	}
	if consumer.config.Nsqd != nil && len(consumer.config.Nsqd) > 0 {
		for _, addr := range consumer.config.Nsqd {
			consumer.logger.Infof("nsq consumer, connecting to %s", addr)
		}
		if err := consumer.client.ConnectToNSQDs(consumer.config.Nsqd); err != nil {
			consumer.logger.Infof("nsq consumer, error connecting to nsqd %s", consumer.config.Nsqd)
			return err
		}
	}

	if !consumer.started {
		panic("nsq consumer, failed to start consumer")
	}

	consumer.started = true

	<-consumer.client.StopChan

	return nil
}

// Stop ...
func (consumer *SimpleNSQConsumer) Stop(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if !consumer.started {
		return nil
	}

	consumer.client.Stop()
	consumer.started = false

	return nil
}
