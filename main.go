package main

import (
	"log"
	"net/http"
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

func listen(conf Config) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/view/"+frontPageTitle, http.StatusFound)
	})

	// Operations on pages
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/delete/", makeHandler(deleteHandler))

	// View list of all pages
	http.HandleFunc("/pages", pagesHandler)
	http.Handle("/static/", http.StripPrefix("/static",
		http.FileServer(http.Dir(staticPath))))

	return http.ListenAndServe(conf.Address, nil)
}

func main() {
	config := Config{
		Address:  addr,
		DataPath: dataPath,
	}
	log.Fatal(listen(config))
}
