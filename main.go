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
const header = templateFolder + "/header.html"
const footer = templateFolder + "/footer.html"

// Render pased template file
func renderTemplate(w http.ResponseWriter, tmpl string) {
	tmpl = filepath.Join(templateFolder, tmpl+".html")
	t, err := template.ParseFiles(tmpl, footer, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.NotFound(w, r)
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
			"/css":     "css/",
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
