package main

import (
	"fmt"
	mgr "github.com/joaosoft/go-manager"
	"github.com/joaosoft/go-manager/nsq"
	"github.com/joaosoft/go-manager/sqlcon"
	"github.com/joaosoft/go-manager/web"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	nsqlib "github.com/nsqio/go-nsq"
	"net/http"
	"os"
)

// EXAMPLE PROCESS
type DummyProcess struct{}

func (instance *DummyProcess) Start() error {
	return nil
}

func (instance *DummyProcess) Stop() error {
	return nil
}

// EXAMPLE NSQ HANDLER
type DummyNSQHandler struct{}

func (instance *DummyNSQHandler) HandleMessage(msg *nsqlib.Message) error {
	return nil
}

// EXAMPLE WEB SERVER HANDLER
func exampleWebServerHandler(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	log.Info(fmt.Sprintf("Web Server requested with id '%s'", id))
	return c.String(http.StatusOK, id)
}

func main() {
	//
	// MANAGER
	//
	manager, _ := mgr.NewManager()

	//
	// PROCESSES
	//
	_ = manager.AddProcess("process_1", &DummyProcess{})

	//
	// SQL CONNECTION
	//
	sqlConfig := sqlcon.NewConfig("localhost", "postgres", 1, 2)
	sqlConnection, _ := manager.NewSQLConnection(sqlConfig)
	_ = manager.AddConnection("conn_1", sqlConnection)

	//
	// NSQ
	//
	nsqConfig := &nsq.Config{
		Topic:   "topic_1",
		Channel: "channel_2",
		Lookupd: []string{"http://localhost:4151"},
	}

	// Consumer
	nsqConsumer, _ := manager.NewNSQConsumer(nsqConfig, &DummyNSQHandler{})
	manager.AddProcess("consumer_1", nsqConsumer)

	// Producer
	nsqProducer, _ := manager.NewNSQProducer(nsqConfig)
	nsqProducer.Publish("topic_1", []byte("body"), 3)
	manager.AddProcess("producer_1", nsqProducer)

	//
	// CONFIGURATION
	//
	dir, _ := os.Getwd()
	simpleConfig, _ := manager.NewSimpleConfig(dir+"/getting_started/system/", "config", "json")
	manager.AddConfig("config_1", simpleConfig)

	// Get configuration by path
	fmt.Println("a: ", manager.GetConfig("config_1").Get("a"))
	fmt.Println("caa: ", manager.GetConfig("config_1").Get("c.ca.caa"))

	// Get configuration by tag
	fmt.Println("a: ", manager.GetConfig("config_1").Get("a"))
	fmt.Println("caa: ", manager.GetConfig("config_1").Get("c.ca.caa"))

	//
	// HTTP SERVER
	//
	configWebServer := web.NewConfig("localhost:8081")
	webServer, _ := manager.NewWEBServer(configWebServer)
	webServer.AddRoute(http.MethodGet, "/example/:id", exampleWebServerHandler)
	manager.AddProcess("web_server_1", webServer)

	manager.Start()
}
