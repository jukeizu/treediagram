FROM golang:1.11 as build
ADD . /go/src/github.com/jukeizu/treediagram
WORKDIR /go/src/github.com/jukeizu/treediagram
RUN make deps build-linux

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
RUN addgroup -S treediagram && adduser -S -G treediagram treediagram
USER treediagram
COPY --from=build /go/src/github.com/jukeizu/treediagram/bin/treediagram .
ENTRYPOINT ["./treediagram"]
