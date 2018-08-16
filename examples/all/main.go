package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"manager"

	"sync"

	"github.com/joaosoft/logger"
	"github.com/labstack/echo"
	"github.com/nsqio/go-nsq"
	"github.com/streadway/amqp"
)

var log = logger.NewLogDefault("manager", logger.InfoLevel)

func dummy_process() error {
	log.Info("hello, i'm executing the dummy process")
	return nil
}

// --------- dummy nsq ---------
type dummy_nsq_handler struct{}

func (dummy *dummy_nsq_handler) HandleMessage(msg *nsq.Message) error {
	log.Infof("executing the handle message of NSQ with [ message: %s ]", string(msg.Body))
	return nil
}

// --------- dummy web http ---------
func dummy_web_http_handler(w http.ResponseWriter, r *http.Request) {
	type Example struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	example := Example{Id: "123", Name: "joao", Age: 29}
	jsonIndent, _ := json.MarshalIndent(example, "", "    ")
	w.Write(jsonIndent)
}

// --------- dummy web echo ---------
func dummy_web_echo_handler(ctx echo.Context) error {
	type Example struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	return ctx.JSON(http.StatusOK, Example{Id: ctx.Param("id"), Name: "joao", Age: 29})
}

func work_handler(id string, data interface{}) error {
	log.Infof("work with the id %s and data %s done!", id, data.(string))
	return nil
}

func rabbit_consumer_handler(message amqp.Delivery) error {
	log.Errorf("\nA IMPRIMIR MENSAGEM %s", string(message.Body))
	return nil
}

