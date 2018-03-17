package main

import (
	"io"
	"net/http"

	"go-manager/services"

	"encoding/json"

	"time"

	"os"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	nsqlib "github.com/nsqio/go-nsq"
)

// --------- dummy process ---------
type dummy_process struct{}

func (manager *dummy_process) Start() error {
	return nil
}

func (manager *dummy_process) Stop() error {
	return nil
}

func (manager *dummy_process) Started() bool {
	return true
}

// --------- dummy nsq ---------
type dummy_nsq_handler struct{}

func (dummy *dummy_nsq_handler) HandleMessage(msg *nsqlib.Message) error {
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

func dummy_web_echo_handler(ctx echo.Context) error {
	type Example struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	return ctx.JSON(http.StatusOK, Example{Id: ctx.Param("id"), Name: "joao", Age: 29})
}

func main() {
	//
	// manager
	manager, _ := gomanager.NewManager()

	//
	// processes
	_ = manager.AddProcess("process_1", &dummy_process{})

	//
	// nsq producer
	nsqConfigConsumer := gomanager.NewNSQConfig("topic_1", "channel_1", []string{"127.0.0.1:4161"}, []string{"127.0.0.1:4150"})
	nsqProducer, _ := gomanager.NewSimpleNSQProducer(nsqConfigConsumer)
	manager.AddNSQProducer("nsq_producer_1", nsqProducer)
	nsqProducer = manager.GetNSQProducer("nsq_producer_1")
	nsqProducer.Publish("topic_1", []byte("MENSAGEM ENVIADA PARA A NSQ"), 3)

	log.Info("waiting 1 seconds...")
	<-time.After(time.Duration(1) * time.Second)

	//
	// nsq consumer
	nsqConfigConsumer := gomanager.NewNSQConfig("topic_1", "channel_1", []string{"127.0.0.1:4161"}, []string{"127.0.0.1:4150"})
	nsqConsumer, _ := gomanager.NewSimpleNSQConsumer(nsqConfigConsumer, &dummy_nsq_handler{})
	manager.AddProcess("nsq_consumer_1", nsqConsumer)

	//
	// configuration
	type DummyConfig struct {
		App  string `json:"app"`
		User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		} `json:"user"`
	}
	dir, _ := os.Getwd()
	obj := &DummyConfig{}
	simpleConfig, _ := gomanager.NewSimpleConfig(dir+"/bin/launcher/data/config.json", obj)
	manager.AddConfig("config_1", simpleConfig)
	config := manager.GetConfig("config_1")

	jsonIndent, _ := json.MarshalIndent(config.Get(), "", "    ")
	log.Infof("CONFIGURATION: %s", jsonIndent)

	//
	// web - with http
	web := gomanager.NewSimpleWebHttp(":8081")
	if err := manager.AddWeb("web_http", web); err != nil {
		log.Error("error adding web process to manager")
	}
	web = manager.GetWeb("web_http")
	web.AddRoute(http.MethodGet, "/web_http", dummy_web_http_handler)

	//
	// web - with echo
	web = gomanager.NewSimpleWebEcho(":8082")
	if err := manager.AddWeb("web_echo", web); err != nil {
		log.Error("error adding web process to manager")
	}
	web = manager.GetWeb("web_echo")
	web.AddRoute(http.MethodGet, "/web_echo/:id", dummy_web_echo_handler)
	go web.Start() // starting this because of the gateway

	log.Info("waiting 1 seconds...")
	<-time.After(time.Duration(1) * time.Second)

	//
	// gateway
	headers := map[string][]string{"Content-Type": {"application/json"}}
	var body io.Reader

	gateway := gomanager.NewSimpleGateway()
	manager.AddGateway("gateway_1", gateway)
	gateway = manager.GetGateway("gateway_1")
	status, bytes, err := gateway.Request(http.MethodGet, "http://127.0.0.1:8082", "/web_echo/123", headers, body)
	log.Infof("status: %d, response: %s, error? %t", status, string(bytes), err != nil)

	manager.Start()
}
