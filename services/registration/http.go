package registration

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/shawntoffel/services-core/transport"
)

type AddRequest struct {
	Command Command `json:"command"`
}

type DisableRequest struct {
	Id string `json:"id"`
}

type DisableResponse struct {
	Id string `json:"id"`
}

func DecodeAddRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	addRequest := AddRequest{}

	err = json.NewDecoder(r.Body).Decode(&addRequest)

	if err != nil {
		return nil, err

	}

	return addRequest, nil
}

func DecodeDisableRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	disableRequest := DisableRequest{}

	err = json.NewDecoder(r.Body).Decode(&disableRequest)

	if err != nil {
		return nil, err

	}

	return disableRequest, nil
}

func DecodeQueryRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	query := CommandQuery{}

	err = json.NewDecoder(r.Body).Decode(&query)

	if err != nil {
		return nil, err

	}

	return query, nil
}

func MakeHandler(s Service, logger log.Logger) http.Handler {
	addMessageHandler := transport.NewDefaultServer(
		logger,
		MakeAddRequestEndpoint(s),
		DecodeAddRequest,
	)

	disableMessageHandler := transport.NewDefaultServer(
		logger,
		MakeDisableRequestEndpoint(s),
		DecodeDisableRequest,
	)

	queryMessageHandler := transport.NewDefaultServer(
		logger,
		MakeQueryRequestEndpoint(s),
		DecodeQueryRequest,
	)

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/registration/").Subrouter()

	subrouter.Handle("/add", addMessageHandler).Methods("POST")
	subrouter.Handle("/disable", disableMessageHandler).Methods("POST")
	subrouter.Handle("/query", queryMessageHandler).Methods("POST")

	return subrouter
}
