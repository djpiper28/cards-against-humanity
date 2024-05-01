package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

var Logger = log.New(os.Stderr)
