package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

// A spider with capabilities such as building a sitemap, testing links,
// and taking screenshots for a site.
type Spider struct {
	client  *http.Client
	app     *App
	visited *Set[string, int]
	current Page
}

// Returns a new spider with an HTTP client.
func NewSpider(app *App) *Spider {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &Spider{
		client:  c,
		app:     app,
		visited: &Set[string, int]{},
		current: Page{},
	}
}

// Enables the spider to crawl starting from the root URL.
func (spider *Spider) Crawl(root *url.URL) *Sitemap {
	sitemap := NewSitemap(spider.app.logger)
	page := *NewPage(root)
	page.kind = Internal
	page.parentURL = root.String()
	_, err := sitemap.AddRoot(page)
	if err != nil {
		spider.app.logger.Error(
			"error adding root page to the sitemap",
			"page", page.request.URL.String(),
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
	// fetch a page
	spider.current = node.GetElement()
	spider.fetch()

	// return early if node is external
	// we don't need to scrape anchor tags from an external node
	if spider.current.kind != Internal {
		return
	}

	// add root page as internal link to Set of links
	if reflect.DeepEqual(sitemap.Root(), spider.current) {
		(*spider.visited)[spider.current.request.URL.String()] = Internal
		spider.current.links[spider.current.request.URL.String()] = Internal
	}

	// parse page to get tree
	doc, err := html.Parse(spider.current.response.Body)
	if err != nil {
		spider.app.logger.Error(
			"error parsing a page",
			"page", spider.current.request.URL.String(),
			"error", err,
		)
		os.Exit(-1)
	}

	// collect each url on the current page
	spider.collect(doc)

	// populate the tree with Set of internal and external links
	for p, t := range spider.current.links {
		if t == Internal { // link is internal
			link, err := url.Parse(fmt.Sprintf("%s%s", sitemap.Root().GetElement().request.URL.String(), p))
			if err != nil {
				spider.app.logger.Error(
					"error parsing a page URL",
					"page", link,
					"error", err,
				)
			}
			page := *NewPage(link)
			page.kind = Internal
			page.parentURL = node.GetElement().request.URL.String()
			child := sitemap.AddChild(node, page)
			spider.walk(sitemap, child)

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
			page.kind = External
			page.parentURL = node.GetElement().request.URL.String()
			child := sitemap.AddChild(node, page)
			spider.walk(sitemap, child)
		}
	}
}

// collect is recursively called in the walk function to visit each anchor, img, or
// script tag on the spider.current.
func (spider *Spider) collect(n *html.Node) {
	switch spider.app.command {

	// sitemap and screenshot command collects links only from anchor tags
	case SITEMAP:
		fallthrough
	case SCREENSHOT:
		// node is an anchor tag
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr { // iterate tag attributes
				if a.Key == "href" { // attribute is an href
					spider.store(a)
					spider.app.logger.Info(
						"collected a page",
						"tag", n.Data,
						"attribute", a.Key,
						"page", a.Val,
					)
					break // skip the remaining attributes
				}
			}
		}

	// test command collects links from anchor, link, img, and script tags
	case TEST:
		// node is an anchor tag or a link tag
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "link") {
			for _, a := range n.Attr { // iterate tag attributes
				if a.Key == "href" { // attribute is an href
					spider.store(a)
					spider.app.logger.Info(
						"collected a page",
						"tag", n.Data,
						"attribute", a.Key,
						"page", a.Val,
					)
					break // skip the remaining attributes
				} else if a.Key == "data-href" { // attribute is a data-href
					spider.store(a)
					spider.app.logger.Info(
						"collected a page",
						"tag", n.Data,
						"attribute", a.Key,
						"page", a.Val,
					)
					break // skip the remaining attributes
				}
			}
		}
		// node is an img tag or script tag
		if n.Type == html.ElementNode && (n.Data == "img" || n.Data == "script") {
			for _, a := range n.Attr { // iterate tag attributes
				if a.Key == "src" { // attribute is a src
					spider.store(a)
					spider.app.logger.Info(
						"collected a page",
						"tag", n.Data,
						"attribute", a.Key,
						"page", a.Val,
					)
					break // skip the remaining attributes
				}
			}
		}
	}

	// visit each link on the current page
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		spider.collect(c)
	}
}

// The spider will store a link in temporary storage as it crawls.
func (spider *Spider) store(attr html.Attribute) {
	var link string
	if strings.HasPrefix(attr.Val, "/") { // link is internal
		link = strings.TrimSuffix(strings.TrimSpace(attr.Val), "/")
		if !spider.visited.Contains(link) { // the link was not visited yet
			(*spider.visited)[link] = Internal
			spider.current.links[link] = Internal // add internal link to Set of links
		}
	} else if !strings.HasPrefix(attr.Val, "#") { // link is external
		link = strings.TrimSuffix(strings.TrimSpace(attr.Val), "/")
		if !spider.visited.Contains(link) { // the link was not visited yet
			(*spider.visited)[link] = External
			spider.current.links[link] = External // add external link to Set of links
		}
	}
}

