package api

import (
	"github.com/jukeizu/client-base"
)

type Client interface {
	Treediagram() Treediagram
}

type client struct {
	ClientConfig treediagramclient.ClientConfig
}

func NewClient(config treediagramclient.ClientConfig) Client {
	return &client{config}
}
