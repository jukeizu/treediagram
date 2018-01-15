package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func MakeTreediagramRequestEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, requestType interface{}) (interface{}, error) {
		request := requestType.(TreediagramRequest)

		response, err := service.Request(request)

		return response, err
	}
}
