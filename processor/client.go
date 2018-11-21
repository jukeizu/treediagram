package processor

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
)

type Client struct{}

func (c Client) Do(command Command) (*processing.Response, error) {
	body, err := json.Marshal(command.Request)
	if err != nil {
		return nil, errors.New("could not marshal command request to json: " + err.Error())
	}

	resp, err := http.Post(command.Intent.Endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New("error sending command request: " + err.Error())
	}

	defer resp.Body.Close()

	reply := &processing.Response{}

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
