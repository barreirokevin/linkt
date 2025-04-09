package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Contains all visited links
var visited Set = Set{}

// TODO: feature: test for HTTP 20x from link to ensure it's not a broken link
func main() {
	sitemap := NewTree[Page]() // contains the sitemap we are building
	page := Page{
		Request: &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Scheme: "https",
				Host:   "me-251369428744.us-central1.run.app",
				Path:   "/",
			},
		},
		Links:  Set{},
		Status: Unknown,
	}

	// add root to sitemap and Set of visited links
	_, err := sitemap.AddRoot(page)
	if err != nil {
		slog.Error(
			"error adding root page to the sitemap",
			"page", sitemap.Root().GetElement().Request.URL.String(),
			"error", err,
		)
		os.Exit(-1)
	}
	visited[page.Request.URL.String()] = Internal

	// construct http client
	client := &http.Client{}

	// call spider to start crawling from the root
	spider(client, sitemap, sitemap.Root())

	fmt.Println("\nfinished crawling root\n")
	printPreorderIndent(sitemap, sitemap.Root(), 0)
}

func spider(client *http.Client, sitemap *Tree[Page], page *Node[Page]) {
	// get page
	time.Sleep(1 * time.Second) // INFO: temporary
	resp, err := client.Do(page.GetElement().Request)
	if err != nil {
		slog.Error(
			"error getting the page",
			"page", page.GetElement().Request.URL.String(),
			"error", err,
		)
		os.Exit(0)
	}
	slog.Info(
		"got the page",
		"page", page.GetElement().Request.URL.String(),
		"status", resp.Status,
	)

	// parse page to get tree
	doc, err := html.Parse(resp.Body)
	if err != nil {
		slog.Error(
			"error parsing a page",
			"page", page.GetElement().Request.URL.String(),
			"error", err,
		)
		os.Exit(-1)
	}

	// define crawl func
	var link string
	var crawl func(*html.Node)
	crawl = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" { // node is an anchor tag
			for _, a := range n.Attr { // iterate anchor tag attributes
				if a.Key == "href" { // attribute is an href
					if strings.HasPrefix(a.Val, "/") { // href is an internal link
						link = strings.TrimSuffix(strings.TrimSpace(a.Val), "/")
						if !wasVisited(link) { // the link was not visited yet
							visited[link] = Internal
							page.GetElement().Links[link] = Internal // add internal link to Set of links
						}
					} else if !strings.HasPrefix(a.Val, "#") { // href is an external link
						link = strings.TrimSuffix(strings.TrimSpace(a.Val), "/")
						if !wasVisited(link) { // the link was not visited yet
							visited[link] = External
							page.GetElement().Links[link] = External
						}
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawl(c)
		}
	}

	// crawl the page
	crawl(doc)

	// populate the tree with Set of internal and external links
	var e Page
	for p, t := range page.GetElement().Links {
		if t == Internal { // link is internal
			e = Page{
				Request: &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme: "https",
						Host:   "me-251369428744.us-central1.run.app",
						Path:   p,
					},
				},
				Links:  Set{},
				Status: Internal,
			}
			sitemap.AddChild(page, e)
		} else { // link is external
			u, err := url.Parse(p)
			if err != nil {
				slog.Error(
					"error parsing a page URL",
					"page", p,
					"error", err,
				)
			}
			e = Page{
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    u,
				},
				Links:  Set{},
				Status: External,
			}
			sitemap.AddChild(page, e)
		}
	}

	// call the spider on the child if it has an internal link
	for _, child := range page.Children() {
		if child.GetElement().Status == Internal {
			spider(client, sitemap, child)
		}
	}
}

func wasVisited(link string) bool {
	for k, _ := range visited {
		if strings.ToLower(k) == strings.ToLower(link) {
			return true
		}
	}
	return false
}

// INFO: for debugging purposes only
func printPreorderIndent(t *Tree[Page], n *Node[Page], d int) {
	indent := strings.Repeat(" ", d*4)
	fmt.Printf("%s%+v\n", indent, n.GetElement().Request.URL.String())
	for _, c := range n.Children() {
		printPreorderIndent(t, c, d+1)
	}
}
