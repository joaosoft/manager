package main

import (
	"fmt"
	mgr "github.com/joaosoft/go-manager"
	"github.com/joaosoft/go-manager/nsq"
	nsqlib "github.com/nsqio/go-nsq"
)

// EXAMPLE NSQ HANDLER
type DummyNSQHandler struct{}

func (instance *DummyNSQHandler) HandleMessage(message *nsqlib.Message) error {
	fmt.Println("THIS IS THE RECEIVED MESSAGE: ", string(message.Body))
	return nil
}

func main() {
	//
	// MANAGER
	//
	manager, _ := mgr.NewManager()

	//
	// NSQ
	//

	// Consumer
	nsqConsumerConfig := &nsq.Config{
		Topic:   "topic_1",
		Channel: "channel_2",
		Nsqd:    []string{"localhost:4150"},
	}
	nsqConsumer, _ := manager.NewNSQConsumer(nsqConsumerConfig, &DummyNSQHandler{})
	manager.AddProcess("consumer_1", nsqConsumer)

	// Producer
	nsqProducerConfig := &nsq.Config{
		Topic:   "topic_1",
		Channel: "channel_2",
		Lookupd: []string{"localhost:4150"},
	}
	nsqProducer, _ := manager.NewNSQProducer(nsqProducerConfig)
	nsqProducer.Publish("topic_1", []byte("MENSAGEM ENVIADA PARA A NSQ"), 3)
	manager.AddProcess("producer_1", nsqProducer)

	manager.Start()
}
