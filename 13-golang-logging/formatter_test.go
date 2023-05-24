package golanglogging

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestFormatter(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.Info("Hello Logging")
	logger.Warn("Hellog Logging")
	logger.Error("Hello Logging")
}
