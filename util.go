package main

import "net/url"

// Returns true if value is a valid URL, otherwise it returns false.
func isValidURL(value string) bool {
	url, err := url.Parse(value)
	if err != nil {
		return false
	}
	if url.Scheme == "" || url.Host == "" {
		return false
	}
	return true
}
