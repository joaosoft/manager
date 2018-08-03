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
	key        string
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
		key:        bindingKey,
		tag:        tag,
		handler:    handler,
		done:       make(chan error),
	}

	var err error

	log.Infof("dialing %s", config.Uri)
	consumer.connection, err = amqp.Dial(config.Uri)
	if err != nil {
		return nil, fmt.Errorf("dial: %s", err)
	}

	log.Infof("got connection, getting channel")
	consumer.channel, err = consumer.connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %s", err)
	}

	log.Infof("got channel, declaring exchange (%s)", config.Exchange)
	if err = consumer.channel.ExchangeDeclare(
		config.Exchange,     // name of the exchange
		config.ExchangeType, // type
		true,                // durable
		false,               // delete when complete
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	); err != nil {
		return nil, fmt.Errorf("exchange declare: %s", err)
	}

	log.Infof("declared exchange, declaring queue (%s)", queue)
	state, err := consumer.channel.QueueDeclare(
		queue, // name of the queue
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue declare: %s", err)
	}

	log.Infof("declared queue (%d messages, %d consumers), binding to exchange (bindingKey '%s')", state.Messages, state.Consumers, bindingKey)

	if err = consumer.channel.QueueBind(
		queue,           // name of the queue
		bindingKey,      // bindingKey
		config.Exchange, // sourceExchange
		false,           // noWait
		nil,             // arguments
	); err != nil {
		return nil, fmt.Errorf("queue bind: %s", err)
	}

	return consumer, nil
}

func (consumer *RabbitmqConsumer) Start() error {
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
