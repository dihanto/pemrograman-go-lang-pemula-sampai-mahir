package golanglogging

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestEntry(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.Info("Hello Logging")
	logger.WithField("username", "Dihanto").Info("Hello Logging")

	entry := logrus.NewEntry(logger)
	entry.WithField("username", "to")
	entry.Info("Hello Entry")
}
