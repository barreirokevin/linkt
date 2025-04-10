package main

import (
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"golang.org/x/net/html"
)

func spider(client *http.Client, sitemap *Tree[Page], page *Node[Page], visited *Set, logger *slog.Logger) {
	// get page
	resp, err := client.Do(page.GetElement().Request)
	if err != nil {
		logger.Error(
			"error getting the page",
			"page", page.GetElement().Request.URL.String(),
			"error", err,
		)
		os.Exit(0)
	}
	logger.Info(
		"got the page",
		"page", page.GetElement().Request.URL.String(),
		"status", resp.Status,
	)

	// add root page as internal link to Set of links
	if reflect.DeepEqual(sitemap.Root(), page) {
		(*visited)[page.GetElement().Request.URL.String()] = Internal
		page.GetElement().Links[page.GetElement().Request.URL.String()] = Internal
	}

	// parse page to get tree
	doc, err := html.Parse(resp.Body)
	if err != nil {
		logger.Error(
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
						if !contains(link, visited) { // the link was not visited yet
							(*visited)[link] = Internal
							page.GetElement().Links[link] = Internal // add internal link to Set of links
						}
					} else if !strings.HasPrefix(a.Val, "#") { // href is an external link
						link = strings.TrimSuffix(strings.TrimSpace(a.Val), "/")
						if !contains(link, visited) { // the link was not visited yet
							(*visited)[link] = External
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
						Scheme: sitemap.Root().GetElement().Request.URL.Scheme,
						Host:   sitemap.Root().GetElement().Request.URL.Host,
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
				logger.Error(
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
			// TODO: call spider in its own go routine so that each spider crawls concurrently
			spider(client, sitemap, child, visited, logger)
		}
	}
}
