package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/Paspartout/gowiki/static"
	"github.com/microcosm-cc/bluemonday"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var validTitle = regexp.MustCompile(`^([a-zA-Z0-9]+)$`)
var validPath = regexp.MustCompile(`^/(((view|delete)/([a-zA-Z0-9]+))|((edit|save)/([a-zA-Z0-9]*)))$`)
var linkRegex = regexp.MustCompile(`\[([a-zA-Z0-9]+)\]`)

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	// Markdown rendering
	bodyRendered := blackfriday.Run(p.Body)
	// Sanitize html
	bodyRendered = bluemonday.UGCPolicy().SanitizeBytes(bodyRendered)

	// Interlinking
	bodyRendered = linkRegex.ReplaceAllFunc(bodyRendered,
		func(link []byte) []byte {
			linkTitle := string(link)
			linkTitle = linkTitle[1 : len(linkTitle)-1]
			linkStr := "<a href=\"" + linkTitle + "\">" + linkTitle + "</a>"
			return []byte(linkStr)
		})

	renderedPage := &RenderedPage{
		Title: p.Title,
		Body:  template.HTML(bodyRendered)}

	renderTemplate(w, "view", renderedPage)
}

// Handles editing pages or creating a new page
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil && os.IsNotExist(err) {
		renderTemplate(w, "new", title)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "edit", p)
}

// Handles saving and moving pages
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := strings.Replace(r.FormValue("body"), "\r", "", -1)
	newTitle := r.FormValue("title")
	if title == "" {
		title = newTitle // use form title for creating a new page
	}

	// Check for valid title before saving
	if !validTitle.MatchString(title) {
		http.Error(w, "Title name is invalid: "+title, http.StatusBadRequest)
		return
	}

	// Create or Overwrite page
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Rename/Move page if title was changed
	if newTitle != title {
		err := p.rename(newTitle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		title = newTitle
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, title string) {
	deletionConfirmed := r.FormValue("Confirmed") == "True"
	p := Page{Title: title}

	if deletionConfirmed {
		err := p.remove()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/view/"+frontPageTitle, http.StatusFound)
	} else {
		renderTemplate(w, "delete", p)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		// log.Printf("%#v\n", m)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		// m[4]+m[7] is the content of the capture groups that eventually contain
		// the page title in /edit/title and /save/title but always contain
		// the page title in /view/title and /delete/title
		fn(w, r, m[4]+m[7])
	}
}

func pagesHandler(w http.ResponseWriter, r *http.Request) {
	dataFiles, err := ioutil.ReadDir(dataPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter for page files
	pages := make([]string, 0, len(dataFiles))
	for _, f := range dataFiles {
		fName := f.Name()
		if !f.IsDir() && fName[len(fName)-3:] == extension {
			pages = append(pages, fName[:len(fName)-3])
		}
	}

	renderTemplate(w, "pages", pages)
}

func listen(conf Config) error {
	// TODO: Refactor model
	dataPath = conf.DataPath

	err := initTemplates(conf.UseLocal)
	if err != nil {
		log.Fatal("error initializing templates:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/view/"+frontPageTitle, http.StatusFound)
	})

	// Operations on pages
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/delete/", makeHandler(deleteHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))

	// View list of all pages
	http.HandleFunc("/pages", pagesHandler)
	http.Handle("/static/",
		http.FileServer(static.FS(conf.UseLocal)))

	return http.ListenAndServe(conf.Address, nil)
}
