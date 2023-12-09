package applog

import (
	"context"
	"log/slog"
	"os"
)

type Logger interface {
	Info(ctx context.Context, message string)
	Error(ctx context.Context, message string)
}

func NewLogger() Logger {
	return &logger{
		inner: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

var _ Logger = (*logger)(nil)

type logger struct {
	inner *slog.Logger
}

func (l *logger) Info(ctx context.Context, message string) {
	l.inner.InfoContext(ctx, message)
}

func (l *logger) Error(ctx context.Context, message string) {
	l.inner.ErrorContext(ctx, message)
}

var DefaultLogger = NewLogger()

func Info(ctx context.Context, message string) {
	DefaultLogger.Info(ctx, message)
}

func Error(ctx context.Context, message string) {
	DefaultLogger.Error(ctx, message)
}
