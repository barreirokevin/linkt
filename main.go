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
		fmt.Fprintf(flag.CommandLine.Output(), "    -t, --test\t\tTest for broken links.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -s, --screenshot\tTake screenshots of a site.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -d, --debug\t\tShow debug logs.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -v, --version\tShow the version number.\n\n")
	}

	// setup and parse CLI flags
	var helpFlag, sitemapFlag, testFlag, screenshotFlag, versionFlag, debugFlag bool
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
	// TODO:

	case screenshotFlag:
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
