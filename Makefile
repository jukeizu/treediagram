VERSION=$(shell git describe --tags)
BUILD=GOARCH=amd64 go build
PROTOFILES=$(wildcard api/*/*.proto)
PBFILES=$(patsubst %.proto,%.pb.go, $(PROTOFILES))

.PHONY: all deps test proto build clean $(PROTOFILES)

all: deps test build 
deps:
	go get -t -v ./...

test:
	go vet ./...
	go test -v -race ./...

build:
	$(BUILD) -o bin/treediagram-$(VERSION) ./cmd/...

build-linux:
	CGO_ENABLED=0 GOOS=linux $(BUILD) -a -installsuffix cgo -o bin/treediagram ./cmd/...

proto: $(PBFILES)

%.pb.go: %.proto
	cd $(dir $<) && protoc $(notdir $<) --go_out=plugins=grpc:.

clean:
	@find bin -type f ! -name '*.toml' -delete -print
