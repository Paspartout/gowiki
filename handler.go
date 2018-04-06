package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var validTitle = regexp.MustCompile(`^([a-zA-Z0-9]+)$`)
var validPath = regexp.MustCompile(`^/(view|edit|save|delete)/([a-zA-Z0-9]+)$`)
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

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := strings.Replace(r.FormValue("body"), "\r", "", -1)
	newTitle := r.FormValue("title")

	// save/overwrite page
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// rename/move page if title was changed
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
		if m == nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, m[2])
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
