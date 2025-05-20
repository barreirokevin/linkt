package main

import (
	"net/http"
	"net/url"
)

const ( // Type of page
	Unknown = iota - 1
	Internal
	External
)

// Maintains information about a page.
type Page struct {
	// HTTP request to GET this page
	request *http.Request
	// set of links on this Page
	links Set[string, int]
	// The kind of a apge, i.e. whther it is an internal page, and external
	// page, or unknown. Internal is equivlant to integer 0, External is
	// equivalent to integer 1, and Unknown is equivalent to integer -1.
	kind        int
	response    *http.Response
	requestTime string
	parentURL   string
}

// Returns a new page.
func NewPage(link *url.URL) *Page {
	return &Page{
		request: &http.Request{
			Method: http.MethodGet,
			URL:    link,
		},
		links: Set[string, int]{},
		kind:  Unknown,
	}
}
