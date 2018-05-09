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

type publishing struct {
	client *client
}

type Publishing interface {
	SendMessage(SendMessageRequest) (Response, error)
}

func (c *client) Publishing() Publishing {
	return &publishing{c}
}

func (p *publishing) SendMessage(request SendMessageRequest) (Response, error) {
	response := Response{}

	err := restclient.Post(p.client.ClientConfig.BaseUrl+"/message", request, &response)

	return response, err
}
