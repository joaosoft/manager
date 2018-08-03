package manager

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitmqConsumer struct {
	config     *RabbitmqConfig
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      string
	bindingKey string
	tag        string
	handler    RabbitmqHandler
	done       chan error
	started    bool
}

func NewRabbitmqConsumer(config *RabbitmqConfig, queue, bindingKey, tag string, handler RabbitmqHandler) (*RabbitmqConsumer, error) {
	consumer := &RabbitmqConsumer{
		config:     config,
		connection: nil,
		channel:    nil,
		queue:      queue,
		bindingKey: bindingKey,
		tag:        tag,
		handler:    handler,
		done:       make(chan error),
	}

	return consumer, nil
}

func (consumer *RabbitmqConsumer) Start() error {
	var err error

	consumer.connection, err = consumer.config.Connect()
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	defer func(err error) {
		if err != nil && consumer.connection != nil {
			consumer.connection.Close()
		}
	}(err)

	log.Infof("got connection, getting channel")
	consumer.channel, err = consumer.connection.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	log.Infof("got channel, declaring exchange (%s)", consumer.config.Exchange)
	if err = consumer.channel.ExchangeDeclare(
		consumer.config.Exchange,     // name of the exchange
		consumer.config.ExchangeType, // type
		true,  // durable
		false, // delete when complete
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		consumer.connection.Close()
		return fmt.Errorf("exchange declare: %s", err)
	}

	log.Infof("declared exchange, declaring queue (%s)", consumer.queue)
	state, err := consumer.channel.QueueDeclare(
		consumer.queue, // name of the queue
		true,           // durable
		false,          // delete when usused
		false,          // exclusive
		false,          // noWait
		nil,            // arguments
	)
	if err != nil {
		consumer.connection.Close()
		return fmt.Errorf("queue declare: %s", err)
	}

	log.Infof("declared queue (%d messages, %d consumers), binding to exchange (bindingKey '%s')", state.Messages, state.Consumers, consumer.bindingKey)

	if err = consumer.channel.QueueBind(
		consumer.queue,           // name of the queue
		consumer.bindingKey,      // bindingKey
		consumer.config.Exchange, // sourceExchange
		false, // noWait
		nil,   // arguments
	); err != nil {
		consumer.connection.Close()
		return fmt.Errorf("queue bind: %s", err)
	}

	log.Infof("queue bound to exchange, starting consume (consumer tag '%s')", consumer.tag)
	deliveries, err := consumer.channel.Consume(
		consumer.queue, // name
		consumer.tag,   // consumerTag,
		false,          // noAck
		false,          // exclusive
		false,          // noLocal
		false,          // noWait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("queue consume: %s", err)
	}

	go consumer.handle(deliveries, consumer.done)

	consumer.started = true

	return nil
}

func (consumer *RabbitmqConsumer) Started() bool {
	return consumer.started
}

func (consumer *RabbitmqConsumer) Stop() error {
	// will close() the deliveries channel
	if err := consumer.channel.Cancel(consumer.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := consumer.connection.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Infof("AMQP shutdown OK")

	consumer.started = false

	// wait for handle() to exit
	return <-consumer.done
}

func (consumer *RabbitmqConsumer) handle(deliveries <-chan amqp.Delivery, done chan error) {
	for delivery := range deliveries {
		if err := consumer.handler(delivery); err != nil {
			delivery.Ack(false)
		} else {
			delivery.Ack(false)
		}
		log.Infof("got %dB delivery: [%v] %s", len(delivery.Body), delivery.DeliveryTag, delivery.Body)
	}

	log.Infof("handle: deliveries channel closed")
	done <- nil
}
