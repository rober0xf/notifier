package config

import (
	"encoding/json"
	"log"
	"net/http"
)

func TestMail(w http.ResponseWriter, r *http.Request) {
	err := MailSender.SendMail([]string{"reciever@gmail.com"}, "sending test", "this is the body")
	if err != nil {
		http.Error(w, "error sending the mail", http.StatusInternalServerError)
		log.Printf("el error es: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "mail sent successfully"}
	json.NewEncoder(w).Encode(response)
}
