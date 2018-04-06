package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

var dataPath = "data/"

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

// Removes a page
func (p *Page) remove() error {
	filename := dataPath + p.Title + extension
	return os.Remove(filename)
}

// Renames the page to the new title
func (p *Page) rename(newTitle string) error {
	if !validTitle.MatchString(newTitle) {
		return fmt.Errorf("new title \"%s\" is invalid", newTitle)
	}

	filename := dataPath + p.Title + extension
	newFileanme := dataPath + newTitle + extension
	err := os.Rename(filename, newFileanme)
	if err != nil {
		return err
	}

	p.Title = newTitle
	return nil
}

// Loads a page using its title
func loadPage(title string) (*Page, error) {
	filename := dataPath + title + extension
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
