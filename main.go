package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
}

var APP_ENV = os.Getenv("APP_ENV")

func init() {
	if APP_ENV == "" {
		APP_ENV = "dev"
	}
}

const appRoot = "/go/src/app/"
const templateFolder = appRoot + "template/"
const nav = templateFolder + "/block/nav.html"
const layout = templateFolder + "_layout.html"

// Write error response to the client
func clientError(w http.ResponseWriter, err error) {
	httpError := err.Error()
	if APP_ENV != "dev" {
		httpError = "Sorry, there was an error"

		log.Println(err.Error()) // Log error if it's not displayed to the client
	}

	http.Error(w, httpError, http.StatusInternalServerError)
}

// Render pased template file
func renderTemplate(w http.ResponseWriter, tmpl string) {
	tmpl = filepath.Join(templateFolder, "block", tmpl+".html")
	t, err := template.ParseFiles(layout, nav, tmpl)
	if err != nil {
		clientError(w, err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var p = r.URL.Path

	if p == "/" {
		renderTemplate(w, "index")
	} else if strings.HasPrefix(p, "/api") {
		apiHandler(w, r)
	} else if strings.HasPrefix(p, "/adoption") { // Handle any fake adoption pages
		var template = "adoption"
		renderTemplate(w, template)
	} else {
		template := r.URL.Path[1:]
		renderTemplate(w, template)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)

	if APP_ENV == "dev" { // Serve static content if app environment is dev (in production nginx will serve)
		// Create static content handlers
		var staticPaths = map[string]string{
			"/js/":     "js/",
			"/css/":    "css/",
			"/images/": "images/",
			"/fonts/":  "fonts/",
		}
		for pattern, path := range staticPaths {
			fs := http.FileServer(http.Dir(path))
			http.Handle(pattern, http.StripPrefix(pattern, fs))
		}
	}

	http.ListenAndServe(":8080", nil)
}
