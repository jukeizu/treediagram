package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jukeizu/treediagram-handler"
	"net/http"
)

type TreediagramRequest struct {
	handler.Request
}

type TreediagramResponse struct {
	Id string `json:id`
}

func DecodeTreediagramRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	treediagramRequest := TreediagramRequest{}

	err = json.NewDecoder(r.Body).Decode(&treediagramRequest)

	if err != nil {
		return nil, err
	}

	return treediagramRequest, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func MakeHandler(s Service, logger log.Logger) http.Handler {
	opts := []httpTransport.ServerOption{
		httpTransport.ServerErrorLogger(logger),
		httpTransport.ServerErrorEncoder(EncodeError),
	}

	sendMessageHandler := httpTransport.NewServer(
		MakeTreediagramRequestEndpoint(s),
		DecodeTreediagramRequest,
		EncodeResponse,
		opts...,
	)
	router := mux.NewRouter()

	router.Handle("/treediagram", sendMessageHandler).Methods("POST")

	return router
}
