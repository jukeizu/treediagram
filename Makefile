BUILD=GOARCH=amd64 go build -v

all: deps test build

deps:
	go get -t -v ./...

test:
	go vet ./...
	go test -v -race ./...

build:
	for CMD in `ls cmd/services`; do $(BUILD) -o bin/$$CMD-service-$$TRAVIS_TAG ./cmd/services/$$CMD; done
