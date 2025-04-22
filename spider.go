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

// A spider with capabilities such as building a sitemap, testing links,
// and taking screenshots for a site.
type Spider struct {
	client  *http.Client
	logger  *slog.Logger
	visited *Set
	options *Options
}

// Returns a new spider with an HTTP client.
func NewSpider(logger *slog.Logger, options *Options) *Spider {
	c := &http.Client{}
	return &Spider{client: c, logger: logger, visited: &Set{}, options: options}
}

// Enables the spider to build a sitemap starting from the root URL.
func (s *Spider) BuildSitemap(root *url.URL) *Sitemap {
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
	s.walk(sitemap, sitemap.Root())
	return sitemap
}

// Enables the spider to test a sitemap and report on the HTTP status code for each link.
func (s *Spider) TestLinks(root *url.URL) {
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
	s.walk(sitemap, sitemap.Root())
}

// TODO:
func (s *Spider) TestImages(root *url.URL) {}

// TODO:
func (s *Spider) TakeScreenshots(root *url.URL) {}

// Enables the spider to recursively walk through the elements on a page. As the spider
// walks through the elements on the page it builds a sitemap with the anchor tags it
// encounters. The spider reports its work based on the action specified.
func (s *Spider) walk(sitemap *Sitemap, node *Node[Page]) {
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

	// display test results
	if s.options.test && s.options.links {
		switch status := resp.StatusCode; {
		case status >= 100 && status <= 199:
			fmt.Printf("\t%s[%s]%s\t%s\n", Blue, resp.Status, Reset, currentPage.URL())
		case status >= 200 && status <= 299:
			fmt.Printf("\t%s[%s]%s\t%s\n", Green, resp.Status, Reset, currentPage.URL())
		case status >= 300 && status <= 399:
			fmt.Printf("\t%s[%s]%s\t%s\n", Yellow, resp.Status, Reset, currentPage.URL())
		case status >= 400 && status <= 499:
			fmt.Printf("\t%s[%s]%s\t%s\n", Red, resp.Status, Reset, currentPage.URL())
		case status >= 500 && status <= 599:
			fmt.Printf("\t%s[%s]%s\t%s\n", Red, resp.Status, Reset, currentPage.URL())
		}
	}

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

	// step through the page
	s.step(doc, &currentPage)

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

	// walk on an internal link
	for _, child := range node.Children() {
		if child.GetElement().Type() == Internal {
			s.walk(sitemap, child)
		}
	}
}

// Step is recursively called in the walk function to visit each anchor tag on the currentPage.
func (s *Spider) step(n *html.Node, currentPage *Page) {
	var link string
	if n.Type == html.ElementNode && n.Data == "a" { // node is an anchor tag
		for _, a := range n.Attr { // iterate anchor tag attributes
			if a.Key == "href" { // attribute is an href
				if strings.HasPrefix(a.Val, "/") { // href is an internal link
					link = strings.TrimSuffix(strings.TrimSpace(a.Val), "/")
					if !s.visited.Contains(link) { // the link was not visited yet
						(*s.visited)[link] = Internal
						currentPage.Links()[link] = Internal // add internal link to Set of links
					}
				} else if !strings.HasPrefix(a.Val, "#") { // href is an external link
					link = strings.TrimSuffix(strings.TrimSpace(a.Val), "/")
					if !s.visited.Contains(link) { // the link was not visited yet
						(*s.visited)[link] = External
						currentPage.Links()[link] = External
					}
				}
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s.step(c, currentPage)
	}
}
