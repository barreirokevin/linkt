package main

import (
	"flag"
	"fmt"
)

// The values for the options available when executing linkt.
type Options struct {
	help       bool
	sitemap    bool
	test       bool
	screenshot bool
	version    bool
	debug      bool
	links      bool
	images     bool
	url        string
}

// Creates and returns Options which contains the values specified.
func NewOptions() *Options {
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
	options := &Options{}
	flag.BoolVar(&options.help, "help", false, "")
	flag.BoolVar(&options.help, "h", false, "")
	flag.BoolVar(&options.sitemap, "sitemap", false, "")
	flag.BoolVar(&options.sitemap, "m", false, "")
	flag.BoolVar(&options.test, "test", false, "")
	flag.BoolVar(&options.test, "t", false, "")
	flag.BoolVar(&options.screenshot, "screeshot", false, "")
	flag.BoolVar(&options.screenshot, "s", false, "")
	flag.BoolVar(&options.version, "version", false, "")
	flag.BoolVar(&options.version, "v", false, "")
	flag.BoolVar(&options.debug, "debug", false, "")
	flag.BoolVar(&options.debug, "d", false, "")
	flag.StringVar(&options.url, "url", "", "")
	flag.BoolVar(&options.links, "l", false, "")
	flag.BoolVar(&options.links, "links", false, "")
	flag.BoolVar(&options.images, "i", false, "")
	flag.BoolVar(&options.images, "images", false, "")
	flag.Parse()

	// define test flag help message
	if options.test {
		flag.Usage = func() {
			fmt.Fprintf(flag.CommandLine.Output(), "\nUsage: linkt --test [options...] --url <url>\n\n")
			fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
			fmt.Fprintf(flag.CommandLine.Output(), "    -l, --links\t\tTest for broken links.\n")
			fmt.Fprintf(flag.CommandLine.Output(), "    -i, --images\tTest for missing images.\n\n")
		}
	}

	return options
}
