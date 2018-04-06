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

// Config stores the configuration for the wiki that is parsed
// from the command line
type Config struct {
	Address  string // Adress to bind to
	DataPath string // Path to md files
}

func parseConfig() Config {
	address := flag.String("address",
		":8080", "The address to listen to")
	dataPath := flag.String("path",
		"data/", "Path to the folder that contains the document files")
	flag.Parse()
	return Config{Address: *address, DataPath: *dataPath}
}

func main() {
	config := parseConfig()
	log.Fatal(listen(config))
}
