package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
)

// Represents an instance of linkt.
type App struct {
	command string
	url     string
	options *Options
	logger  *slog.Logger
}

// Creates and returns a new app with the services needed to run it.
func NewApp() *App {
	command := ""
	url := ""
	if len(os.Args) > 1 {
		// command is 2nd to last element
		command = os.Args[len(os.Args)-2]
		command = strings.ToLower(command)
		command = strings.TrimSpace(command)
		// url is last element
		url = os.Args[len(os.Args)-1]
		url = strings.ToLower(url)
		url = strings.TrimSpace(url)
	}
	options := NewOptions()
	app := &App{command: command, options: options, url: url}
	app.logger = NewLogger(options.debug)
	return app
}

// Runs the app with its services.
func (app *App) Run() {
	switch app.command {
	case "sitemap":
		app.Sitemap()
	case "test":
		app.Test()
	case "screenshot":
		app.Screenshot()
	case "help":
		app.Help()
	default:
		switch {
		case app.options.version:
			app.Version()
		default:
			helpMsg := "\nUsage: linkt [options] <command> [<args>]\n\n"
			helpMsg += "Commands:\n"
			helpMsg += "\tsitemap\t\t\tBuild a sitemap with URL as the root.\n"
			helpMsg += "\ttest\t\t\tRun a test against the URL.\n"
			helpMsg += "\tscreenshot\t\tTake screenshots of all the pages on a site.\n"
			helpMsg += "\thelp <command>\t\tDisplay help for a command.\n\n"
			helpMsg += "Options:\n"
			helpMsg += "\t-v, --version\t\tShow the version number.\n\n"
			fmt.Print(helpMsg)
		}
	}
	os.Exit(0)
}

// Executes the sitemap command for linkt.
func (app *App) Sitemap() {
	helpMsg := ""
	switch {
	case app.options.print:
		root, err := url.Parse(app.url)
		if err != nil || root.Scheme == "" || root.Host == "" {
			app.logger.Error("missing or invalid URL", "url", app.url, "error", err)
			os.Exit(0)
		}
		done := make(chan bool)
		if !app.options.debug {
			go sitemapAnimation(done)
		}
		spider := NewSpider(app)
		sitemap := spider.BuildSitemap(root)
		if !app.options.debug {
			done <- true
		}
		sitemap.Print()
		os.Exit(0)
	case app.options.xml:
		if app.options.directory == "" {
			helpMsg = "\nUsage: linkt --xml --dir <path> [options] sitemap <url>\n\n"
			helpMsg += "Options:\n"
			helpMsg += "\t--dir <path>\t\tThe directory to store the XML file.\n"
			helpMsg += "\t--debug\t\t\tShow debug logs.\n\n"
			fmt.Print(helpMsg)
			os.Exit(0)
		}
		root, err := url.Parse(app.url)
		if err != nil || root.Scheme == "" || root.Host == "" {
			app.logger.Error("missing or invalid URL", "url", app.url, "error", err)
			os.Exit(0)
		}
		done := make(chan bool)
		if !app.options.debug {
			go sitemapAnimation(done)
		}
		spider := NewSpider(app)
		sitemap := spider.BuildSitemap(root)
		sitemap.XML(app.options.directory)
		if !app.options.debug {
			done <- true
		}
		os.Exit(0)
	default:
		helpMsg = "\nUsage: linkt [options] sitemap <url>\n\n"
		helpMsg += "Options:\n"
		helpMsg += "\t--xml\t\t\tSave the sitemap to an XML file.\n"
		helpMsg += "\t--print\t\t\tPrint the sitemap to standard output.\n"
		helpMsg += "\t--dir <path>\t\tThe directory to store the XML file.\n"
		helpMsg += "\t--debug\t\t\tShow debug logs.\n\n"
		fmt.Print(helpMsg)
		os.Exit(0)
	}
}

