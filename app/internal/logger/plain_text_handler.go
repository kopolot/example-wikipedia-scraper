package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

type PlainTextHandler struct {
	w io.Writer
	l slog.Level
}

func NewPlainTextHandler(w io.Writer, level slog.Level) *PlainTextHandler {
	return &PlainTextHandler{w: w, l: level}
}

func (h *PlainTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.l
}

func (h *PlainTextHandler) Handle(ctx context.Context, r slog.Record) error {
	ts := r.Time.Format(time.RFC3339)
	level := r.Level.String()
	msg := r.Message
	fmt.Fprintf(h.w, "[LOG] [%s] %s: %s", ts, level, msg)
	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(h.w, " %s=%v", a.Key, a.Value)
		return true
	})
	fmt.Fprintln(h.w)
	return nil
}

func (h *PlainTextHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *PlainTextHandler) WithGroup(_ string) slog.Handler      { return h }
