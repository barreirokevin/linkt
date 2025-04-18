package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

// Builds a sitemap for the URL root. The logger outputs logs to the logs directory in the
// current working directory.
func sitemap(root *url.URL, done chan bool, logger *slog.Logger) {
	visited := &Set{}          // Contains all visited links
	sitemap := NewTree[Page]() // contains the sitemap we are building
	page := Page{
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    root,
		},
		Links:  Set{},
		Status: Unknown,
	}
	// add root to sitemap and Set of visited links
	_, err := sitemap.AddRoot(page)
	if err != nil {
		logger.Error(
			"error adding root page to the sitemap",
			"page", sitemap.Root().GetElement().Request.URL.String(),
			"error", err,
		)
		os.Exit(-1)
	}
	// construct http client
	client := &http.Client{}
	// call spider to start crawling from the root
	spider(client, sitemap, sitemap.Root(), visited, logger)
	// send signal to dots animation that sitemap is done
	done <- true
	// print the sitemap
	// printPreorderIndent(sitemap, sitemap.Root(), -1)
}

// Prints the sitemap to stdout with a preorder traversal of tree t.
func printPreorderIndent(t *Tree[Page], n *Node[Page], d int) {
	if reflect.DeepEqual(t.Root(), n) {
		fmt.Printf("\r%s%+v%s\n", Orange, n.GetElement().Request.URL.String(), Reset)
	} else if len(n.Children()) == 0 && reflect.DeepEqual(n.GetParent().Children()[len(n.GetParent().Children())-1], n) {
		indent := strings.Repeat(" ", d*4)
		if d > 0 && d%2 == 0 {
			fmt.Printf("%s│%s └─── %s%+v\n", Orange, indent, Reset, n.GetElement().Request.URL.String())
		} else if d > 0 && d%2 != 0 {
			fmt.Printf("%s│%s└─── %s%+v\n", Orange, indent, Reset, n.GetElement().Request.URL.String())
		} else {
			fmt.Printf("%s%s└─── %s%+v\n", Orange, indent, Reset, n.GetElement().Request.URL.String())
		}
	} else {
		indent := strings.Repeat(" ", d*4)
		if d > 0 && d%2 == 0 {
			fmt.Printf("%s│%s ├─── %s%+v\n", Orange, indent, Reset, n.GetElement().Request.URL.String())
		} else if d > 0 && d%2 != 0 {
			fmt.Printf("%s│%s├─── %s%+v\n", Orange, indent, Reset, n.GetElement().Request.URL.String())
		} else {
			fmt.Printf("%s%s├─── %s%+v\n", Orange, indent, Reset, n.GetElement().Request.URL.String())
		}
	}

	for _, c := range n.Children() {
		printPreorderIndent(t, c, d+1)
	}
}
