package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
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
	pc, fn, line, _ := runtime.Caller(1)
	httpError := fmt.Sprintf("%s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)

	if APP_ENV != "dev" {
		log.Println(httpError) // Log error if it's not displayed to the client

		httpError = "Sorry, there was an error"
	}

	http.Error(w, httpError, http.StatusInternalServerError)
}

// Render passed template file
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
	} else if p == "/adoption" {
		http.Redirect(w, r, "https://africanpygmyhedgehogclub.co.uk/club-rescue", 301)
	} else {
		template := r.URL.Path[1:]
		renderTemplate(w, template)
	}
}

func main() {
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		fh, err := os.Open("images/favicon.ico")
		defer fh.Close()
		if err != nil {
			clientError(w, err)
			return
		}

		io.Copy(w, fh) // Send file to client
	})
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
