package main

import (
	"html/template"
	"net/http"

	"github.com/paspartout/gowiki/tmpl"
)

var templateMap map[string]*template.Template

func initTemplates(useLocal bool) error {
	const (
		templatePath   = "/tmpl/"
		templateBase   = "/tmpl/layout/base.html"
		templateEnding = ".html"
	)
	templates := []string{"view", "edit", "delete", "new", "pages"}

	templateMap = make(map[string]*template.Template)
	for _, tpl := range templates {
		newTmpl := template.New(tpl + templateEnding)
		// First load base template
		_, err := newTmpl.Parse(tmpl.FSMustString(useLocal, templateBase))
		if err != nil {
			return err
		}
		// Then add the specific one over it
		_, err = newTmpl.Parse(tmpl.FSMustString(useLocal, templatePath+tpl+templateEnding))
		if err != nil {
			return err
		}

		templateMap[tpl] = newTmpl
	}

	return nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templateMap[tmpl].ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
