package gomanager

import (
	"io"
	"math/rand"
	"net/http"
	"encoding/json"
	"time"
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/nsqio/go-nsq"
)

// --------- dummy process ---------
func dummy_process() error {
	log.Info("hello, i'm exetuting the dummy process")
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

func usage() {
	//
	// manager
	manager := NewManager()

	//
	// manager: processes
	process := NewSimpleProcess(dummy_process)
	if err := manager.AddProcess("process_1", process); err != nil {
		log.Errorf("MAIN: error on processes %s", err)
	}

	//
	// nsq producer
	nsqConfigProducer := NewNSQConfig("topic_1", "channel_1", []string{"127.0.0.1:4150"}, []string{"127.0.0.1:4161"}, 30, 5)
	nsqProducer, _ := NewSimpleNSQProducer(nsqConfigProducer)
	manager.AddNSQProducer("nsq_producer_1", nsqProducer)
	nsqProducer = manager.GetNSQProducer("nsq_producer_1")
	nsqProducer.Publish("topic_1", []byte("MENSAGEM ENVIADA PARA A NSQ"), 3)

	log.Info("waiting 1 seconds...")
	<-time.After(time.Duration(1) * time.Second)

	//
	// manager: nsq consumer
	nsqConfigConsumer := NewNSQConfig("topic_1", "channel_1", []string{"127.0.0.1:4161"}, []string{"127.0.0.1:4150"}, 30, 5)
	nsqConsumer, _ := NewSimpleNSQConsumer(nsqConfigConsumer, &dummy_nsq_handler{})
	manager.AddProcess("nsq_consumer_1", nsqConsumer)

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
	simpleConfig, _ := NewSimpleConfig(dir+"/bin/launcher/data/config.json", obj)
	manager.AddConfig("config_1", simpleConfig)
	config := manager.GetConfig("config_1")

	jsonIndent, _ := json.MarshalIndent(config.Get(), "", "    ")
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
	web := NewSimpleWebHttp(":8081")
	if err := manager.AddWeb("web_http", web); err != nil {
		log.Error("error adding web process to manager")
	}
	web = manager.GetWeb("web_http")
	web.AddRoute(http.MethodGet, "/web_http", dummy_web_http_handler)

	// web - with echo
	web = NewSimpleWebEcho(":8082")
	if err := manager.AddWeb("web_echo", web); err != nil {
		log.Error("error adding web process to manager")
	}
	web = manager.GetWeb("web_echo")
	web.AddRoute(http.MethodGet, "/web_echo/:id", dummy_web_echo_handler)
	go web.Start() // starting this because of the gateway

	log.Info("waiting 1 seconds...")
	<-time.After(time.Duration(1) * time.Second)

	//
	// manager: gateway
	headers := map[string][]string{"Content-Type": {"application/json"}}
	var body io.Reader

	gateway := NewSimpleGateway()
	manager.AddGateway("gateway_1", gateway)
	gateway = manager.GetGateway("gateway_1")
	status, bytes, err := gateway.Request(http.MethodGet, "http://127.0.0.1:8082", "/web_echo/123", headers, body)
	log.Infof("status: %d, response: %s, error? %t", status, string(bytes), err != nil)

	//
	// manager: database

	// database - postgres
	postgresConfig := NewDBConfig("postgres", "postgres://user:password@localhost:7001?sslmode=disable")
	postgresConn := NewSimpleDB(postgresConfig)
	manager.AddDB("postgres", postgresConn)

	// database - mysql
	mysqlConfig := NewDBConfig("mysql", "root:password@tcp(127.0.0.1:7002)/mysql")
	mysqlConn := NewSimpleDB(mysqlConfig)
	manager.AddDB("mysql", mysqlConn)

	//
	// manager: redis
	redisConfig := NewRedisConfig("127.0.0.1", 7100, 0, "")
	redisConn := NewSimpleRedis(redisConfig)
	manager.AddRedis("redis", redisConn)

	//
	// manager: workqueue
	workqueueConfig := NewWorkListConfig("queue_001", 1, 2, time.Second*2, FIFO)
	workqueue := NewSimpleWorkList(workqueueConfig, work_handler)
	manager.AddWorkList("queue_001", workqueue)
	workqueue = manager.GetWorkList("queue_001")
	for i := 1; i <= 1000; i++ {
		go workqueue.AddWork(fmt.Sprintf("PROCESS: %d", i), fmt.Sprintf("THIS IS MY MESSAGE %d", i))
	}
	if err := workqueue.Start(); err != nil {
		log.Errorf("MAIN: error on workqueue %s", err)
	}

	manager.Start()
}
