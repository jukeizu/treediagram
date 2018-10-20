package registry

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	pb "github.com/jukeizu/treediagram/api/registration"
	"github.com/shawntoffel/services-core/transport"
)

type httpBinding struct {
	logger  log.Logger
	Service pb.RegistrationServer
}

func NewHttpBinding(logger log.Logger, s pb.RegistrationServer) httpBinding {
	return httpBinding{logger, s}
}

func (h httpBinding) MakeHandler() http.Handler {
	addMessageHandler := transport.NewDefaultServer(
		h.logger,
		h.addCommandEndpoint,
		DecodeAddRequest,
	)

	disableMessageHandler := transport.NewDefaultServer(
		h.logger,
		h.disableCommandEndpoint,
		DecodeDisableRequest,
	)

	queryMessageHandler := transport.NewDefaultServer(
		h.logger,
		h.queryCommandEndpoint,
		DecodeQueryRequest,
	)

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/registration/").Subrouter()

	subrouter.Handle("/add", addMessageHandler).Methods("POST")
	subrouter.Handle("/disable", disableMessageHandler).Methods("POST")
	subrouter.Handle("/query", queryMessageHandler).Methods("POST")

	return subrouter
}

func (h httpBinding) addCommandEndpoint(ctx context.Context, r interface{}) (interface{}, error) {
	request := r.(pb.AddCommandRequest)

	return h.Service.AddCommand(ctx, &request)
}

func (h httpBinding) disableCommandEndpoint(ctx context.Context, r interface{}) (interface{}, error) {
	request := r.(pb.DisableCommandRequest)

	return h.Service.DisableCommand(ctx, &request)
}

func (h httpBinding) queryCommandEndpoint(ctx context.Context, r interface{}) (interface{}, error) {
	request := r.(pb.QueryCommandsRequest)

	return h.Service.QueryCommands(ctx, &request)
}

func DecodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := pb.AddCommandRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	return request, err
}

func DecodeDisableRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := pb.DisableCommandRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	return request, err
}

func DecodeQueryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := pb.QueryCommandsRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	return request, err
}
