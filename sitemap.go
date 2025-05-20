package main

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Sitemap struct {
	Tree[Page]
	logger *slog.Logger
}

// Returns an empty sitemap.
func NewSitemap(logger *slog.Logger) *Sitemap {
	return &Sitemap{logger: logger}
}

// Returns the tree as a string that displays the hiearachy.
func (s *Sitemap) String() string {
	str := ""

	var preorderIndent func(s *Sitemap, n *Node[Page], d int)
	preorderIndent = func(s *Sitemap, n *Node[Page], d int) {
		if reflect.DeepEqual(s.Root(), n) {
			// the current node is the root
			str += fmt.Sprintf("%s\n", n.GetElement().URL())

		} else if len(n.Children()) == 0 &&
			reflect.DeepEqual(n.GetParent().Children()[len(n.GetParent().Children())-1], n) {
			// the current node is the last child
			indent := strings.Repeat(" ", d*4)
			if d > 0 && d%2 == 0 {
				str += fmt.Sprintf("│%s └─── %+v\n", indent, n.GetElement().URL())
			} else if d > 0 && d%2 != 0 {
				str += fmt.Sprintf("│%s└─── %+v\n", indent, n.GetElement().URL())
			} else {
				str += fmt.Sprintf("%s└─── %+v\n", indent, n.GetElement().URL())
			}

		} else {
			// the current node is not the last child
			indent := strings.Repeat(" ", d*4)
			if d > 0 && d%2 == 0 {
				str += fmt.Sprintf("│%s ├─── %+v\n", indent, n.GetElement().URL())
			} else if d > 0 && d%2 != 0 {
				str += fmt.Sprintf("│%s├─── %+v\n", indent, n.GetElement().URL())
			} else {
				str += fmt.Sprintf("%s├─── %+v\n", indent, n.GetElement().URL())
			}
		}

		// recursively call preorderIndent for each child
		for _, c := range n.Children() {
			preorderIndent(s, c, d+1)
		}
	}

	// begin preorderIndent with the tree's root
	preorderIndent(s, s.Root(), -1)
	return str
}

// Prints the sitemap to standard out.
func (s *Sitemap) Print() {
	fmt.Printf("\n%s\n", s.String())
}

// Writes each link in the sitemap to an XML file and stores that file at directory dir.
func (s *Sitemap) XML(dir string) {
	// create sitemap XML file
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		s.logger.Error("directory not found", "error", err)
		os.Exit(0)
	}
	path := filepath.Join(dir, "sitemap.xml")
	file, err := os.Create(path)
	if err != nil {
		s.logger.Error("sitemap file not created", "error", err)
		os.Exit(0)
	}
	defer file.Close()

	// create Set of all links in the sitemap
	type url struct {
		Link string `xml:"loc"`
	}
	allLinks := Set[url, int]{} // Set prevents duplicate links
	// write each entry to the XML file
	file.WriteString(xml.Header)
	file.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	file.WriteString("\n")
	s.Preorder(func(c *Node[Page]) {
		e := url{Link: c.GetElement().URL()}
		if !allLinks.Contains(e) {
			allLinks[e] = 0
			data, err := xml.MarshalIndent(e, "", "  ")
			if err != nil {
				s.logger.Error("could not marshal link to XML", "error", err)
				os.Exit(0)
			}
			file.Write(data)
			file.WriteString("\n")
		}
	})
	file.WriteString("</urlset>")
}
