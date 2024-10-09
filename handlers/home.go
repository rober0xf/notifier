package handlers

import (
	"html/template"
	"net/http"
)

type homeData struct {
	Title    string
	Header   string
	User     string
	IsLogged bool
	Message  string
	Items    []string
}

func HomeHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := homeData{
			Title:    "goapi",
			Header:   "Home Page",
			User:     "Unknown",
			IsLogged: false,
			Message:  "testing front",
			Items:    []string{"item 1", "item 2", "item 3"},
		}

		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, "error rendering", http.StatusInternalServerError)
		}
	}
}
