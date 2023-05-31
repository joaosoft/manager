run:
	go run ./examples/main.go

build:
	mkdir -p bin/
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o bin/dependency main/main.go

fmt:
	go fmt ./...

vet:
	go vet ./*

gometalinter:
	gometalinter ./*

dependency-get:
	./bin/dependency get

dependency-update:
	./bin/dependency update

dependency-reset:
	./bin/dependency reset