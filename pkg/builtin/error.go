package builtin

import (
	"fmt"

	"github.com/jukeizu/contract"
)

type ParseError struct {
	Message string
}

func NewParseError(format string, a ...interface{}) ParseError {
	return ParseError{Message: fmt.Sprintf(format, a...)}
}

func (e ParseError) Error() string {
	return e.Message
}

func FormatErrorResponse(err error) (*contract.Response, error) {
	if err == nil {
		return nil, nil
	}

	switch err.(type) {
	case ParseError:
		return contract.StringResponse(err.Error()), nil
	}

	return nil, err
}
