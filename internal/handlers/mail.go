package handlers

import (
	"encoding/json"
	"net/http"

	database "github.com/rober0xf/notifier/internal/db"
)

func TestMail(w http.ResponseWriter, r *http.Request) {
	err := database.MailSender.SendMail([]string{"test@gmail.com"}, "sending test", "this is the body")
	if err != nil {
		http.Error(w, "error sending the mail", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "mail sent successfully"}
	json.NewEncoder(w).Encode(response)
}