// Performs an HTTP request to get the current page.
func (spider *Spider) fetch() {
	// verify page URL contains valid URL
	url, err := url.Parse(spider.current.request.URL.String())
	if err != nil || url.Scheme == "" || url.Host == "" {
		spider.app.logger.Info("invalid URL", "url", url)
		return // skip the remaining code
	}
	// delay the http request
	delay := time.Duration(spider.app.options.delay) * time.Millisecond
	time.Sleep(delay)
	// start timer to get request time
	start := time.Now()
	spider.current.response, err = spider.client.Do(spider.current.request)
	spider.current.requestTime = fmt.Sprintf("%d ms", time.Since(start).Milliseconds())
	if err != nil {
		spider.app.logger.Error(
			"error getting the page",
			"page", spider.current.request.URL.String(),
			"error", err,
		)
		os.Exit(0)
	}
	spider.app.logger.Info(
		"fetched a page",
		"page", spider.current.request.URL.String(),
		"status", spider.current.response.Status,
		"request time", spider.current.requestTime,
	)
	// process the response
	spider.process()
}

// Performs an action based on the commands and options the spider received
// when the app was executed.
func (spider *Spider) process() {
	switch spider.app.command {
	case TEST:
		// print link test result to standard out
		switch status := spider.current.response.StatusCode; {
		case status >= 100 && status <= 199:
			fmt.Printf(
				"\n%s\n\tStatus\t\t\t%s%s%s\n\tRequest Time\t\t%s%s%s\n\tParent URL\t\t%s%s%s\n",
				spider.current.request.URL.String(),
				Blue,
				spider.current.response.Status,
				Reset,
				Faint,
				spider.current.requestTime,
				Reset,
				Faint,
				spider.current.parentURL,
				Reset,
			)
		case status >= 200 && status <= 299:
			fmt.Printf(
				"\n%s\n\tStatus\t\t\t%s%s%s\n\tRequest Time\t\t%s%s%s\n\tParent URL\t\t%s%s%s\n",
				spider.current.request.URL.String(),
				Green,
				spider.current.response.Status,
				Reset,
				Faint,
				spider.current.requestTime,
				Reset,
				Faint,
				spider.current.parentURL,
				Reset,
			)
		case status >= 300 && status <= 399:
			fmt.Printf(
				"\n%s\n\tStatus\t\t\t%s%s%s\n\tRequest Time\t\t%s%s%s\n\tParent URL\t\t%s%s%s\n",
				spider.current.request.URL.String(),
				Yellow,
				spider.current.response.Status,
				Reset,
				Faint,
				spider.current.requestTime,
				Reset,
				Faint,
				spider.current.parentURL,
				Reset,
			)
		case status >= 400 && status <= 499:
			fallthrough
		case status >= 500 && status <= 599:
			fmt.Printf(
				"\n%s\n\tStatus\t\t\t%s%s%s\n\tRequest Time\t\t%s%s%s\n\tParent URL\t\t%s%s%s\n",
				spider.current.request.URL.String(),
				Red,
				spider.current.response.Status,
				Reset,
				Faint,
				spider.current.requestTime,
				Reset,
				Faint,
				spider.current.parentURL,
				Reset,
			)
		case status == 999:
			fmt.Printf(
				"\n%s\n\tStatus\t\t\t%s%s%s\n\tRequest Time\t\t%s%s%s\n\tParent URL\t\t%s%s%s\n",
				spider.current.request.URL.String(),
				Purple,
				spider.current.response.Status,
				Reset,
				Faint,
				spider.current.requestTime,
				Reset,
				Faint,
				spider.current.parentURL,
				Reset,
			)
		}
		if spider.app.options.json {
			var status string
			if spider.current.response.StatusCode == 999 {
				status = "999 Request Denied"
			} else {
				status = spider.current.response.Status
			}
			r := NewRecord(
				spider.current.request.URL.String(),
				status,
				spider.current.requestTime,
			)
			spider.app.JSON = append(spider.app.JSON, r)
		}

	case SCREENSHOT:
		// create screenshot file
		processedURL := strings.ReplaceAll(spider.current.request.URL.String(), "/", "-")
		processedURL = strings.ReplaceAll(processedURL, ":", "")
		filename := fmt.Sprintf("%s.jpeg", processedURL)
		path := filepath.Join(spider.app.options.directory, filename)
		file, err := os.Create(path)
		if err != nil {
			spider.app.logger.Error(
				"error creating image file",
				"error", err,
				"filename", filename,
			)
			os.Exit(1)
		}
		defer file.Close()
		// navigate to page and take a screenshot
		ctx, cancel := chromedp.NewContext(
			context.Background(),
			// Uncomment to see browser UI (headless=false)
			// chromedp.WithDebugf(log.Printf),
		)
		defer cancel()
		var buf []byte
		err = chromedp.Run(ctx,
			chromedp.Navigate(spider.current.request.URL.String()),
			// Wait until page is fully loaded
			chromedp.WaitVisible("body", chromedp.ByQuery),
			// Take a screenshot of the entire page
			chromedp.FullScreenshot(&buf, 90),
		)
		if err != nil {
			spider.app.logger.Error(
				"error running chromedp",
				"error", err,
				"url", spider.current.request.URL.String(),
			)
			os.Exit(1)
		}
		// write the screenshot to file
		if _, err := file.Write(buf); err != nil {
			spider.app.logger.Error(
				"error writing data to image file",
				"error", err,
				"filename", filename,
			)
		}
	}
}
