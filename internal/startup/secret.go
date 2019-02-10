package startup

import (
	"errors"
	"io/ioutil"
)

func ReadSecretFromFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.New("could not read secrets file: " + filename + ": " + err.Error())
	}

	return string(b), nil
}
