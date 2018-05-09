package publishing

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func MakeSendMessageEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, requestType interface{}) (interface{}, error) {
		request := requestType.(SendMessageRequest)

		response, err := service.SendMessage(request)

		return response, err
	}
}
