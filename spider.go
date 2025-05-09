package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

// A spider with capabilities such as building a sitemap, testing links,
// and taking screenshots for a site.
type Spider struct {
	client  *http.Client
	app     *App
	visited *Set[string, int]
	current struct {
		page     Page
		response *http.Response
	}
}

// Returns a new spider with an HTTP client.
func NewSpider(app *App) *Spider {
	c := &http.Client{}
	return &Spider{
		client:  c,
		app:     app,
		visited: &Set[string, int]{},
		current: struct {
			page     Page
			response *http.Response
		}{},
	}
}

// Enables the spider to crawl starting from the root URL.
func (spider *Spider) Crawl(root *url.URL) *Sitemap {
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

// Enables the spider to recursively walk through the elements on a page. As the spider
// walks through the elements on the page it builds a sitemap with the anchor tags it
// encounters. The spider reports its work based on the action specified.
func (spider *Spider) walk(sitemap *Sitemap, node *Node[Page]) {
	// get page
	spider.current.page = node.GetElement()
	var err error
	spider.current.response, err = spider.client.Do(spider.current.page.Request())
	if err != nil {
		spider.app.logger.Error(
			"error getting the page",
			"page", spider.current.page.URL(),
			"error", err,
		)
		os.Exit(0)
	}
	spider.app.logger.Info(
		"got the page",
		"page", spider.current.page.URL(),
		"status", spider.current.response.Status,
	)

	// execute command on page
	spider.execute()

	// add root page as internal link to Set of links
	if reflect.DeepEqual(sitemap.Root(), spider.current.page) {
		(*spider.visited)[spider.current.page.URL()] = Internal
		spider.current.page.Links()[spider.current.page.URL()] = Internal
	}

	// parse page to get tree
	doc, err := html.Parse(spider.current.response.Body)
	if err != nil {
		spider.app.logger.Error(
			"error parsing a page",
			"page", spider.current.page.URL(),
			"error", err,
		)
		os.Exit(-1)
	}

	// step through the page
	spider.step(doc)

	// populate the tree with Set of internal and external links
	for p, t := range spider.current.page.Links() {
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

// Step is recursively called in the walk function to visit each anchor tag on the spider.current.page.
func (spider *Spider) step(n *html.Node) {
	var link string
	if n.Type == html.ElementNode && n.Data == "a" { // node is an anchor tag
		for _, a := range n.Attr { // iterate anchor tag attributes
			if a.Key == "href" { // attribute is an href
				if strings.HasPrefix(a.Val, "/") { // href is an internal link
					link = strings.TrimSuffix(strings.TrimSpace(a.Val), "/")
					if !spider.visited.Contains(link) { // the link was not visited yet
						(*spider.visited)[link] = Internal
						spider.current.page.Links()[link] = Internal // add internal link to Set of links
					}
				} else if !strings.HasPrefix(a.Val, "#") { // href is an external link
					link = strings.TrimSuffix(strings.TrimSpace(a.Val), "/")
					if !spider.visited.Contains(link) { // the link was not visited yet
						(*spider.visited)[link] = External
						spider.current.page.Links()[link] = External
					}
				}
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		spider.step(c)
	}
}

// Performs an action based on the commands and options the spider received
// when the app was executed.
func (spider *Spider) execute() {
	switch spider.app.command {
	case "test":
		switch {
		case spider.app.options.links:
			switch status := spider.current.response.StatusCode; {
			case status >= 100 && status <= 199:
				fmt.Printf(
					"\t%s[%s]%s\t%s\n",
					Blue,
					spider.current.response.Status,
					Reset,
					spider.current.page.URL())
			case status >= 200 && status <= 299:
				fmt.Printf(
					"\t%s[%s]%s\t%s\n",
					Green,
					spider.current.response.Status,
					Reset,
					spider.current.page.URL())
			case status >= 300 && status <= 399:
				fmt.Printf(
					"\t%s[%s]%s\t%s\n",
					Yellow,
					spider.current.response.Status,
					Reset,
					spider.current.page.URL())
			case status >= 400 && status <= 499:
				fmt.Printf(
					"\t%s[%s]%s\t%s\n",
					Red,
					spider.current.response.Status,
					Reset,
					spider.current.page.URL())
			case status >= 500 && status <= 599:
				fmt.Printf(
					"\t%s[%s]%s\t%s\n",
					Red,
					spider.current.response.Status,
					Reset,
					spider.current.page.URL())
			}

		case spider.app.options.images:
			// TODO:

		default:
			// TODO:

		}

	case "screenshot":
		ctx, cancel := chromedp.NewContext(
			context.Background(),
			// Uncomment to see browser UI (headless=false)
			// chromedp.WithDebugf(log.Printf),
		)
		defer cancel()
		// Create screenshot.png file
		screenshotFile := "test_screenshot.jpeg"
		file, err := os.Create(screenshotFile)
		if err != nil {
			spider.app.logger.Error("error creating image file", "error", err)
			os.Exit(1)
		}
		defer file.Close()
		// Navigate to website and take a screenshot
		var buf []byte
		err = chromedp.Run(ctx,
			chromedp.Navigate("https://www.google.com"),
			// Wait until page is fully loaded
			chromedp.WaitVisible("body", chromedp.ByQuery),
			// Take a screenshot of the entire page
			chromedp.FullScreenshot(&buf, 90),
		)
		if err != nil {
			spider.app.logger.Error("error running chromedp", "error", err)
			os.Exit(1)
		}
		// Write the screenshot to file
		if _, err := file.Write(buf); err != nil {
			spider.app.logger.Error("error writing data to image file", "error", err)
		}
	}
}
