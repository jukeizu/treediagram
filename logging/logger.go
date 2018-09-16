package logging

import (
	"os"

	"github.com/go-kit/kit/log"
)

func NewLogger(component string, version string) log.Logger {
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "component", component)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "version", version)

	return logger
}
