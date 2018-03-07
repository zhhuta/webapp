package main

import (
	"html/template"
	"net/http"
	"time"
	//Google appengine import
	"google.golang.org/appengine"
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

type templateParams struct {
	Date string
	Time string
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
	curentDate := time.Now().Local()
	params.Date = curentDate.Format("2006-02-01")
	params.Time = curentDate.Format("3:04 PM")

	if r.Method == "GET" {
		indexTemplate.Execute(w, params)
	}
}