func main() {
	//
	// manager
	m := manager.NewManager()

	//
	// manager: processes
	process := manager.NewSimpleProcess(dummy_process)
	if err := m.AddProcess("process_1", process); err != nil {
		log.Errorf("MAIN: error on processes %s", err)
	}

	//
	// nsq rabbitmqProducer
	nsqConfigProducer := manager.NewNSQConfig("topic_1", "channel_1", []string{"127.0.0.1:4150"}, []string{"127.0.0.1:4161"}, 30, 5)
	nsqProducer, _ := manager.NewSimpleNSQProducer(nsqConfigProducer)
	m.AddNSQProducer("nsq_producer_1", nsqProducer)
	nsqProducer = m.GetNSQProducer("nsq_producer_1")
	nsqProducer.Publish("topic_1", []byte("MENSAGEM ENVIADA PARA A NSQ"), 3)

	log.Info("waiting 1 seconds...")
	<-time.After(time.Duration(1) * time.Second)

	//
	// manager: nsq rabbitmqconsumer
	nsqConfigConsumer := manager.NewNSQConfig("topic_1", "channel_1", []string{"127.0.0.1:4161"}, []string{"127.0.0.1:4150"}, 30, 5)
	nsqConsumer, _ := manager.NewSimpleNSQConsumer(nsqConfigConsumer, &dummy_nsq_handler{})
	m.AddProcess("nsq_consumer_1", nsqConsumer)

	//
	// manager: configuration
	type dummy_config struct {
		App  string `json:"app"`
		User struct {
			Name   string `json:"name"`
			Age    int    `json:"age"`
			Random int    `json:"random"`
		} `json:"user"`
	}
	dir, _ := os.Getwd()
	obj := &dummy_config{}
	simpleConfig, _ := manager.NewSimpleConfig(dir+"/examples/all/data/config.json", obj)
	m.AddConfig("config_1", simpleConfig)
	config := m.GetConfig("config_1")

	jsonIndent, _ := json.MarshalIndent(config.GetObj(), "", "    ")
	log.Infof("CONFIGURATION: %s", jsonIndent)

	// allows to set a new configuration and save in the file
	n := rand.Intn(9000)
	obj.User.Random = n
	log.Infof("MAIN: Random: %d", n)
	config.Set(obj)
	if err := config.Save(); err != nil {
		log.Error("MAIN: error whe saving configuration file")
	}

	//
	// manager: web

	// web - with http
	web := manager.NewSimpleWebHttp(":8081")
	if err := m.AddWeb("web_http", web); err != nil {
		log.Error("error adding web process to manager")
	}
	web = m.GetWeb("web_http")
	web.AddRoute(http.MethodGet, "/web_http", dummy_web_http_handler)

	// web - with echo
	web = manager.NewSimpleWebEcho(":8082")
	if err := m.AddWeb("web_echo", web); err != nil {
		log.Error("error adding web process to manager")
	}
	web = m.GetWeb("web_echo")
	web.AddRoute(http.MethodGet, "/web_echo/:id", dummy_web_echo_handler)
	go web.Start(&sync.WaitGroup{}) // starting this because of the gateway

	log.Info("waiting 1 seconds...")
	<-time.After(time.Duration(1) * time.Second)

	//
	// manager: gateway
	headers := map[string][]string{"Content-Type": {"application/json"}}

	gateway := manager.NewSimpleGateway()
	m.AddGateway("gateway_1", gateway)
	gateway = m.GetGateway("gateway_1")
	status, bytes, err := gateway.Request(http.MethodGet, "http://127.0.0.1:8082", "/web_echo/123", headers, nil)
	log.Infof("status: %d, response: %s, error? %t", status, string(bytes), err != nil)

	//
	// manager: database

	// database - postgres
	postgresConfig := manager.NewDBConfig("postgres", "postgres://user:password@localhost:7001?sslmode=disable")
	postgresConn := manager.NewSimpleDB(postgresConfig)
	m.AddDB("postgres", postgresConn)

	// database - mysql
	mysqlConfig := manager.NewDBConfig("mysql", "root:password@tcp(127.0.0.1:7002)/mysql")
	mysqlConn := manager.NewSimpleDB(mysqlConfig)
	m.AddDB("mysql", mysqlConn)

	//
	// manager: redis
	redisConfig := manager.NewRedisConfig("127.0.0.1", 7100, 0, "")
	redisConn := manager.NewSimpleRedis(redisConfig)
	m.AddRedis("redis", redisConn)

	//
	// manager: workqueue
	workqueueConfig := manager.NewWorkListConfig("queue_001", 1, 2, time.Second*2, manager.FIFO)
	workqueue := manager.NewSimpleWorkList(workqueueConfig, work_handler)
	m.AddWorkList("queue_001", workqueue)
	workqueue = m.GetWorkList("queue_001")
	for i := 1; i <= 1000; i++ {
		workqueue.AddWork(fmt.Sprintf("PROCESS: %d", i), fmt.Sprintf("THIS IS MY MESSAGE %d", i))
	}
	wgQueue := sync.WaitGroup{}
	wgQueue.Add(1)
	if err := workqueue.Start(&wgQueue); err != nil {
		log.Errorf("MAIN: error on workqueue %s", err)
	}

	//
	// manager: rabbitmq rabbitmqProducer
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s%s", "root", "password", "localhost", "5673", "/local")
	exchange := "example"
	exchangeType := "direct"
	queue := "test-queue"
	bindingKey := "test-key"
	consumerTag := "simple-rabbitmqconsumer"
	configRabbitmq := manager.NewRabbitmqConfig(uri, exchange, exchangeType)

	rabbitmqProducer, err := manager.NewRabbitmqProducer(configRabbitmq)
	if err != nil {
		log.Errorf("%s", err)
	}

	wgRabitMqProducer := sync.WaitGroup{}
	wgRabitMqProducer.Add(1)
	if err := rabbitmqProducer.Start(&wgRabitMqProducer); err != nil {
		log.Errorf("%s", err)
	}
	m.AddRabbitmqProducer("rabbitmq_producer", rabbitmqProducer)

	err = rabbitmqProducer.Publish(bindingKey, []byte(`teste do joao`), true)
	if err != nil {
		log.Errorf("%s", err)
	}

	//
	// manager: rabbitmq rabbitmqconsumer
	rabbitmqconsumer, err := manager.NewRabbitmqConsumer(configRabbitmq, queue, bindingKey, consumerTag, rabbit_consumer_handler)
	if err != nil {
		log.Errorf("%s", err)
	}
	m.AddRabbitmqConsumer("rabbitmq_consumer", rabbitmqconsumer)

	m.Start()
}
