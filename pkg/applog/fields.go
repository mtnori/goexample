package applog

import (
	"context"
	"encoding/json"
	"log/slog"
	"maps"
)

type Fields map[string]any

func (f Fields) Merge(other Fields) Fields {
	// 使いやすいように nil レシーバのハンドリングもしておく
	if f == nil {
		return other
	}
	// 新しいものを作ってマージ
	clone := maps.Clone(f)
	for key, value := range other {
		clone[key] = value
	}
	return clone
}

func (f Fields) String() string {
	// 文字列への変換ロジックを実装
	j, err := json.Marshal(f)
	if err != nil {
		return err.Error()
	}
	return string(j)
}

type contextKey struct{}

func WithFields(ctx context.Context, fields Fields) context.Context {
	// nil レシーバのハンドリングをしているのでいきなり Merge を呼び出して OK
	return context.WithValue(ctx, contextKey{}, contextualFields(ctx).Merge(fields))
}

// Fields を取り出すのはこのパッケージだけの責務なので非公開関数で問題なし
func contextualFields(ctx context.Context) Fields {
	// コンテキストに設定されていないときは nil を返す
	f, _ := ctx.Value(contextKey{}).(Fields)
	return f
}

type Handler struct {
	parent slog.Handler
}

func WithCustomHandler(parent slog.Handler) *Handler {
	return &Handler{
		parent: parent,
	}
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	record.Add(slog.Any("fields", contextualFields(ctx)))
	return h.parent.Handle(ctx, record)
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.parent.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h.parent.WithAttrs(attrs)}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h.parent.WithGroup(name)}
}

var _ slog.Handler = (*Handler)(nil)
