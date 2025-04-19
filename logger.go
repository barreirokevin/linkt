package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
)

// ANSI color codes
const (
	Reset      = "\033[0m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Purple     = "\033[35m"
	Cyan       = "\033[36m"
	Orange     = "\033[38;5;215m"
	Faint      = "\u001b[2m"
	ResetFaint = "\u001b[22m"
)

func NewLogger(debugFlag bool) *slog.Logger {
	logLevel := &slog.LevelVar{}
	opts := PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: logLevel,
		},
	}
	if debugFlag { // set log level to debug
		logLevel.Set(slog.LevelDebug)
	} else { // otherwise, default log level is error
		logLevel.Set(slog.LevelError)
	}
	handler := NewPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)
	return logger
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewTextHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
	return h
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	var level string
	switch r.Level {
	case slog.LevelDebug:
		level = fmt.Sprintf("%s[%s]%s", Purple, r.Level.String(), Reset)
	case slog.LevelInfo:
		level = fmt.Sprintf("%s[%s]%s", Blue, r.Level.String(), Reset)
	case slog.LevelWarn:
		level = fmt.Sprintf("%s[%s]%s", Yellow, r.Level.String(), Reset)
	case slog.LevelError:
		level = fmt.Sprintf("%s[%s]%s", Red, r.Level.String(), Reset)
	}
	var attrs string
	r.Attrs(func(a slog.Attr) bool {
		attrs += fmt.Sprintf("%s%s=%+v%s", Faint, a.Key, a.Value.Any(), ResetFaint)
		return true
	})
	h.l.Println(level, r.Message, attrs)
	return nil
}
