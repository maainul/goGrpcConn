package logging

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	isTerminal = terminal.IsTerminal(int(os.Stdout.Fd()))
	isTest     = strings.HasSuffix(os.Args[0], ".test")
	flagMu     = sync.Mutex{}
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Level = logrus.InfoLevel
	if isTest {
		testing.Init()
		flagMu.Lock()
		if !flag.Parsed() {
			flag.Parse()
		}
		flagMu.Unlock()

		if !testing.Verbose() {
			logger.Level = logrus.FatalLevel
		}
	}

	logger.Out = os.Stdout
	if isTerminal {
		logger.Formatter = &logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		}
	} else {
		logger.Formatter = &FluentdFormatter{
			TimestampFormat: time.RFC3339,
		}
	}

	logger.SetReportCaller(true)
	return logger
}

func WithError(err error, logger logrus.FieldLogger) *logrus.Entry {
	return logger.WithError(err).WithField("stacktrace", fmt.Sprintf("%+v", err))
}
