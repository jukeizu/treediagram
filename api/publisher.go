package api

import (
	"github.com/jukeizu/treediagram/services/publishing/storage"
	restclient "github.com/shawntoffel/go-restclient"
)

type SendMessageRequest struct {
	storage.Message
	storage.Request
}

type Response struct {
	Id string `json:"id"`
}

type publisher struct {
	client *client
}

type Publisher interface {
	SendMessage(SendMessageRequest) (Response, error)
}

func (c *client) Publisher() Publisher {
	return &publisher{c}
}

func (p *publisher) SendMessage(request SendMessageRequest) (Response, error) {
	response := Response{}

	err := restclient.Post(p.client.ClientConfig.BaseUrl+"/message", request, &response)

	return response, err
}
