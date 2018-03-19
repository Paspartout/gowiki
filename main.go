package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

const (
	extension      = ".md"
	templatePath   = "tmpl/"
	dataPath       = "data/"
	staticPath     = "./static/"
	frontPageTitle = "FrontPage"
)

var templates = template.Must(template.ParseFiles(
	templatePath+"edit.html",
	templatePath+"view.html"))

var validPath = regexp.MustCompile(`^/(view|edit|save)/([a-zA-Z0-9]+)$`)
var linkRegex = regexp.MustCompile(`\[([a-zA-Z0-9]+)\]`)

// Page represents a page of the wiki
type Page struct {
	Title string
	Body  []byte
}

// RenderedPage represents a page that has been rendered to html
type RenderedPage struct {
	Title string
	Body  template.HTML
}

func (p *Page) save() error {
	filename := dataPath + p.Title + extension
	err := ioutil.WriteFile(filename, p.Body, 0600)
	_, isPerr := err.(*os.PathError)
	if err != nil && isPerr {
		// Try to fix path error by making dataPath directory
		err = os.Mkdir(dataPath, 0700)
		if err != nil {
			return err
		}
		log.Printf("Creating %s directory for pages", dataPath)
		return p.save()
	} else if err != nil {
		return err
	}
	return nil
}

func loadPage(title string) (*Page, error) {
	filename := dataPath + title + extension
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/view/"+frontPageTitle, http.StatusFound)
	})

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.Handle("/static/", http.StripPrefix("/static",
		http.FileServer(http.Dir(staticPath))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
