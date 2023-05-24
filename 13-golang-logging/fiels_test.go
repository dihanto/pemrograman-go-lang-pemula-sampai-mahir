package golanglogging

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestFiled(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithField("username", "Dihanto").Info("Hello Logging")

	logger.WithField("username", "Di").
		WithField("name", "Kurniawan").
		Info("Hello Logging")
}

func TestFileds(t *testing.T) {

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithFields(logrus.Fields{
		"username": "Di",
		"name":     "to",
	}).Info("Hello Logging")
}
