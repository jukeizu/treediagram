package processor

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Client struct{}

func (c Client) Do(request interface{}, endpoint string) (string, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return "", errors.New("could not marshal request to json: " + err.Error())
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", errors.New("error sending request: " + err.Error())
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
