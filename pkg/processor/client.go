package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jukeizu/contract"
)

type Client struct{}

func (c Client) Do(request interface{}, endpoint string) (string, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("could not marshal request to json: %s", err.Error())
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error sending request: %s", err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", buildError(resp)
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func buildError(resp *http.Response) error {
	if resp.StatusCode != http.StatusInternalServerError {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	errorMessage := contract.Error{}
	err := json.NewDecoder(resp.Body).Decode(&errorMessage)
	if err != nil {
		return fmt.Errorf("could not decode error: %s", err.Error())
	}

	return fmt.Errorf(errorMessage.Message)
}
