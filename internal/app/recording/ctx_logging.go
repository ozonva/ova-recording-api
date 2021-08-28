package recording

import (
	"context"
	"github.com/sirupsen/logrus"
)

type LoggerKey string
const loggerKeyName = LoggerKey("Logger")

func AddValue(ctx context.Context, key string, value interface{}) context.Context {
	log := ctx.Value(loggerKeyName)
	if log == nil {
		log = logrus.NewEntry(logrus.StandardLogger())
	}
	log = log.(*logrus.Entry).WithField(key, value)
	ctx = context.WithValue(ctx, loggerKeyName, log)
	ctx = context.WithValue(ctx, key, value)
	return ctx
}

func GetLogger(ctx context.Context) *logrus.Entry {
	log := ctx.Value(loggerKeyName)
	if log == nil {
		logrus.Errorf("Cannot get logger from context by key %s", loggerKeyName)
		log = logrus.NewEntry(logrus.StandardLogger())
	}
	return log.(*logrus.Entry)
}
