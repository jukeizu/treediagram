package publisher

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jukeizu/treediagram/services/publisher/storage"
)

type SendMessageRequest struct {
	storage.Message
	storage.Request
}

type Response struct {
	Id string `json:id`
}

func DecodeMessageRequest(_ context.Context, r *http.Request) (outputRequest interface{}, err error) {
	request := SendMessageRequest{}

	err = json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		return nil, err
	}

	return request, nil
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
		MakeSendMessageEndpoint(s),
		DecodeMessageRequest,
		EncodeResponse,
		opts...,
	)

	router := mux.NewRouter()

	router.Handle("/message", sendMessageHandler).Methods("POST")

	return router
}
