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

func sitemap(root *url.URL, logger *slog.Logger) {
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
	// print the sitemap
	printPreorderIndent(sitemap, sitemap.Root(), -1)
}

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
