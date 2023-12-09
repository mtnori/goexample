package applog

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
)

type _TraceIDContextKey struct{}

var traceIDContextKey = _TraceIDContextKey{}

func WithTraceID(ctx context.Context) context.Context {
	traceID := uuid.NewString()
	return context.WithValue(ctx, traceIDContextKey, traceID)
}

type TraceIDHandler struct {
	parent slog.Handler
}

func WithTraceIDHandler(parent slog.Handler) *TraceIDHandler {
	return &TraceIDHandler{
		parent: parent,
	}
}

func (h *TraceIDHandler) Handle(ctx context.Context, record slog.Record) error {
	record.Add(slog.String("traceID", ctx.Value(traceIDContextKey).(string)))
	return h.parent.Handle(ctx, record)
}
func (h *TraceIDHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.parent.Enabled(ctx, level)
}

func (h *TraceIDHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TraceIDHandler{h.parent.WithAttrs(attrs)}
}

func (h *TraceIDHandler) WithGroup(name string) slog.Handler {
	return &TraceIDHandler{h.parent.WithGroup(name)}
}

var _ slog.Handler = (*TraceIDHandler)(nil)
