package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"
)

func main() {
	// define custom usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage: linkt [options...] <url>\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -m, --sitemap\tBuild a sitemap.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -t, --test\t\tTest for broken links.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -s, --screenshot\tTake screenshots of a site.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "    -v, --version\tShow the version number.\n\n")
	}

	// setup and parse CLI flags
	var helpFlag, sitemapFlag, testFlag, screenshotFlag, versionFlag bool
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
	flag.Parse()

	switch {
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
		fmt.Printf("%s%s%s\n", Orange, logo, Reset)
		os.Exit(0)

	case sitemapFlag: // build sitemap
		var root *url.URL
		if len(os.Args) == 3 { // verify URL is valid
			if isValidURL(os.Args[1]) {
				root, _ = url.Parse(os.Args[1])
			} else if isValidURL(os.Args[2]) {
				root, _ = url.Parse(os.Args[2])
			} else { // invalid URL
				fmt.Printf("%sinvalid URL%s\n", Red, Reset)
				os.Exit(0)
			}
		}
		// create logs directory if it doesn't exist
		os.Mkdir("logs", 0777)

		// create logger
		now := time.Now()
		log, err := os.OpenFile(fmt.Sprintf("logs/%d.log", now.UnixNano()), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("%serror creating log: %s%s\n", Red, err, Reset)
			os.Exit(-1)
		}
		defer log.Close()
		logger := slog.New(slog.NewTextHandler(log, &slog.HandlerOptions{Level: slog.LevelDebug}))
		done := make(chan bool)
		go sitemapAnimation(done)
		spider := NewSpider(logger)
		spider.DoSitemap(root, done)
		os.Exit(0)

	case testFlag: // run test
	// TODO:

	case screenshotFlag: // take screenshots
	// TODO:

	default: // show help
		flag.Usage()
		os.Exit(0)
	}
}
