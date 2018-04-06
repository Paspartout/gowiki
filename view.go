package main

import (
	"html/template"
	"net/http"
)

const (
	templatePath = "tmpl/"
	templateBase = "tmpl/layout/base.html"
)

var templateMap = map[string]*template.Template{
	"view": template.Must(
		template.ParseFiles(templateBase, templatePath+"view.html")),
	"edit": template.Must(
		template.ParseFiles(templateBase, templatePath+"edit.html")),
	"delete": template.Must(
		template.ParseFiles(templateBase, templatePath+"delete.html")),
	"pages": template.Must(
		template.ParseFiles(templateBase, templatePath+"pages.html")),
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templateMap[tmpl].ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
