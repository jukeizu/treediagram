VERSION=$(shell git describe --tags)
BUILD=GOARCH=amd64 go build

.PHONY: all deps test proto build clean

all: deps test build

deps:
	go get -t -v ./...

test:
	go vet ./...
	go test -v -race ./...

proto:
	cd api/registration && protoc registration.proto --go_out=plugins=grpc:.
	cd api/receiving && protoc receiving.proto --go_out=plugins=grpc:.
	cd api/publishing && protoc publishing.proto --go_out=plugins=grpc:.
	cd api/user && protoc user.proto --go_out=plugins=grpc:.
	cd api/scheduling && protoc scheduling.proto --go_out=plugins=grpc:.

build:
	$(BUILD) -o bin/treediagram-$(VERSION) ./cmd/...

build-linux:
	CGO_ENABLED=0 GOOS=linux $(BUILD) -a -installsuffix cgo -o bin/treediagram ./cmd/...


clean:
	@find bin -type f ! -name '*.toml' -delete -print
