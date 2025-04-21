```
 ___       ___  ________   ___  __    _________
|\  \     |\  \|\   ___  \|\  \|\  \ |\___   ___\
\ \  \    \ \  \ \  \\ \  \ \  \/  /|\|___ \  \_|
 \ \  \    \ \  \ \  \\ \  \ \   ___  \   \ \  \
  \ \  \____\ \  \ \  \\ \  \ \  \\ \  \   \ \  \
   \ \_______\ \__\ \__\\ \__\ \__\\ \__\   \ \__\
    \|_______|\|__|\|__| \|__|\|__| \|__|    \|__|
```

linkt is a command-line tool to perform a variety of actions with a URL, such as building a sitemap, testing for broken links, testing for missing images, and taking screenshots of each page.

## Usage

```bash
Usage: linkt [options...] --url <url>

Options:
    -m, --sitemap       Build a sitemap.
    -t, --test          Test for broken links.
    -s, --screenshot    Take screenshots of a site.
    -d, --debug         Show debug logs.
    -v, --version       Show the version number.
```

## Install

1. Download the latest source code:

   ```bash
   git clone https://github.com/barreirokevin/linkt.git
   ```

1. Navigate to the linkt directroy:

   ```bash
   cd <path-to-linkt>
   ```

1. Build the source code:

   ```bash
   go build -o bin/linkt *.go
   ```

1. Execute linkt:

   ```bash
   ./bin/linkt
   ```
