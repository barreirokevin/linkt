package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"golang.org/x/net/html"
)

// Called on a page and begins crawling therefrom to obtain all anchor tags within the domain of page.
func spider(client *http.Client, sitemap *Sitemap, node *Node[Page], visited *Set, logger *slog.Logger) {
	// get page
	currentPage := node.GetElement()
	resp, err := client.Do(currentPage.Request())
	if err != nil {
		logger.Error(
			"error getting the page",
			"page", currentPage.URL(),
			"error", err,
		)
		os.Exit(0)
	}
	logger.Info(
		"got the page",
		"page", currentPage.URL(),
		"status", resp.Status,
	)

	// add root page as internal link to Set of links
	if reflect.DeepEqual(sitemap.Root(), currentPage) {
		(*visited)[currentPage.URL()] = Internal
		currentPage.Links()[currentPage.URL()] = Internal
	}

	// parse page to get tree
	doc, err := html.Parse(resp.Body)
	if err != nil {
		logger.Error(
			"error parsing a page",
			"page", currentPage.URL(),
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
							currentPage.Links()[link] = Internal // add internal link to Set of links
						}
					} else if !strings.HasPrefix(a.Val, "#") { // href is an external link
						link = strings.TrimSuffix(strings.TrimSpace(a.Val), "/")
						if !contains(link, visited) { // the link was not visited yet
							(*visited)[link] = External
							currentPage.Links()[link] = External
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
	for p, t := range currentPage.Links() {
		if t == Internal { // link is internal
			link, err := url.Parse(fmt.Sprintf("%s%s", sitemap.Root().GetElement().URL(), p))
			if err != nil {
				logger.Error(
					"error parsing a page URL",
					"page", link,
					"error", err,
				)
			}
			page := *NewPage(link)
			page.SetType(Internal)
			sitemap.AddChild(node, page)

		} else { // link is external
			link, err := url.Parse(p)
			if err != nil {
				logger.Error(
					"error parsing a page URL",
					"page", p,
					"error", err,
				)
			}
			page := *NewPage(link)
			page.SetType(External)
			sitemap.AddChild(node, page)
		}
	}

	// call the spider on the child if it has an internal link
	for _, child := range node.Children() {
		if child.GetElement().Type() == Internal {
			spider(client, sitemap, child, visited, logger)
		}
	}
}
