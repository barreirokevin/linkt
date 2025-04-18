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

type Sitemap struct {
	Tree[Page]
	logger *slog.Logger
}

// Returns an empty sitemap.
func NewSitemap(logger *slog.Logger) *Sitemap {
	return &Sitemap{logger: logger}
}

// Builds a sitemap for the URL root. The logger outputs logs to the logs directory in the
// current working directory.
func (s *Sitemap) Build(root *url.URL, done chan bool) {
	visited := &Set{} // contains all visited links
	page := *NewPage(root)
	// add root to sitemap and Set of visited links
	_, err := s.AddRoot(page)
	if err != nil {
		s.logger.Error(
			"error adding root page to the sitemap",
			"page", page.URL(),
			"error", err,
		)
		os.Exit(-1)
	}
	// construct http client
	client := &http.Client{}
	// call spider to start crawling from the root
	spider(client, s, s.Root(), visited, s.logger)
	// send signal to sitemap animation that sitemap is done
	done <- true
	// print the sitemap to stdout
	fmt.Printf("%s\n", s.String())
}

// Returns the tree as a string that displays the hiearachy.
func (s *Sitemap) String() string {
	str := ""

	var preorderIndent func(s *Sitemap, n *Node[Page], d int)
	preorderIndent = func(s *Sitemap, n *Node[Page], d int) {
		if reflect.DeepEqual(s.Root(), n) {
			// the current node is the root
			str += fmt.Sprintf("%s\n", n.GetElement().URL())

		} else if len(n.Children()) == 0 &&
			reflect.DeepEqual(n.GetParent().Children()[len(n.GetParent().Children())-1], n) {
			// the current node is the last child
			indent := strings.Repeat(" ", d*4)
			if d > 0 && d%2 == 0 {
				str += fmt.Sprintf("│%s └─── %+v\n", indent, n.GetElement().URL())
			} else if d > 0 && d%2 != 0 {
				str += fmt.Sprintf("│%s└─── %+v\n", indent, n.GetElement().URL())
			} else {
				str += fmt.Sprintf("%s└─── %+v\n", indent, n.GetElement().URL())
			}

		} else {
			// the current node is not the last child
			indent := strings.Repeat(" ", d*4)
			if d > 0 && d%2 == 0 {
				str += fmt.Sprintf("│%s ├─── %+v\n", indent, n.GetElement().URL())
			} else if d > 0 && d%2 != 0 {
				str += fmt.Sprintf("│%s├─── %+v\n", indent, n.GetElement().URL())
			} else {
				str += fmt.Sprintf("%s├─── %+v\n", indent, n.GetElement().URL())
			}
		}

		// recursively call preorderIndent for each child
		for _, c := range n.Children() {
			preorderIndent(s, c, d+1)
		}
	}

	// begin preorderIndent with the tree's root
	preorderIndent(s, s.Root(), -1)
	return str
}
