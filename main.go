package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	options := NewOptions()
	logger := NewLogger(options.debug)

	switch {
	case options.sitemap:
		switch {
		case options.print:
			root := isValidURL(options.url, logger)
			done := make(chan bool)
			if !options.debug { // start sitemap animation
				go sitemapAnimation(done)
			}
			// spawn a spider to build the sitemap
			spider := NewSpider(logger, options)
			sitemap := spider.BuildSitemap(root)
			if !options.debug { // stop sitemap animation
				done <- true
			}
			sitemap.Print()
			// exit the program successfully
			os.Exit(0)

		case options.xml:
			if options.directory == "" {
				flag.Usage()
				os.Exit(0)
			}
			root := isValidURL(options.url, logger)
			done := make(chan bool)
			if !options.debug { // start sitemap animation
				go sitemapAnimation(done)
			}
			// spawn a spider to build the sitemap
			spider := NewSpider(logger, options)
			sitemap := spider.BuildSitemap(root)
			sitemap.XML(options.directory)
			if !options.debug { // stop sitemap animation
				done <- true
			}

			// exit the program successfully
			os.Exit(0)

		default: // show help for sitemap options
			flag.Usage()
			os.Exit(0)
		}

	case options.test:
		switch {
		case options.links:
			root := isValidURL(options.url, logger)
			// spawn a spider to test for broken links
			spider := NewSpider(logger, options)
			spider.TestLinks(root)
			// exit the program successfully
			os.Exit(0)

		case options.images:
			fmt.Printf("%s[UNDER CONSTRUCTION]%s -i and --images is not available yet.\n\n", Orange, Reset)
			root := isValidURL(options.url, logger)
			// spawn a spider to test for broken links
			spider := NewSpider(logger, options)
			spider.TestImages(root)
			// exit the program successfully
			os.Exit(0)

		default: // show help for test options
			flag.Usage()
			os.Exit(0)
		}

	case options.screenshot:
		fmt.Printf("%s[UNDER CONSTRUCTION]%s -s and --screenshot is not available yet.\n\n", Orange, Reset)
		root := isValidURL(options.url, logger)
		// spawn a spider to test for broken links
		spider := NewSpider(logger, options)
		spider.TakeScreenshots(root)
		// exit the program successfully
		os.Exit(0)

	case options.version: // show version
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
