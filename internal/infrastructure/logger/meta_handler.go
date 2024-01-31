package logger

import (
	"context"

	"golang.org/x/exp/slog"
)

type MetaHandler struct {
	handlers []slog.Handler
}

func NewMetaHandler() *MetaHandler {
	return &MetaHandler{
		handlers: []slog.Handler{},
	}
}

// Add a new handler to the collection
func (h *MetaHandler) Add(handler slog.Handler) {
	h.handlers = append(h.handlers, handler)
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *MetaHandler) Enabled(ctx context.Context, level slog.Level) bool {
	var ret = true
	for i, _ := range h.handlers {
		if ret2 := h.handlers[i].Enabled(ctx, level); !ret2 {
			ret = ret2
		}
	}
	return ret
}

// WithAttrs returns a new TextHandler whose attributes consists
// of h's attributes followed by attrs.
func (h *MetaHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var newMeta = &MetaHandler{handlers: []slog.Handler{}}
	for i := range h.handlers {
		h.handlers[i] = h.handlers[i].WithAttrs(attrs)
	}
	return newMeta
}

// WithGroup add a new group to the structured log.
func (h *MetaHandler) WithGroup(name string) slog.Handler {
	var newMeta = &MetaHandler{handlers: []slog.Handler{}}
	for i := range h.handlers {
		h.handlers[i] = h.handlers[i].WithGroup(name)
	}
	return newMeta
}

// Handle formats its argument Record as a single line of space-separated
// key=value items.
func (h *MetaHandler) Handle(ctx context.Context, r slog.Record) error {
	var err error
	for i := range h.handlers {
		if err2 := h.handlers[i].Handle(ctx, r); err2 != nil {
			err = err2
		}
	}
	return err
}
