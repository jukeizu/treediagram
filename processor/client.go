package processor

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
)

type Client struct{}

func (c Client) Do(request interface{}, endpoint string) (*processingpb.Response, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, errors.New("could not marshal request to json: " + err.Error())
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New("error sending request: " + err.Error())
	}

	defer resp.Body.Close()

	reply := &processingpb.Response{}

	err = c.decodeJSON(resp.Body, reply)
	if err != nil {
		return nil, errors.New("failed to decode JSON response body: " + err.Error())
	}

	return reply, nil
}

func (c Client) decodeJSON(body io.Reader, into interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.New("error reading response body: " + err.Error())
	}

	if b == nil || len(b) < 1 {
		return errors.New("body did not contain any bytes")
	}

	return json.Unmarshal(b, &into)
}
