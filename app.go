package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// String for each command
const SITEMAP = "sitemap"
const TEST = "test"
const SCREENSHOT = "screenshot"
const HELP = "help"

// Represents an instance of linkt.
type App struct {
	command string
	url     string
	options *Options
	logger  *slog.Logger
	JSON    []Record `json:"results"`
}

// Represents a record in the JSON file with test results.
type Record struct {
	URL    string `json:"url"`
	Status string `json:"status"`
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
	app := &App{command: command, options: options, url: url, JSON: []Record{}}
	app.logger = NewLogger(options.debug)
	return app
}

// Runs the app with its services.
func (app *App) Run() {
	switch app.command {
	case SITEMAP:
		app.Sitemap()
	case TEST:
		app.Test()
	case SCREENSHOT:
		app.Screenshot()
	case HELP:
		app.Help()
	default:
		switch {
		case app.options.version:
			app.Version()
		default:
			app.Help()
		}
	}
	os.Exit(0)
}

// Executes the sitemap command for linkt.
func (app *App) Sitemap() {
	helpMsg := ""
	switch {
	case app.options.print:
		root, err := url.Parse(strings.TrimSuffix(app.url, "/"))
		if err != nil || root.Scheme == "" || root.Host == "" {
			app.logger.Error("missing or invalid URL", "url", app.url, "error", err)
			os.Exit(0)
		}
		done := make(chan bool)
		if !app.options.debug {
			go app.Progress(done)
		}
		spider := NewSpider(app)
		sitemap := spider.Crawl(root)
		if !app.options.debug {
			done <- true
		}
		sitemap.Print()
		os.Exit(0)
	case app.options.xml:
		if app.options.directory == "" {
			helpMsg = "\nUsage: linkt --xml --dir <path> [options] sitemap <url>\n\n"
			helpMsg += "Options:\n"
			helpMsg += "\t--delay <milliseconds>\t\tThe amount of time to delay each HTTP request.\n"
			helpMsg += "\t--debug\t\t\tShow debug logs.\n\n"
			fmt.Print(helpMsg)
			os.Exit(0)
		}
		root, err := url.Parse(strings.TrimSuffix(app.url, "/"))
		if err != nil || root.Scheme == "" || root.Host == "" {
			app.logger.Error("missing or invalid URL", "url", app.url, "error", err)
			os.Exit(0)
		}
		done := make(chan bool)
		if !app.options.debug {
			go app.Progress(done)
		}
		spider := NewSpider(app)
		sitemap := spider.Crawl(root)
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
		helpMsg += "\t--delay <milliseconds>\t\tThe amount of time to delay each HTTP request.\n"
		helpMsg += "\t--debug\t\t\tShow debug logs.\n\n"
		fmt.Print(helpMsg)
		os.Exit(0)
	}
}

// Tests a site for broken links, namely links that return a 4xx or 5xx HTTP error.
func (app *App) Test() {
	helpMsg := ""
	var err error
	var file *os.File
	switch {
	case app.options.json:
		if app.options.directory == "" {
			helpMsg = "\nUsage: linkt --json --dir <path> [options] test <url>\n\n"
			helpMsg += "Options:\n"
			helpMsg += "\t--delay <milliseconds>\t\tThe amount of time to delay each HTTP request.\n"
			helpMsg += "\t--debug\t\t\tShow debug logs.\n\n"
			fmt.Print(helpMsg)
			os.Exit(0)
		} else {
			// create directory and file to store json file
			if err := os.MkdirAll(app.options.directory, os.ModePerm); err != nil {
				app.logger.Error("directory not found", "error", err)
				os.Exit(1)
			}
			processedURL := strings.ReplaceAll(app.url, "/", "-")
			processedURL = strings.ReplaceAll(processedURL, ":", "")
			filename := fmt.Sprintf("%s.json", processedURL)
			path := filepath.Join(app.options.directory, filename)
			file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				app.logger.Error(
					"error creating JSON file",
					"error", err,
					"filename", filename,
				)
				os.Exit(1)
			}
			defer file.Close()
		}
		root, err := url.Parse(strings.TrimSuffix(app.url, "/"))
		if err != nil || root.Scheme == "" || root.Host == "" {
			app.logger.Error("missing or invalid URL", "url", app.url, "error", err)
			os.Exit(0)
		}
		spider := NewSpider(app)
		spider.Crawl(root)
		if app.options.json {
			data, err := json.Marshal(app.JSON)
			if err != nil {
				app.logger.Error("error encoding the test results into JSON", "error", err)
				os.Exit(1)
			}
			file.WriteString(fmt.Sprintf(`{"root":"%s","results":`, root.String()))
			file.Write(data)
			file.WriteString("}")
		}
		os.Exit(0)

	default:
		root, err := url.Parse(strings.TrimSuffix(app.url, "/"))
		if err != nil || root.Scheme == "" || root.Host == "" {
			app.logger.Error("missing or invalid URL", "url", app.url, "error", err)
			os.Exit(0)
		}
		spider := NewSpider(app)
		spider.Crawl(root)
		os.Exit(0)
	}
}

