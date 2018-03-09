package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"google.golang.org/appengine"
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

type templateParams struct {
	Date    string
	Time    string
	Notice  string
	Warning string
	Name    string
}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	params := templateParams{}
	currentDate := time.Now().Local()
	params.Date = currentDate.Format("2006-02-01")
	params.Time = currentDate.Format("3:04 PM")

	if r.Method == "GET" {
		indexTemplate.Execute(w, params)
	}

	if r.Method == "POST" {

		name := r.FormValue("name")
		params.Name = name
		if name == "" {
			name = "Anonymous"
		}

		message := r.FormValue("message")
		if r.FormValue("message") == "" {
			w.WriteHeader(http.StatusBadRequest)

			params.Warning = "No message"
			indexTemplate.Execute(w, params)
			return
		}

		params.Notice = fmt.Sprintf("Message from %s: %s", name, message)
		indexTemplate.Execute(w, params)
	}
}
