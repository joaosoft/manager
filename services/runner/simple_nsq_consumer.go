package runner

import (
	"fmt"
	"github.com/labstack/gommon/log"
	nsqlib "github.com/nsqio/go-nsq"
	"time"
)

// IConsumer costumer interface
type IConsumer interface {
	Start() error
	Stop() error
}

// Consumer ... costumer structure
type Consumer struct {
	client  *nsqlib.Consumer
	started bool
	handler IHandler
	config  *Config
}

// NewConsumer ... creates a new costumer
func NewConsumer(config *Config, handler IHandler) (IConsumer, error) {
	log.Infof("nsq consumer, creating manager [topic:%s][channel:%s]", config.Topic, config.Channel)

	// Creating nsq configuration
	nsqConfig := nsqlib.NewConfig()
	nsqConfig.MaxAttempts = config.MaxAttempts
	nsqConfig.DefaultRequeueDelay = time.Duration(config.RequeueDelay) * time.Second
	nsqConfig.MaxInFlight = config.MaxInFlight
	nsqConfig.ReadTimeout = 120 * time.Second

	nsqConsumer, err := nsqlib.NewConsumer(config.Topic, config.Channel, nsqConfig)
	nsqConsumer.AddHandler(handler)
	if err != nil {
		panic(err)
	}

	manager := &Consumer{
		client:  nsqConsumer,
		config:  config,
		started: false,
		handler: handler,
	}

	log.Infof("nsq consumer, manager [topic:%s][channel:%s] created", config.Topic, config.Channel)

	return manager, nil
}

// HandleMessage ... handle message queue
func (manager *Consumer) HandleMessage(message *nsqlib.Message) error {
	message.DisableAutoResponse()

	if err := manager.handler.HandleMessage(message); err != nil {
		return err
	}

	return nil
}

// Start ... start's nsq costumer
func (manager *Consumer) Start() error {
	if manager.handler == nil {
		return fmt.Errorf("nsq consumer, no handler configured")
	}

	if manager.config.Lookupd != nil && len(manager.config.Lookupd) > 0 {
		manager.started = true
		for _, addr := range manager.config.Lookupd {
			log.Infof("nsq consumer, manager connecting to %s", addr)
		}
		if err := manager.client.ConnectToNSQLookupds(manager.config.Lookupd); err != nil {
			log.Infof("nsq consumer, error connecting to loookupd %s", manager.config.Nsqd)
			log.Error(err)
			return err
		}
	}
	if manager.config.Nsqd != nil && len(manager.config.Nsqd) > 0 {
		manager.started = true
		for _, addr := range manager.config.Nsqd {
			log.Infof("nsq consumer, manager connecting to %s", addr)
		}
		if err := manager.client.ConnectToNSQDs(manager.config.Nsqd); err != nil {
			log.Infof("nsq consumer, error connecting to nsqd %s", manager.config.Nsqd)
			return err
		}
	}

	if !manager.started {
		panic("nsq consumer, failed to start manager")
	}

	<-manager.client.StopChan

	manager.started = false

	return nil
}

// Stop ... stop's nsq costumer
func (manager *Consumer) Stop() error {
	log.Infof("nsq consumer, stopping ")
	manager.client.Stop()
	manager.started = false
	log.Infof("nsq consumer, stopped")

	return nil
}
