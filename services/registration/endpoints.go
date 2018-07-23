package registration

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeAddRequestEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, requestType interface{}) (interface{}, error) {
		request := requestType.(AddRequest)

		response, err := service.Add(request.Command)

		return response, err
	}
}

func MakeDisableRequestEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, requestType interface{}) (interface{}, error) {
		request := requestType.(DisableRequest)

		err := service.Disable(request.Id)

		return DisableResponse{Id: request.Id}, err
	}
}

func MakeQueryRequestEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, requestType interface{}) (interface{}, error) {
		request := requestType.(CommandQuery)

		result, err := service.Query(request)

		return result, err
	}
}
