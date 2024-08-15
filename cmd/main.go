package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"

	"github.com/siddhant-vij/Counter-Templ-HTMX/templates"
)

func main() {
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	pageComponent := templates.Page(0, 0)
	http.Handle("/", templ.Handler(pageComponent))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
