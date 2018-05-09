VERSION=$(shell git describe --tags)
BUILD=GOARCH=amd64 go build -v

all: deps test build

deps:
	go get -t -v ./...

test:
	go vet ./...
	go test -v -race ./...

build:
	for CMD in `ls cmd/services`; do $(BUILD) -o bin/$$CMD-service-$(VERSION) ./cmd/services/$$CMD; done
	for CMD in `ls cmd/listeners`; do $(BUILD) -o bin/$$CMD-listener-$(VERSION) ./cmd/listeners/$$CMD; done
