package log

import (
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/andrew-d/go-webapp-skeleton/conf"
)

type private struct{}

var contextKey private

func NewLogger() *logrus.Logger {
	log := logrus.New()

	if conf.C.Debug {
		log.Level = logrus.DebugLevel
	} else {
		log.Level = logrus.InfoLevel
		// TODO: print JSON in prod?
	}

	return log
}

func FromContext(c context.Context) *logrus.Logger {
	return c.Value(contextKey).(*logrus.Logger)
}

func NewContext(parent context.Context, log *logrus.Logger) context.Context {
	return context.WithValue(parent, contextKey, log)
}
