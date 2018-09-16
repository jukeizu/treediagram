FROM golang:1.11 as build
WORKDIR /go/src/github.com/jukeizu/treediagram
COPY Makefile go.mod go.sum ./
RUN make deps
ADD . .
RUN make build-linux
RUN echo "treediagram:x:100:101:/" > passwd

FROM scratch
COPY --from=build /go/src/github.com/jukeizu/treediagram/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build --chown=100:101 /go/src/github.com/jukeizu/treediagram/bin/treediagram .
USER treediagram
ENTRYPOINT ["./treediagram"]
