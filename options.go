package main

import "flag"

// The values for the options available when executing linkt.
type Options struct {
	version   bool
	debug     bool
	links     bool
	images    bool
	xml       bool
	print     bool
	directory string
}

// Creates and returns Options which contains the values specified.
func NewOptions() *Options {
	options := &Options{}
	flag.BoolVar(&options.version, "version", false, "")
	flag.BoolVar(&options.version, "v", false, "")
	flag.BoolVar(&options.debug, "debug", false, "")
	flag.BoolVar(&options.debug, "d", false, "")
	flag.BoolVar(&options.links, "l", false, "")
	flag.BoolVar(&options.links, "links", false, "")
	flag.BoolVar(&options.images, "i", false, "")
	flag.BoolVar(&options.images, "images", false, "")
	flag.BoolVar(&options.xml, "xml", false, "")
	flag.BoolVar(&options.print, "print", false, "")
	flag.StringVar(&options.directory, "dir", "", "")
	flag.Parse()
	return options
}
