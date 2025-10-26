package mail

import (
	"bufio"
	"log"
	"os"
)

// emails that we do not want
func MustDisposableEmail() (disp_emails []string) {
	file, err := os.Open("internal/services/mail/disposable_emails.txt")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		disp_emails = append(disp_emails, scanner.Text())
	}
	return disp_emails
}
