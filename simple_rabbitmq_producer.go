package manager

import (
	"fmt"

	"time"

	"github.com/streadway/amqp"
)

type RabbitmqProducer struct {
	config     *RabbitmqConfig
	connection *amqp.Connection
	channel    *amqp.Channel
	tag        string
	started    bool
}

func NewRabbitmqProducer(config *RabbitmqConfig) (*RabbitmqProducer, error) {
	return &RabbitmqProducer{
		config: config,
	}, nil
}

func (producer *RabbitmqProducer) Start() error {
	var err error

	producer.connection, err = producer.config.Connect()
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	defer func(err error) {
		if err != nil && producer.connection != nil {
			producer.connection.Close()
		}
	}(err)

	log.Infof("got connection, getting channel")
	producer.channel, err = producer.connection.Channel()

	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	log.Infof("got channel, declaring %q exchange (%s)", producer.config.ExchangeType, producer.config.Exchange)
	if err := producer.channel.ExchangeDeclare(
		producer.config.Exchange,     // name
		producer.config.ExchangeType, // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("exchange declare: %s", err)
	}

	producer.started = true

	return nil
}

func (producer *RabbitmqProducer) Started() bool {
	return producer.started
}

func (producer *RabbitmqProducer) Stop() error {
	// will close() the deliveries channel
	if err := producer.channel.Cancel(producer.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := producer.connection.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Infof("AMQP shutdown OK")

	return nil
}

func (producer *RabbitmqProducer) Publish(routingKey string, body []byte, reliable bool) error {
	msg := amqp.Publishing{
		DeliveryMode:    amqp.Persistent,
		Timestamp:       time.Now(),
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "",
		Body:            body,
		Priority:        0, // 0-9
	}

	log.Infof("declared exchange, publishing %dB body (%s)", len(body), body)
	if err := producer.channel.Publish(
		producer.config.Exchange, // publish to an exchange
		routingKey,               // routing to 0 or more queues
		false,                    // mandatory
		false,                    // immediate
		msg,
	); err != nil {
		return fmt.Errorf("exchange publish: %s", err)
	}

	return nil
}
