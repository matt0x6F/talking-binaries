package logger

import (
	"io"
	"log"
)

// Logger is a logger backed by Go's standard logger
var Logger *logger

type logger struct {
	logger  *log.Logger
	verbose bool
}

func (l logger) Info(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l logger) Debug(format string, v ...interface{}) {
	if l.verbose {
		l.logger.Printf(format, v...)
	}
}

func New(out io.Writer, debug bool) *logger {
	Logger = new(logger)

	Logger.logger = log.New(out, "[main] ", log.Flags())
	Logger.verbose = debug

	return Logger
}
