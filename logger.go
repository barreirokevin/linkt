package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"time"
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

// Creates a new logger. By default the log level is set to slog.LevelError. If debugFlag
// is true, then the log level is set to slog.LevelDebug.
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

// Custom handler for the logger to modify the styling of the log output.
type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

// Options to specify when using a PrettyHandler.
type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

// Returns a PrettyHandler for the logger to putput logs with custom styling.
func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewTextHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
	return h
}

// Customizes the styling of the log output.
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

// Prints text to stdout in such a manner that it gives the impression the text is moving
// while Spider is building the sitemap.
func sitemapAnimation(done chan bool) {
	for {
		select {
		case <-done:
			fmt.Printf("\n%s[SUCCESS]%s sitemap was created!\n", Green, Reset)
			return
		default:
			dots := []string{".  ", ".. ", "...", " ..", "  .", "   "}
			for _, s := range dots {
				fmt.Printf("\r%s[PENDING]%s collecting links%s%s%s\n", Orange, Reset, Orange, s, Reset)
				time.Sleep((1 * time.Second) / 4)
			}
		}
	}
}
