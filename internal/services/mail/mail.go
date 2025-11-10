package mail

import (
	_ "embed"
)

//go:embed disposable_emails.txt
var disposableEmails []byte

// emails that we do not want
func MustDisposableEmail() (disp_emails []string) {
	// file, err := os.Open("internal/services/mail/disposable_emails.txt")
	// if err != nil {
	// 	log.Panic(err)
	// }
	// defer file.Close()
	//
	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	disp_emails = append(disp_emails, scanner.Text())
	// }
	for _, email := range disposableEmails {
		disp_emails = append(disp_emails, string(email))
	}
	return disp_emails
}
