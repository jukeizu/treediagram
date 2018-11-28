package intent

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	pb "github.com/jukeizu/treediagram/api/protobuf-spec/intent"
	"github.com/shawntoffel/services-core/transport"
)

type httpBinding struct {
	logger  log.Logger
	Service pb.IntentRegistryServer
}

func NewHttpBinding(logger log.Logger, s pb.IntentRegistryServer) httpBinding {
	return httpBinding{logger, s}
}

func (h httpBinding) MakeHandler() http.Handler {
	addMessageHandler := transport.NewDefaultServer(
		h.logger,
		h.addIntentEndpoint,
		DecodeAddRequest,
	)

	disableMessageHandler := transport.NewDefaultServer(
		h.logger,
		h.disableIntentEndpoint,
		DecodeDisableRequest,
	)

	queryMessageHandler := transport.NewDefaultServer(
		h.logger,
		h.queryIntentEndpoint,
		DecodeQueryRequest,
	)

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/intent/").Subrouter()

	subrouter.Handle("/add", addMessageHandler).Methods("POST")
	subrouter.Handle("/disable", disableMessageHandler).Methods("POST")
	subrouter.Handle("/query", queryMessageHandler).Methods("POST")

	return subrouter
}

func (h httpBinding) addIntentEndpoint(ctx context.Context, r interface{}) (interface{}, error) {
	request := r.(pb.AddIntentRequest)

	return h.Service.AddIntent(ctx, &request)
}

func (h httpBinding) disableIntentEndpoint(ctx context.Context, r interface{}) (interface{}, error) {
	request := r.(pb.DisableIntentRequest)

	return h.Service.DisableIntent(ctx, &request)
}

func (h httpBinding) queryIntentEndpoint(ctx context.Context, r interface{}) (interface{}, error) {
	request := r.(pb.QueryIntentsRequest)

	return h.Service.QueryIntents(ctx, &request)
}

func DecodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := pb.AddIntentRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	return request, err
}

func DecodeDisableRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := pb.DisableIntentRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	return request, err
}

func DecodeQueryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := pb.QueryIntentsRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	return request, err
}
