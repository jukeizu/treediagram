package api

import restclient "github.com/shawntoffel/go-restclient"

type TreediagramRequest struct {
	Source        string `json:"source"`
	CorrelationId string `json:"correlationId"`
	Bot           User   `json:"bot"`
	Author        User   `json:"author"`
	ChannelId     string `json:"channelId"`
	ServerId      string `json:"serverId"`
	Mentions      Users  `json:"mentions"`
	Content       string `json:"content"`
}

type Users []User

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TreediagramResponse struct {
	Id string `json:"id"`
}

type treediagram struct {
	client *client
}

type Treediagram interface {
	Request(TreediagramRequest) (TreediagramResponse, error)
}

func (c *client) Treediagram() Treediagram {
	return &treediagram{c}
}

func (t *treediagram) Request(request TreediagramRequest) (TreediagramResponse, error) {
	response := TreediagramResponse{}

	err := restclient.Post(t.client.ClientConfig.BaseUrl+"/treediagram", request, &response)

	return response, err
}
