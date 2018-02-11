package command

import ()

type Request struct {
	Source    string `json:"source"`
	UserId    string `json:"userId"`
	ChannelId string `json:"channelId"`
	Content   string `json:"content"`
}

type Response struct {
	Source    string `json:"source"`
	UserId    string `json:"userId"`
	ChannelId string `json:"channelId"`
	Content   string `json:"content"`
}

type Command interface {
	IsCommand(request Request) bool
	Handle(request Request) Response
}

func buildResponse(request Request) Response {
	response := Response{
		Source:    request.Source,
		UserId:    request.UserId,
		ChannelId: request.ChannelId,
	}

	return response
}
