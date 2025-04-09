package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	// define custom usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage: linkt [options...] <url>\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -s, --sitemap\tBuild a sitemap.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -t, --test\t\tTest for broken links.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -v, --version\tShow the version number.\n\n")
	}

	// setup and parse CLI flags
	var helpFlag, sitemapFlag, testFlag, versionFlag bool
	flag.BoolVar(&helpFlag, "help", false, "")
	flag.BoolVar(&helpFlag, "h", false, "")
	flag.BoolVar(&sitemapFlag, "sitemap", false, "")
	flag.BoolVar(&sitemapFlag, "s", false, "")
	flag.BoolVar(&testFlag, "test", false, "")
	flag.BoolVar(&testFlag, "t", false, "")
	flag.BoolVar(&versionFlag, "version", false, "")
	flag.BoolVar(&versionFlag, "v", false, "")
	flag.Parse()

	if len(os.Args) < 2 || helpFlag { // help flag is set
		flag.Usage()
		os.Exit(0)
	}

	if len(os.Args) == 2 && versionFlag { // version flag is set
		fmt.Printf("%slinkt v0.0.1%s\n", Orange, Reset)
		os.Exit(0)
	}

	var root *url.URL
	if len(os.Args) == 3 { // verify URL is valid
		if isValidURL(os.Args[1]) {
			root, _ = url.Parse(os.Args[1])
		} else if isValidURL(os.Args[2]) {
			root, _ = url.Parse(os.Args[2])
		} else { // invalid URL
			fmt.Printf("%sInvalid URL%s\n", Red, Reset)
			os.Exit(0)
		}
	}

	if len(os.Args) == 3 && sitemapFlag { // sitemap flag is set
		sitemap(root) // build sitemap
		os.Exit(0)
	}

	if len(os.Args) == 3 && testFlag { // test flag is set
		// TODO: call test process
		// TODO: os.Exit(0)
	}
}

func isValidURL(value string) bool {
	url, err := url.Parse(value)
	if err != nil {
		return false
	}
	if url.Scheme == "" || url.Host == "" {
		return false
	}
	return true
}

func sitemap(root *url.URL) {
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
		slog.Error(
			"error adding root page to the sitemap",
			"page", sitemap.Root().GetElement().Request.URL.String(),
			"error", err,
		)
		os.Exit(-1)
	}
	(*visited)[page.Request.URL.String()] = Internal
	// construct http client
	client := &http.Client{}
	// call spider to start crawling from the root
	spider(client, sitemap, sitemap.Root(), visited)
	// print the sitemap
	printPreorderIndent(sitemap, sitemap.Root(), 0)
}

func spider(client *http.Client, sitemap *Tree[Page], page *Node[Page], visited *Set) {
	// get page
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
						if !wasVisited(link, visited) { // the link was not visited yet
							(*visited)[link] = Internal
							page.GetElement().Links[link] = Internal // add internal link to Set of links
						}
					} else if !strings.HasPrefix(a.Val, "#") { // href is an external link
						link = strings.TrimSuffix(strings.TrimSpace(a.Val), "/")
						if !wasVisited(link, visited) { // the link was not visited yet
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
			spider(client, sitemap, child, visited)
		}
	}
}

func wasVisited(link string, visited *Set) bool {
	for k, _ := range *visited {
		if strings.ToLower(k) == strings.ToLower(link) {
			return true
		}
	}
	return false
}

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Orange = "\033[38;5;215m"
)

func printPreorderIndent(t *Tree[Page], n *Node[Page], d int) {
	var color string
	if len(n.Children()) > 0 {
		color = Cyan
	} else {
		color = Purple
	}

	indent := strings.Repeat(" ", d*4)
	fmt.Printf("%s%s%+v%s\n", color, indent, n.GetElement().Request.URL.String(), Reset)
	for _, c := range n.Children() {
		printPreorderIndent(t, c, d+1)
	}
}
