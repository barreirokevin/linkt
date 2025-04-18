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

// The keys of this map represent a Set, i.e. no duplicate values.
type Set map[string]int

// Maintains information about a page.
type Page struct {
	// HTTP request to GET this page
	request *http.Request
	// set of links on this Page
	links Set
	// The type of a apge, i.e. whther it is an internal page, and external
	// page, or unknown. Internal is equivlant to integer 0, External is
	// equivalent to integer 1, and Unknown is equivalent to integer -1.
	t int
}

// Returns a new page.
func NewPage(link *url.URL) *Page {
	return &Page{
		request: &http.Request{
			Method: http.MethodGet,
			URL:    link,
		},
		links: Set{},
		t:     Unknown,
	}
}

// Sets the type of page to unknown, internal or external, i.e. -1, 0,
// and 1, respectively,
func (p *Page) SetType(t int) {
	p.t = t
}

// Returns the type of page, i.e. -1, 0, or 1. Each value represents
// unknown, internal, or external, respectively,
func (p Page) Type() int {
	return p.t
}

// Returns a Set of links in this page.
func (p Page) Links() Set {
	return p.links
}

// Returns the HTTP request to GET this page.
func (p Page) Request() *http.Request {
	return p.request
}

// Returns a string representing the URL of this Page.
func (p Page) URL() string {
	return p.request.URL.String()
}
