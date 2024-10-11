package handlers

import (
	"goapi/templates"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	name := "Notifier"

	err := templates.HomeTemplate(name).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "error rendering templates", http.StatusInternalServerError)
		return
	}
}
