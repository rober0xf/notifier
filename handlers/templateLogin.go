package handlers

import (
	"bytes"
	"goapi/templates"
	"io"
	"log"
	"net/http"
)

func LoginTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := templates.LoginTemplate().Render(r.Context(), w); err != nil {
				http.Error(w, "error rendering template", http.StatusInternalServerError)
				return
			}
			return
		}

		if r.Method == http.MethodPost {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("error reading request body: %v", err)
				http.Error(w, "error reading request body", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			formData := LoginHandler(w, r)
			if formData != nil {
				log.Printf("message: %s", formData.Message)
				log.Printf("status: %d", formData.Status)
				http.Error(w, formData.Message, formData.Status)
				return
			}

			w.Header().Set("HX-Redirect", "/")
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
