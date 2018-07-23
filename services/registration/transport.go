package registration

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/shawntoffel/services-core/transport"
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

	router.Handle("/add", addMessageHandler).Methods("POST")
	router.Handle("/disable", disableMessageHandler).Methods("POST")
	router.Handle("/query", queryMessageHandler).Methods("POST")

	return router
}
