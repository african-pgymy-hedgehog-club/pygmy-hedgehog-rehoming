package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
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
const templateFolder = appRoot + "/template"
const layout = templateFolder + "/_layout.html"

// Render pased template file
func renderTemplate(w http.ResponseWriter, tmpl string) {
	tmpl = filepath.Join(templateFolder, "block", tmpl+".html")
	t, err := template.ParseFiles(tmpl, layout)
	if err != nil { // If there was an error parsing the templates send an error back to the client
		httpError := err.Error()
		if APP_ENV != "dev" {
			httpError = "Sorry, there was an error"
		}

		http.Error(w, httpError, http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		renderTemplate(w, "index")
	} else {
		template := r.URL.Path[1:]
		renderTemplate(w, template)
	}
}

func mooHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Moo22"))
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/moo", mooHandler)

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
