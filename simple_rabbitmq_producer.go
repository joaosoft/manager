package manager

import (
	"github.com/joaosoft/logger"
	"time"

	"sync"

	"github.com/streadway/amqp"
)

type SimpleRabbitmqProducer struct {
	config     *RabbitmqConfig
	connection *amqp.Connection
	channel    *amqp.Channel
	tag        string
	logger logger.ILogger
	started    bool
}

func (manager *Manager) NewSimpleRabbitmqProducer(config *RabbitmqConfig) (*SimpleRabbitmqProducer, error) {
	return &SimpleRabbitmqProducer{
		config: config,
		logger: manager.logger,
	}, nil
}

func (producer *SimpleRabbitmqProducer) Start(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if producer.started {
		return nil
	}

	var err error
	producer.connection, err = producer.config.Connect()
	if err != nil {
		err = producer.logger.Errorf("dial: %s", err).ToError()
		return err
	}

	defer func(err error) {
		if err != nil {
			if producer.connection != nil {
				producer.connection.Close()
			}
		}
	}(err)

	producer.logger.Infof("got connection, getting channel")
	if producer.channel, err = producer.connection.Channel(); err != nil {
		err = producer.logger.Errorf("channel: %s", err).ToError()
		return err
	}

	producer.logger.Infof("got channel, declaring %q exchange (%s)", producer.config.ExchangeType, producer.config.Exchange)
	if err = producer.channel.ExchangeDeclare(
		producer.config.Exchange,     // name
		producer.config.ExchangeType, // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		err = producer.logger.Errorf("exchange declare: %s", err).ToError()
		return err
	}

	producer.started = true

	return nil
}

func (producer *SimpleRabbitmqProducer) Started() bool {
	producer.started = false
	return producer.started
}

func (producer *SimpleRabbitmqProducer) Stop(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if !producer.started {
		return nil
	}

	// will close() the deliveries channel
	if err := producer.channel.Cancel(producer.tag, true); err != nil {
		err = producer.logger.Errorf("consumer cancel failed: %s", err).ToError()
		return err
	}

	if err := producer.connection.Close(); err != nil {
		err = producer.logger.Errorf("AMQP connection close error: %s", err).ToError()
		return err
	}

	producer.logger.Infof("AMQP shutdown OK")
	producer.started = false

	return nil
}

func (producer *SimpleRabbitmqProducer) Publish(routingKey string, body []byte, reliable bool) error {
	msg := amqp.Publishing{
		DeliveryMode:    amqp.Persistent,
		Timestamp:       time.Now(),
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "",
		Body:            body,
		Priority:        0, // 0-9
	}

	producer.logger.Infof("declared exchange, publishing %dB body (%s)", len(body), body)
	if err := producer.channel.Publish(
		producer.config.Exchange, // publish to an exchange
		routingKey,               // routing to 0 or more queues
		false,                    // mandatory
		false,                    // immediate
		msg,
	); err != nil {
		err = producer.logger.Errorf("exchange publish: %s", err).ToError()
		return err
	}

	return nil
}
