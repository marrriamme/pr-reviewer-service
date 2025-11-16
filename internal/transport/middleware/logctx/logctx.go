package logctx

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey string

const loggerKey contextKey = "logger"

func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value(loggerKey).(*logrus.Entry)
	if !ok {
		return logrus.NewEntry(logrus.StandardLogger())
	}
	return logger
}