// Tests a site for broken links, namely links that return a 4xx or 5xx HTTP error.
func (app *App) Test() {
	switch {
	case app.options.links:
		root, err := url.Parse(app.url)
		if err != nil || root.Scheme == "" || root.Host == "" {
			app.logger.Error("missing or invalid URL", "url", app.url, "error", err)
			os.Exit(0)
		}
		spider := NewSpider(app)
		spider.TestLinks(root)
		os.Exit(0)

	case app.options.images:
		fmt.Printf("%s[UNDER CONSTRUCTION]%s -i and --images is not available yet.\n\n", Orange, Reset)
		root, err := url.Parse(app.url)
		if err != nil || root.Scheme == "" || root.Host == "" {
			app.logger.Error("missing or invalid URL", "url", app.url, "error", err)
			os.Exit(0)
		}
		spider := NewSpider(app)
		spider.TestImages(root)
		os.Exit(0)

	default:
		helpMsg := "\bUsage: linkt [options] test <url>\n\n"
		helpMsg += "Options:\n"
		helpMsg += "\t-l, --links\t\tTest for broken links.\n"
		helpMsg += "\t-i, --images\t\tTest for missing images.\n"
		helpMsg += "\t--debug\t\t\tShow debug logs.\n\n"
		fmt.Print(helpMsg)
		os.Exit(0)
	}
}

// Takes screenshot of each page in a site and saves them to a directory.
func (app *App) Screenshot() {
	if app.options.directory == "" {
		helpMsg := "\nUsage: linkt --dir <path> [options] screenshot <url>\n\n"
		helpMsg += "Options:\n"
		helpMsg += "\t--debug\t\t\tShow debug logs.\n\n"
		fmt.Print(helpMsg)
		os.Exit(0)
	}
	fmt.Printf("%s[UNDER CONSTRUCTION]%s -s and --screenshot is not available yet.\n\n", Orange, Reset)
	root, err := url.Parse(app.url)
	if err != nil || root.Scheme == "" || root.Host == "" {
		app.logger.Error("missing or invalid URL", "url", app.url, "error", err)
		os.Exit(0)
	}
	spider := NewSpider(app)
	spider.TakeScreenshots(root)
	os.Exit(0)
}

// Prints the help message for a corresponding command or option to standard output.
func (app *App) Help() {
	var helpCmd string
	if len(os.Args) > 2 {
		helpCmd = os.Args[2]
		helpCmd = strings.ToLower(helpCmd)
		helpCmd = strings.TrimSpace(helpCmd)
	}

	var helpMsg string
	switch helpCmd {
	case "sitemap":
		helpMsg = "\nUsage: linkt [options] sitemap <url>\n\n"
		helpMsg += "Options:\n"
		helpMsg += "\t--xml\t\t\tSave the sitemap to an XML file.\n"
		helpMsg += "\t--print\t\t\tPrint the sitemap to standard output.\n"
		helpMsg += "\t--dir <path>\t\tThe directory to store the XML file.\n"
		helpMsg += "\t--debug\t\t\tShow debug logs.\n\n"

	case "test":
		helpMsg = "\nUsage: linkt [options] test <url>\n\n"
		helpMsg += "Options:\n"
		helpMsg += "\t-l, --links\t\tTest for broken links.\n"
		helpMsg += "\t-i, --images\t\tTest for missing images.\n"
		helpMsg += "\t--debug\t\t\tShow debug logs.\n\n"

	case "screenshot":
		helpMsg = "\nUsage: linkt --dir <path> [options] screenshot <url>\n\n"
		helpMsg += "Options:\n"
		helpMsg += "\t--debug\t\t\tShow debug logs.\n\n"

	case "help":
		fallthrough
	default:
		helpMsg = "\nUsage: linkt [options] <command> [<args>]\n\n"
		helpMsg += "Commands:\n"
		helpMsg += "\tsitemap\t\t\tBuild a sitemap with URL as the root.\n"
		helpMsg += "\ttest\t\t\tRun a test against the URL.\n"
		helpMsg += "\tscreenshot\t\tTake screenshots of all the pages on a site.\n"
		helpMsg += "\thelp <command>\t\tDisplay help for a command.\n\n"
		helpMsg += "Options:\n"
		helpMsg += "\t-v, --version\t\tShow the version number.\n\n"
	}
	fmt.Printf(helpMsg)
}

// Prints the linkt version and logo to standard output.
func (app *App) Version() {
	version := `
		
 ___       ___  ________   ___  __    _________
|\  \     |\  \|\   ___  \|\  \|\  \ |\___   ___\
\ \  \    \ \  \ \  \\ \  \ \  \/  /|\|___ \  \_|
 \ \  \    \ \  \ \  \\ \  \ \   ___  \   \ \  \
  \ \  \____\ \  \ \  \\ \  \ \  \\ \  \   \ \  \
   \ \_______\ \__\ \__\\ \__\ \__\\ \__\   \ \__\
    \|_______|\|__|\|__| \|__|\|__| \|__|    \|__|    v0.0.1, built with Go 1.23.2
                                                                     

		`
	fmt.Printf("%s\n", version)
}
