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

	if len(os.Args) == 3 && sitemapFlag { // sitemap flag is set
		go animation()
		sitemap(root, logger) // build sitemap
		os.Exit(0)
	}

	if len(os.Args) == 3 && testFlag { // test flag is set
		// TODO: call test process
		// TODO: os.Exit(0)
	}
}
