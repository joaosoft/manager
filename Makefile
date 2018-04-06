env:
	docker-compose up -d postgres mysql nsqd nsqadmin redis
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