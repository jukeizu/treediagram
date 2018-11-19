TAG=$(shell git describe --tags)
VERSION=$(TAG:v%=%)
REPO=jukeizu/treediagram
GO=GO111MODULE=on go
BUILD=GOARCH=amd64 $(GO) build -ldflags="-s -w -X main.Version=$(VERSION)" 
PROTOFILES=$(wildcard api/protobuf-spec/*/*.proto)
PBFILES=$(patsubst %.proto,%.pb.go, $(PROTOFILES))

.PHONY: all deps test proto build clean $(PROTOFILES)

all: deps test build 
deps:
	$(GO) get -d -v ./...

test:
	$(GO) vet ./...
	$(GO) test -v -race ./...

build:
	$(BUILD) -o bin/treediagram-$(VERSION) ./cmd/...

build-linux:
	CGO_ENABLED=0 GOOS=linux $(BUILD) -a -installsuffix cgo -o bin/treediagram ./cmd/...

docker-build:
	docker build -t $(REPO):$(VERSION) .

docker-deploy:
	docker tag $(REPO):$(VERSION) $(REPO):latest
	docker push $(REPO):$(VERSION)
	docker push $(REPO):latest

proto: $(PBFILES)

%.pb.go: %.proto
	cd $(dir $<) && protoc $(notdir $<) --go_out=plugins=grpc:.

clean:
	@find bin -type f ! -name '*.toml' -delete -print
