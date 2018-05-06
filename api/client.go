package api

import (
	"github.com/jukeizu/client-base"
)

type Client interface {
}

type client struct {
	ClientConfig treediagramclient.ClientConfig
}

func NewClient(config treediagramclient.ClientConfig) Client {
	return &client{config}
}
