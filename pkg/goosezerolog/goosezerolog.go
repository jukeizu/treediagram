package goosezerolog

import (
	"fmt"

	"github.com/rs/zerolog"
)

type GooseLogger struct {
	Logger zerolog.Logger
}

func New(zl zerolog.Logger) GooseLogger {
	logger := zl.With().Str("migrator", "goose").Logger()

	return GooseLogger{Logger: logger}
}

func (g GooseLogger) Fatal(v ...interface{}) {
	g.Logger.Fatal().Msg(fmt.Sprint(v...))
}

func (g GooseLogger) Fatalf(format string, v ...interface{}) {
	g.Logger.Fatal().Msg(fmt.Sprintf(format, v...))
}

func (g GooseLogger) Print(v ...interface{}) {
	g.Logger.Info().Msg(fmt.Sprint(v...))
}

func (g GooseLogger) Println(v ...interface{}) {
	g.Print(v...)
}

func (g GooseLogger) Printf(format string, v ...interface{}) {
	g.Logger.Info().Msg(fmt.Sprintf(format, v...))
}
