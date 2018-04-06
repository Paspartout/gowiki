package main

import (
	"log"
)

const (
	extension      = ".md"
	dataPath       = "data/"   // TODO: Make configurable through command line flag
	staticPath     = "static/" // TODO: Pack into executable
	addr           = ":8080"
	frontPageTitle = "FrontPage"
)

// Config stores the configuration for the wiki that is parsed
// from the command line
type Config struct {
	Address  string // Adress to bind to
	DataPath string // Path to md files
}

func main() {
	config := Config{
		Address:  addr,
		DataPath: dataPath,
	}
	log.Fatal(listen(config))
}
