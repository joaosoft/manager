env:
	docker-compose up -d database
	docker-compose up -d nsqd
	docker-compose up -d nsqadmin
	docker-compose up -d redis
	docker-compose up -d elasticsearch
	docker-compose up -d cassandra-1
	docker-compose up -d cassandra-2

start:
	docker-compose build
	docker build -f ${DOCKER_BASE_DOCKERFILE} -t ${DOCKER_REGISTRY}/${PROJECT_NAME}:base .
	docker-compose start
	docker-compose stop

stop:
	docker-compose stop

build:
	go build .