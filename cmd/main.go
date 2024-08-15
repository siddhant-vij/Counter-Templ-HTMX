package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/siddhant-vij/Counter-Templ-HTMX/templates"
)

var globalCount int
var sessionManager *scs.SessionManager

func main() {
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sessionCount := sessionManager.GetInt(r.Context(), "sessionCount")
		pageComponent := templates.Page(globalCount, sessionCount)
		pageComponent.Render(r.Context(), w)
	})

	mux.HandleFunc("/global", func(w http.ResponseWriter, r *http.Request) {
		globalCount++
		gc := templates.GlobalCount(globalCount)
		gc.Render(r.Context(), w)
	})

	mux.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		curCount := sessionManager.GetInt(r.Context(), "sessionCount")
		curCount++
		sessionManager.Put(r.Context(), "sessionCount", curCount)
		sc := templates.SessionCount(curCount)
		sc.Render(r.Context(), w)
	})

	muxWithSessionMiddleware := sessionManager.LoadAndSave(mux)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", muxWithSessionMiddleware))
}
