package manager

import (
	"fmt"
	"time"

	"sync"

	"github.com/nsqio/go-nsq"
)

// SimpleNSQConsumer ...
type SimpleNSQConsumer struct {
	client  *nsq.Consumer
	handler INSQHandler
	config  *NSQConfig
	started bool
}

// NewSimpleNSQConsumer ...
func NewSimpleNSQConsumer(config *NSQConfig, handler INSQHandler) (INSQConsumer, error) {
	log.Infof("nsq consumer, creating consumer [ topic: %s, channel: %s ]", config.Topic, config.Channel)

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

	log.Infof("nsq consumer, consumer [ topic: %s, channel: %s ] created", config.Topic, config.Channel)

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
func (consumer *SimpleNSQConsumer) Start(wg *sync.WaitGroup) error {
	if consumer.handler == nil {
		return fmt.Errorf("nsq consumer, no handler configured")
	}

	if consumer.config.Lookupd != nil && len(consumer.config.Lookupd) > 0 {
		consumer.started = true
		for _, addr := range consumer.config.Lookupd {
			log.Infof("nsq consumer, consumer connecting to %s", addr)
		}
		if err := consumer.client.ConnectToNSQLookupds(consumer.config.Lookupd); err != nil {
			log.Infof("nsq consumer, error connecting to loookupd %s", consumer.config.Nsqd)
			log.Error(err)
			return err
		}
	}
	if consumer.config.Nsqd != nil && len(consumer.config.Nsqd) > 0 {
		consumer.started = true
		for _, addr := range consumer.config.Nsqd {
			log.Infof("nsq consumer, connecting to %s", addr)
		}
		if err := consumer.client.ConnectToNSQDs(consumer.config.Nsqd); err != nil {
			log.Infof("nsq consumer, error connecting to nsqd %s", consumer.config.Nsqd)
			return err
		}
	}

	if !consumer.started {
		panic("nsq consumer, failed to start consumer")
	}

	consumer.started = true
	wg.Done()

	<-consumer.client.StopChan

	return nil
}

// Stop ...
func (consumer *SimpleNSQConsumer) Stop(wg *sync.WaitGroup) error {
	defer wg.Done()

	if consumer.started {
		consumer.client.Stop()
		consumer.started = false
	}

	return nil
}
