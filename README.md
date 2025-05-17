```
 ___       ___  ________   ___  __    _________
|\  \     |\  \|\   ___  \|\  \|\  \ |\___   ___\
\ \  \    \ \  \ \  \\ \  \ \  \/  /|\|___ \  \_|
 \ \  \    \ \  \ \  \\ \  \ \   ___  \   \ \  \
  \ \  \____\ \  \ \  \\ \  \ \  \\ \  \   \ \  \
   \ \_______\ \__\ \__\\ \__\ \__\\ \__\   \ \__\
    \|_______|\|__|\|__| \|__|\|__| \|__|    \|__|
```

linkt is a command-line tool to perform a variety of actions with a URL, such as building a sitemap, testing for broken links, and taking screenshots of each page.

## Usage

```
Usage: linkt [options] <command> [<args>]

Commands:
        sitemap                 Build a sitemap with URL as the root.
        test                    Test for broken links in anchor, image, link, and script tags.
        screenshot              Take screenshots of all the pages on a site.
        help <command>          Display help for a command.

Options:
        -v, --version           Show the version number.
```

## Install

1. Download the latest source code:

   ```
   git clone https://github.com/barreirokevin/linkt.git
   ```

1. Navigate to the linkt directroy:

   ```
   cd <path-to-linkt>
   ```

1. Build the source code:

   ```
   go build -o bin/linkt *.go
   ```

1. Execute linkt:

   ```
   ./bin/linkt
   ```
