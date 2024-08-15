package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/siddhant-vij/Counter-Templ-HTMX/templates"
)

type GlobalState struct {
	Count int
}

var global GlobalState
var sessionManager *scs.SessionManager

func getHandler(w http.ResponseWriter, r *http.Request) {
	sessionCount := sessionManager.GetInt(r.Context(), "sessionCount")
	pageComponent := templates.Page(global.Count, sessionCount)
	pageComponent.Render(r.Context(), w)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Form.Has("global") {
		global.Count++
	}

	if r.Form.Has("session") {
		curCount := sessionManager.GetInt(r.Context(), "sessionCount")
		curCount++
		sessionManager.Put(r.Context(), "sessionCount", curCount)
	}

	getHandler(w, r)
}

func main() {
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postHandler(w, r)
			return
		}
		getHandler(w, r)
	})

	muxWithSessionMiddleware := sessionManager.LoadAndSave(mux)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", muxWithSessionMiddleware))
}
