package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

func main() {
	// define custom usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage: linkt [options...] --url <url>\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -m, --sitemap\tBuild a sitemap.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -t, --test\t\tRun a test against the URL.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -s, --screenshot\tTake screenshots of a site.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -d, --debug\t\tShow debug logs.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -v, --version\tShow the version number.\n\n")
	}

	// setup and parse CLI flags
	var helpFlag,
		sitemapFlag,
		testFlag,
		screenshotFlag,
		versionFlag,
		debugFlag,
		linksFlag,
		imagesFlag bool
	var urlFlag string
	flag.BoolVar(&helpFlag, "help", false, "")
	flag.BoolVar(&helpFlag, "h", false, "")
	flag.BoolVar(&sitemapFlag, "sitemap", false, "")
	flag.BoolVar(&sitemapFlag, "m", false, "")
	flag.BoolVar(&testFlag, "test", false, "")
	flag.BoolVar(&testFlag, "t", false, "")
	flag.BoolVar(&testFlag, "screeshot", false, "")
	flag.BoolVar(&testFlag, "s", false, "")
	flag.BoolVar(&versionFlag, "version", false, "")
	flag.BoolVar(&versionFlag, "v", false, "")
	flag.BoolVar(&debugFlag, "debug", false, "")
	flag.BoolVar(&debugFlag, "d", false, "")
	flag.StringVar(&urlFlag, "url", "", "")
	flag.BoolVar(&linksFlag, "l", false, "")
	flag.BoolVar(&linksFlag, "links", false, "")
	flag.BoolVar(&imagesFlag, "i", false, "")
	flag.BoolVar(&imagesFlag, "images", false, "")
	flag.Parse()

	logger := NewLogger(debugFlag)

	switch {
	case sitemapFlag:
		root, err := url.Parse(urlFlag)
		if err != nil || root.Scheme == "" || root.Host == "" { // verify url
			logger.Error("missing or invalid URL", "url", urlFlag, "error", err)
			os.Exit(0)
		}
		done := make(chan bool)
		if !debugFlag { // start sitemap animation
			go sitemapAnimation(done)
		}
		// spawn a spider to build the sitemap
		spider := NewSpider(logger)
		sitemap := spider.DoSitemap(root)
		if !debugFlag { // stop sitemap animation
			done <- true
		}
		fmt.Printf("\n%s\n", sitemap.String())
		// exit the program successfully
		os.Exit(0)

	case testFlag:
		// define test flag help message
		flag.Usage = func() {
			fmt.Fprintf(flag.CommandLine.Output(), "\nUsage: linkt --test [options...] --url <url>\n\n")
			fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
			fmt.Fprintf(flag.CommandLine.Output(), "    -l, --links\t\tTest for broken links.\n")
			fmt.Fprintf(flag.CommandLine.Output(), "    -i, --images\tTest for missing images.\n\n")
		}

		switch {
		case linksFlag:
			fmt.Printf("%s[UNDER CONSTRUCTION]%s -l and --links is not available yet.\n\n", Orange, Reset)
			// TODO:

		case imagesFlag:
			fmt.Printf("%s[UNDER CONSTRUCTION]%s -i and --images is not available yet.\n\n", Orange, Reset)
			// TODO:

		default: // show help for test option
			flag.Usage()
			os.Exit(0)
		}

	case screenshotFlag:
		fmt.Printf("%s[UNDER CONSTRUCTION]%s -s and --screenshot is not available yet.\n\n", Orange, Reset)
		// TODO:

	case versionFlag: // show version
		logo := `
		
 ___       ___  ________   ___  __    _________
|\  \     |\  \|\   ___  \|\  \|\  \ |\___   ___\
\ \  \    \ \  \ \  \\ \  \ \  \/  /|\|___ \  \_|
 \ \  \    \ \  \ \  \\ \  \ \   ___  \   \ \  \
  \ \  \____\ \  \ \  \\ \  \ \  \\ \  \   \ \  \
   \ \_______\ \__\ \__\\ \__\ \__\\ \__\   \ \__\
    \|_______|\|__|\|__| \|__|\|__| \|__|    \|__|    v0.0.1, built with Go 1.23.2
                                                                     

		`
		fmt.Printf("%s\n", logo)
		os.Exit(0)

	default: // show help
		flag.Usage()
		os.Exit(0)
	}
}
