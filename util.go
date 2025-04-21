package main

import (
	"log/slog"
	"net/url"
	"os"
)

// Returns a URL if value is a valid URL, otherwise the program is stopped.
func isValidURL(value string, logger *slog.Logger) *url.URL {
	link, err := url.Parse(value)
	if err != nil || link.Scheme == "" || link.Host == "" { // verify url
		logger.Error("missing or invalid URL", "url", value, "error", err)
		os.Exit(0)
	}
	return link
}
