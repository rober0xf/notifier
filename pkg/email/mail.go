package email

import (
	_ "embed"
)

//go:embed disposable_emails.txt
var disposableEmails []byte

// emails that we do not want
func MustDisposableEmail() (disp_emails []string) {
	for _, email := range disposableEmails {
		disp_emails = append(disp_emails, string(email))
	}

	return disp_emails
}
