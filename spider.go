package main

import (
	"fmt"
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
	app     *App
	visited *Set[string, int]
}

// Returns a new spider with an HTTP client.
func NewSpider(app *App) *Spider {
	c := &http.Client{}
	return &Spider{client: c, app: app, visited: &Set[string, int]{}}
}

// Enables the spider to build a sitemap starting from the root URL.
func (spider *Spider) BuildSitemap(root *url.URL) *Sitemap {
	sitemap := NewSitemap(spider.app.logger)
	page := *NewPage(root)
	_, err := sitemap.AddRoot(page)
	if err != nil {
		spider.app.logger.Error(
			"error adding root page to the sitemap",
			"page", page.URL(),
			"error", err,
		)
		os.Exit(-1)
	}
	// start recursively building the sitemap
	spider.walk(sitemap, sitemap.Root())
	return sitemap
}

// Enables the spider to test a sitemap and report on the HTTP status code for each link.
func (spider *Spider) TestLinks(root *url.URL) {
	sitemap := NewSitemap(spider.app.logger)
	page := *NewPage(root)
	_, err := sitemap.AddRoot(page)
	if err != nil {
		spider.app.logger.Error(
			"error adding root page to the sitemap",
			"page", page.URL(),
			"error", err,
		)
		os.Exit(-1)
	}
	// start recursively building the sitemap
	spider.walk(sitemap, sitemap.Root())
}

// TODO:
func (spider *Spider) TestImages(root *url.URL) {}

// TODO:
func (spider *Spider) TakeScreenshots(root *url.URL) {}

// Enables the spider to recursively walk through the elements on a page. As the spider
// walks through the elements on the page it builds a sitemap with the anchor tags it
// encounters. The spider reports its work based on the action specified.
func (spider *Spider) walk(sitemap *Sitemap, node *Node[Page]) {
	// get page
	currentPage := node.GetElement()
	resp, err := spider.client.Do(currentPage.Request())
	if err != nil {
		spider.app.logger.Error(
			"error getting the page",
			"page", currentPage.URL(),
			"error", err,
		)
		os.Exit(0)
	}
	spider.app.logger.Info(
		"got the page",
		"page", currentPage.URL(),
		"status", resp.Status,
	)

	// display test results
	if spider.app.command == "test" && spider.app.options.links {
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
		(*spider.visited)[currentPage.URL()] = Internal
		currentPage.Links()[currentPage.URL()] = Internal
	}

	// parse page to get tree
	doc, err := html.Parse(resp.Body)
	if err != nil {
		spider.app.logger.Error(
			"error parsing a page",
			"page", currentPage.URL(),
			"error", err,
		)
		os.Exit(-1)
	}

	// step through the page
	spider.step(doc, &currentPage)

	// populate the tree with Set of internal and external links
	for p, t := range currentPage.Links() {
		if t == Internal { // link is internal
			link, err := url.Parse(fmt.Sprintf("%s%s", sitemap.Root().GetElement().URL(), p))
			if err != nil {
				spider.app.logger.Error(
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
				spider.app.logger.Error(
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

	for _, child := range node.Children() {
		if spider.app.command == "sitemap" && child.GetElement().Type() == Internal {
			// walk on an internal link if building a sitemap
			spider.walk(sitemap, child)
		} else if spider.app.command == "test" && spider.app.options.links {
			// walk on an internal or external link if testing for broken links
			spider.walk(sitemap, child)
		}
	}
}

// Step is recursively called in the walk function to visit each anchor tag on the currentPage.
func (spider *Spider) step(n *html.Node, currentPage *Page) {
	var link string
	if n.Type == html.ElementNode && n.Data == "a" { // node is an anchor tag
		for _, a := range n.Attr { // iterate anchor tag attributes
			if a.Key == "href" { // attribute is an href
				if strings.HasPrefix(a.Val, "/") { // href is an internal link
					link = strings.TrimSuffix(strings.TrimSpace(a.Val), "/")
					if !spider.visited.Contains(link) { // the link was not visited yet
						(*spider.visited)[link] = Internal
						currentPage.Links()[link] = Internal // add internal link to Set of links
					}
				} else if !strings.HasPrefix(a.Val, "#") { // href is an external link
					link = strings.TrimSuffix(strings.TrimSpace(a.Val), "/")
					if !spider.visited.Contains(link) { // the link was not visited yet
						(*spider.visited)[link] = External
						currentPage.Links()[link] = External
					}
				}
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		spider.step(c, currentPage)
	}
}
