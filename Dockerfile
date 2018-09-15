FROM golang:1.11 as build
WORKDIR /go/src/github.com/jukeizu/treediagram
COPY Makefile go.mod go.sum ./
RUN make deps
ADD . /go/src/github.com/jukeizu/treediagram
RUN make build-linux

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
RUN addgroup -S treediagram && adduser -S -G treediagram treediagram
USER treediagram
COPY --from=build /go/src/github.com/jukeizu/treediagram/bin/treediagram .
ENTRYPOINT ["./treediagram"]
