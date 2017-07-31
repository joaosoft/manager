# go-manager
[![Build Status](https://travis-ci.org/joaosoft/go-manager.svg?branch=master)](https://travis-ci.org/golang/dep) | [![Code Climate](https://codeclimate.com/github/joaosoft/go-manager)]
A package framework with application support. 
###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## with support to
* Processes
* Configurations
* NSQ Consumers
* NSQ Producers
* SQL Connections
* Web Servers
* Gateways
* Elastic Search

## Dependecy Management (using dep)

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`

## Usage 

```go
manager, _ := pm.NewManager()
```

>### Processes
```go
// EXAMPLE PROCESS
type DummyProcess struct{}

func (instance *DummyProcess) Start() error {
	return nil
}

func (instance *DummyProcess) Stop() error {
	return nil
}

_ = manager.AddProcess("process_1", &DummyProcess{})
```

>### Configurations
```go
dir, _ := os.Getwd()
simpleConfig, _ := manager.NewSimpleConfig(dir+"/getting_started/system/", "config", "json")
manager.AddConfig("config_1", simpleConfig)

// Get configuration by path
fmt.Println("a: ", manager.GetConfig("config_1").Get("a"))
fmt.Println("caa: ", manager.GetConfig("config_1").Get("c.ca.caa"))

// Get configuration by tag
fmt.Println("a: ", manager.GetConfig("config_1").Get("a"))
fmt.Println("caa: ", manager.GetConfig("config_1").Get("c.ca.caa"))
```

>### NSQ Consumers 
```go
// EXAMPLE NSQ HANDLER
type DummyNSQHandler struct{}

func (instance *DummyNSQHandler) HandleMessage(msg *nsqlib.Message) error {
	return nil
}

nsqConfig := &nsq.Config{
    Topic:   "topic_1",
    Channel: "channel_2",
    Lookupd: []string{"http://localhost:4151"},
}

// Consumer
nsqConsumer, _ := manager.NewNSQConsumer(nsqConfig, &DummyNSQHandler{})
manager.AddProcess("consumer_1", nsqConsumer)
```

>### NSQ Producers
```go
// Producer
nsqProducer, _ := manager.NewNSQProducer(nsqConfig)
manager.AddProcess("producer_1", nsqProducer)
```

>### SQL Connections
```go
sqlConfig := sqlcon.NewConfig("localhost", "postgres", 1, 2)
sqlConnection, _ := manager.NewSQLConnection(sqlConfig)
_ = manager.AddConnection("conn_1", sqlConnection)
```

>### Web Servers
```go
// EXAMPLE WEB SERVER HANDLER
func exampleWebServerHandler(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	log.Info(fmt.Sprintf("Web Server requested with id '%s'", id))
	return c.String(http.StatusOK, id)
}

configWebServer := web.NewConfig("localhost:8081")
	webServer, _ := manager.NewWEBServer(configWebServer)
	webServer.AddRoute(http.MethodGet, "/example/:id", exampleWebServerHandler)
	manager.AddProcess("web_server_1", webServer)
```

>### Gateways
```go
var headers map[string]string
var body io.Reader
configGateway := gateway.NewConfig("http://localhost:8081")
gateway, _ := manager.NewGateway(configGateway)
manager.AddGateway("gateway_1", gateway)
manager.GetGateway("gateway_1")
status, bytes, err := manager.RequestGateway("gateway_1", http.MethodGet, "/example/123456789", headers, body)
fmt.Println("STATUS:", status, "RESPONSE:", string(bytes), "err:", err)
```

>### Elastic Search
```go
configElasticClient := elastic.NewConfig("http://localhost:9200")
elasticClient := manager.NewElasticClient(configElasticClient)
manager.AddElasticClient("elastic_1", elasticClient)
response, err := elasticClient.Search("index", "type", "body")
fmt.Println("RESPONSE:", response, "ERROR:", err)
```

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
