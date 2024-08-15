package main

import (
	"log"
	"net/http"

	"github.com/siddhant-vij/Counter-Templ-HTMX/templates"
)

type GlobalState struct {
	Count int
}

var global GlobalState

func getHandler(w http.ResponseWriter, r *http.Request) {
	pageComponent := templates.Page(global.Count, 0)
	pageComponent.Render(r.Context(), w)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Form.Has("global") {
		global.Count++
	}

	getHandler(w, r)
}

func main() {
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postHandler(w, r)
			return
		}
		getHandler(w, r)
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
