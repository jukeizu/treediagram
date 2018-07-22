package registration

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type Command struct {
	Id             bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Server         string        `json:"server"`
	Name           string        `json:"name"`
	Regex          string        `json:"regex"`
	RequireMention bool          `json:"requireMention"`
	Endpoint       string        `json:"endpoint"`
	Help           string        `json:"help"`
	Enabled        bool          `json:"enabled"`
}

type CommandQuery struct {
	Server   string `json:"server"`
	LastId   string `json:"lastId"`
	PageSize int    `json:"pageSize"`
}

type CommandQueryResult struct {
	Commands []Command `json:"commands"`
	HasMore  bool      `json:"hasMore"`
}

type AddRequest struct {
	Command Command `json:"command"`
}

type DisableRequest struct {
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

	addMessageHandler := httpTransport.NewServer(
		MakeAddRequestEndpoint(s),
		DecodeAddRequest,
		EncodeResponse,
		opts...,
	)

	disableMessageHandler := httpTransport.NewServer(
		MakeDisableRequestEndpoint(s),
		DecodeDisableRequest,
		EncodeResponse,
		opts...,
	)

	queryMessageHandler := httpTransport.NewServer(
		MakeQueryRequestEndpoint(s),
		DecodeQueryRequest,
		EncodeResponse,
		opts...,
	)
	router := mux.NewRouter()

	router.Handle("/add", addMessageHandler).Methods("POST")
	router.Handle("/disable", disableMessageHandler).Methods("POST")
	router.Handle("/query", queryMessageHandler).Methods("POST")

	return router
}
