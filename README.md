# go-manager
[![Build Status](https://travis-ci.org/joaosoft/go-manager.svg?branch=master)](https://travis-ci.org/joaosoft/go-manager) | [![codecov](https://codecov.io/gh/joaosoft/go-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/go-manager)

A package that allows you to have all your processes and data organized and with control.
After a read of the project https://gitlab.com/mandalore/go-app extracted some concepts. 

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Processes
* Configurations (with reload and write options)
* NSQ Consumers
* NSQ Producers
* Database Connections
* Web Servers
* Gateways
* Redis Connections
* Work Queues (with FIFO and LIFO modes)

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/go-manager/service
```

## Usage 
This examples are available in the project at [go-manager/bin/launcher/main.go](https://github.com/joaosoft/go-manager/tree/master/bin/launcher/main.go)

```go
//
// manager
manager := gomanager.NewManager()

// ADD ALL THE PROCESSES YOU WANT...
// ...

// DONT FORGET TO START YOUR MANAGER!
manager.Start()

```

>### Processes
```go
// --------- dummy process ---------
func dummy_process() error {
	log.Info("hello, i'm exetuting the dummy process")
	return nil
}

//
// manager: processes
process := gomanager.NewSimpleProcess(dummy_process)
if err := manager.AddProcess("process_1", process); err != nil {
    log.Errorf("MAIN: error on processes %s", err)
}
```

>### Configurations
```go
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
simpleConfig, _ := gomanager.NewSimpleConfig(dir+"/bin/launcher/data/config.json", obj)
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
```

>### NSQ Consumers 
```go
// --------- dummy nsq ---------
type dummy_nsq_handler struct{}

func (dummy *dummy_nsq_handler) HandleMessage(msg *nsqlib.Message) error {
	log.Infof("executing the handle message of NSQ with [ message: %s ]", string(msg.Body))
	return nil
}

//
// manager: nsq consumer
nsqConfigConsumer := gomanager.NewNSQConfig("topic_1", "channel_1", []string{"127.0.0.1:4161"}, []string{"127.0.0.1:4150"})
nsqConsumer, _ := gomanager.NewSimpleNSQConsumer(nsqConfigConsumer, &dummy_nsq_handler{})
manager.AddProcess("nsq_consumer_1", nsqConsumer)
```

>### NSQ Producers
```go
//
// nsq producer
nsqConfigProducer := gomanager.NewNSQConfig("topic_1", "channel_1", []string{"127.0.0.1:4150"}, []string{"127.0.0.1:4161"})
nsqProducer, _ := gomanager.NewSimpleNSQProducer(nsqConfigProducer)
manager.AddNSQProducer("nsq_producer_1", nsqProducer)
nsqProducer = manager.GetNSQProducer("nsq_producer_1")
nsqProducer.Publish("topic_1", []byte("MENSAGEM ENVIADA PARA A NSQ"), 3)
```

>### Database Connections
```go
//
// manager: database

// database - postgres
postgresConfig := gomanager.NewDBConfig("postgres", "postgres://user:password@localhost:7001?sslmode=disable")
postgresConn := gomanager.NewSimpleDB(postgresConfig)
manager.AddDB("postgres", postgresConn)

// database - mysql
mysqlConfig := gomanager.NewDBConfig("mysql", "root:password@tcp(127.0.0.1:7002)/mysql")
mysqlConn := gomanager.NewSimpleDB(mysqlConfig)
manager.AddDB("mysql", mysqlConn)
```

>### Web Servers
```go
//
// manager: web

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

// web - with http
web := gomanager.NewSimpleWebHttp(":8081")
if err := manager.AddWeb("web_http", web); err != nil {
    log.Error("error adding web process to manager")
}
web = manager.GetWeb("web_http")
web.AddRoute(http.MethodGet, "/web_http", dummy_web_http_handler)

// --------- dummy web echo ---------
func dummy_web_echo_handler(ctx echo.Context) error {
	type Example struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	return ctx.JSON(http.StatusOK, Example{Id: ctx.Param("id"), Name: "joao", Age: 29})
}

// web - with echo
web = gomanager.NewSimpleWebEcho(":8082")
if err := manager.AddWeb("web_echo", web); err != nil {
    log.Error("error adding web process to manager")
}
web = manager.GetWeb("web_echo")
web.AddRoute(http.MethodGet, "/web_echo/:id", dummy_web_echo_handler)
```

>### Gateways
```go
//
// manager: gateway
headers := map[string][]string{"Content-Type": {"application/json"}}
var body io.Reader

gateway := gomanager.NewSimpleGateway()
manager.AddGateway("gateway_1", gateway)
gateway = manager.GetGateway("gateway_1")
status, bytes, err := gateway.Request(http.MethodGet, "http://127.0.0.1:8082", "/web_echo/123", headers, body)
log.Infof("status: %d, response: %s, error? %t", status, string(bytes), err != nil)
```

>### Redis Connections
```go
//
// manager: redis
redisConfig := gomanager.NewRedisConfig("127.0.0.1", 7100, 0, "")
redisConn := gomanager.NewSimpleRedis(redisConfig)
manager.AddRedis("redis", redisConn)
```

>### Work Queues
```go
func work_handler(id string, data interface{}) error {
	log.Infof("work with the id %s and data %s done!", id, data.(string))
	return nil
}

//
// manager: workqueue
workqueueConfig := gomanager.NewWorkQueueConfig("queue_001", 1, 2, time.Second*2, gomanager.FIFO)
workqueue := gomanager.NewSimpleWorkQueue(workqueueConfig, work_handler)
manager.AddWorkQueue("queue_001", workqueue)
workqueue = manager.GetWorkQueue("queue_001")
for i := 1; i <= 1000; i++ {
    go workqueue.AddWork(fmt.Sprintf("PROCESS: %d", i), fmt.Sprintf("THIS IS MY MESSAGE %d", i))
}
if err := workqueue.Start(); err != nil {
    log.Errorf("MAIN: error on workqueue %s", err)
}
```

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