// Takes screenshot of each page in a site and saves them to a directory.
func (app *App) Screenshot() {
	if app.options.directory == "" {
		helpMsg := "\nUsage: linkt --dir <path> [options] screenshot <url>\n\n"
		helpMsg += "Options:\n"
		helpMsg += "\t--delay <milliseconds>\t\tThe amount of time to delay each HTTP request.\n"
		helpMsg += "\t--debug\t\t\tShow debug logs.\n\n"
		fmt.Print(helpMsg)
		os.Exit(0)
	}
	// create directory to store screenshots
	if err := os.MkdirAll(app.options.directory, os.ModePerm); err != nil {
		app.logger.Error("directory not found", "error", err)
		os.Exit(1)
	}
	root, err := url.Parse(strings.TrimSuffix(app.url, "/"))
	if err != nil || root.Scheme == "" || root.Host == "" {
		app.logger.Error("missing or invalid URL", "url", app.url, "error", err)
		os.Exit(1)
	}
	done := make(chan bool)
	if !app.options.debug {
		go app.Progress(done)
	}
	spider := NewSpider(app)
	spider.Crawl(root)
	if !app.options.debug {
		done <- true
	}
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
	case SITEMAP:
		helpMsg = "\nUsage: linkt [options] sitemap <url>\n\n"
		helpMsg += "Options:\n"
		helpMsg += "\t--xml\t\t\t\tSave the sitemap to an XML file.\n"
		helpMsg += "\t--print\t\t\t\tPrint the sitemap to standard output.\n"
		helpMsg += "\t--dir <path>\t\t\tThe directory to store the XML file.\n"
		helpMsg += "\t--delay <milliseconds>\t\tThe amount of time to delay each HTTP request.\n"
		helpMsg += "\t--debug\t\t\t\tShow debug logs.\n\n"

	case TEST:
		helpMsg = "\nUsage: linkt [options] test <url>\n\n"
		helpMsg += "Options:\n"
		helpMsg += "\t--json\t\t\t\tSave the test results to a JSON file.\n"
		helpMsg += "\t--dir <path>\t\t\tThe directory to store the JSON file.\n"
		helpMsg += "\t--delay <milliseconds>\t\tThe amount of time to delay each HTTP request.\n"
		helpMsg += "\t--debug\t\t\t\tShow debug logs.\n\n"

	case SCREENSHOT:
		helpMsg = "\nUsage: linkt --dir <path> [options] screenshot <url>\n\n"
		helpMsg += "Options:\n"
		helpMsg += "\t--delay <milliseconds>\t\tThe amount of time to delay each HTTP request.\n"
		helpMsg += "\t--debug\t\t\t\tShow debug logs.\n\n"

	case HELP:
		fallthrough
	default:
		helpMsg = "\nUsage: linkt [options] <command> [<args>]\n\n"
		helpMsg += "Commands:\n"
		helpMsg += "\tsitemap\t\t\tBuild a sitemap with URL as the root.\n"
		helpMsg += "\ttest\t\t\tTest for broken links in anchor, image, link, and script tags.\n"
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

// Prints text to stdout in such a manner that it gives the impression the text is moving
// while the spider is working.
func (app *App) Progress(done chan bool) {
	for {
		switch app.command {
		case SITEMAP:
			select {
			case <-done:
				fmt.Printf(
					"\n%s[SUCCESS]%s sitemap was created!\n",
					Green, Reset)
				if app.options.xml {
					fmt.Printf(
						"\nsitemap was saved to %s%s/sitemap.xml%s\n\n",
						Green, app.options.directory, Reset)
				}
				return
			default:
				dots := []string{".  ", ".. ", "...", " ..", "  .", "   "}
				for _, s := range dots {
					fmt.Printf(
						"\r%s[PENDING]%s collecting links%s%s%s",
						Orange, Reset, Orange, s, Reset)
					time.Sleep((1 * time.Second) / 4)
				}
			}
		case SCREENSHOT:
			select {
			case <-done:
				fmt.Printf(
					"\n%s[SUCCESS]%s screenshots were taken!\n",
					Green, Reset)
				fmt.Printf(
					"\nscreenshots were saved to %s%s%s\n\n",
					Green, app.options.directory, Reset)
				return
			default:
				dots := []string{".  ", ".. ", "...", " ..", "  .", "   "}
				for _, s := range dots {
					fmt.Printf(
						"\r%s[PENDING]%s taking screenshots%s%s%s",
						Orange, Reset, Orange, s, Reset)
					time.Sleep((1 * time.Second) / 4)
				}
			}
		}
	}
}
