package main

import (
	"flag"
	"log"
)

const (
	extension      = ".md"
	staticPath     = "static/"
	frontPageTitle = "FrontPage"
)

// Generate go file that embeds the static files
//go:generate esc -o static/static.go -ignore '[A-Za-z]*\.go' -pkg static static

// Generate go file that embeds the template files
//go:generate esc -o tmpl/tmpl.go -ignore '[A-Za-z]*\.go' -pkg tmpl tmpl

// Config stores the configuration for the wiki that is parsed
// from the command line
type Config struct {
	Address  string // Adress to bind to
	DataPath string // Path to md files
	UseLocal bool   // True if user wants to use local static files e.g. for development
}

func parseConfig() Config {
	address := flag.String("address",
		":8080", "The address to listen to")
	dataPath := flag.String("path",
		"data/", "Path to the folder that contains the document files")
	useLocal := flag.Bool("local", false,
		"Use local static files and templates instead of embedded ones.")
	flag.Parse()
	return Config{Address: *address, DataPath: *dataPath, UseLocal: *useLocal}
}

func main() {
	config := parseConfig()
	log.Fatal(listen(config))
}
