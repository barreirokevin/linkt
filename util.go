package main

import (
	"fmt"
	"net/url"
	"time"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Orange = "\033[38;5;215m"
)

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

func sitemapAnimation(done chan bool) {
	for {
		select {
		case <-done:
			fmt.Printf("\n%s=>%s done!\n\n", Green, Reset)
			return
		default:
			dots := []string{".  ", ".. ", "...", " ..", "  .", "   "}
			for _, s := range dots {
				fmt.Printf("\r%s=>%s collecting links%s%s%s", Orange, Reset, Orange, s, Reset)
				time.Sleep((1 * time.Second) / 4)
			}
		}
	}
}

func contains(value string, set *Set) bool {
	_, found := (*set)[value]
	if found {
		return true
	}
	return false
}
