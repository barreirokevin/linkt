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

const ( // actions a spider can take
	SITEMAP = iota
	TEST
	SCREENSHOTS
)

// A spider with capabilities such as building a sitemap, testing links,
// and taking screenshots for a site.
type Spider struct {
	client  *http.Client
	logger  *slog.Logger
	visited *Set
}

// Returns a new spider with an HTTP client.
func NewSpider(logger *slog.Logger) *Spider {
	c := &http.Client{}
	return &Spider{client: c, logger: logger, visited: &Set{}}
}

// Enables the spider to build a sitemap starting from the root URL.
func (s *Spider) DoSitemap(root *url.URL, done chan bool) {
	sitemap := NewSitemap(s.logger)
	page := *NewPage(root)
	_, err := sitemap.AddRoot(page)
	if err != nil {
		s.logger.Error(
			"error adding root page to the sitemap",
			"page", page.URL(),
			"error", err,
		)
		os.Exit(-1)
	}
	// start recursively building the sitemap
	s.build(sitemap, sitemap.Root())
	// send signal to sitemap animation that sitemap is done
	done <- true
	// print the sitemap to stdout
	fmt.Printf("%s\n", sitemap.String())
}

// Called on a page and begins crawling therefrom to obtain all anchor tags within the domain of page to build a sitemap.
func (s *Spider) build(sitemap *Sitemap, node *Node[Page]) {
	// get page
	currentPage := node.GetElement()
	resp, err := s.client.Do(currentPage.Request())
	if err != nil {
		s.logger.Error(
			"error getting the page",
			"page", currentPage.URL(),
			"error", err,
		)
		os.Exit(0)
	}
	s.logger.Info(
		"got the page",
		"page", currentPage.URL(),
		"status", resp.Status,
	)

	// add root page as internal link to Set of links
	if reflect.DeepEqual(sitemap.Root(), currentPage) {
		(*s.visited)[currentPage.URL()] = Internal
		currentPage.Links()[currentPage.URL()] = Internal
	}

	// parse page to get tree
	doc, err := html.Parse(resp.Body)
	if err != nil {
		s.logger.Error(
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
						if !contains(link, s.visited) { // the link was not visited yet
							(*s.visited)[link] = Internal
							currentPage.Links()[link] = Internal // add internal link to Set of links
						}
					} else if !strings.HasPrefix(a.Val, "#") { // href is an external link
						link = strings.TrimSuffix(strings.TrimSpace(a.Val), "/")
						if !contains(link, s.visited) { // the link was not visited yet
							(*s.visited)[link] = External
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
				s.logger.Error(
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
				s.logger.Error(
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
			s.build(sitemap, child)
		}
	}
}
