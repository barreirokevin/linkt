package main

import "net/http"

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
	Request *http.Request
	// set of links on this Page
	Links Set
	// The status of a apge, i.e. whther it is an internal page, and external
	// page, or unknown. Internal is equivlant to integer 0, External is
	// equivalent to integer 1, and Unknown is equivalent to integer -1.
	Status int
}

// Returns a string representing the URL of this Page.
func (p Page) URL() string {
	return p.Request.URL.String()
}
