package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

func newLogger() *log.Logger {
	logger := log.New(os.Stderr)
	logger.SetReportTimestamp(true)
	logger.SetReportCaller(true)
	return logger
}

var Logger = newLogger()
