package main

import (
	"fmt"
	"net/url"
	"time"
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
			fmt.Printf("\n%s[SUCCESS]%s sitemap was created!\n", Green, Reset)
			return
		default:
			dots := []string{".  ", ".. ", "...", " ..", "  .", "   "}
			for _, s := range dots {
				fmt.Printf("\r%s[PENDING]%s collecting links%s%s%s", Orange, Reset, Orange, s, Reset)
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
