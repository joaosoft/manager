env:
	docker-compose up -d manager.postgres manager.mysql manager.nsqd manager.nsqadmin manager.redis manager.rabbitmq manager.dynamodb manager.elasticsearch
run:
	go run ./bin/launcher/main.go

start:
	docker-compose build
	docker-compose up api

stop:
	docker-compose stop

build:
	go build ./...

fmt:
	go fmt ./...

vet:
	go vet ./*

gometalinter:
	gometalinter ./*