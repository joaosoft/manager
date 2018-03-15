env:
	docker-compose up -d database
	docker-compose up -d nsqd
	docker-compose up -d nsqadmin
	docker-compose up -d redis
	docker-compose up -d elasticsearch
	docker-compose up -d cassandra-1
	docker-compose up -d cassandra-2
run:
	go run ./bin/launcher/main.go

start:
	docker-compose build
	docker-compose up api

stop:
	docker-compose stop

build:
	go build .

fmt:
	go fmt ./...

vet:
	go vet ./*

gometalinter:
	gometalinter ./*