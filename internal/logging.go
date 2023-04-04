package internal

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogging() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	writer := newConsoleWriter()
	log.Logger = log.Output(writer).With().Timestamp().Logger()
}

func newConsoleWriter() io.Writer {
	if os.Getenv("LAMBDA_TASK_ROOT") != "" {
		return os.Stdout
	}
	return zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: "2006-01-02T15:04:05.999Z07:00",
	}
}
