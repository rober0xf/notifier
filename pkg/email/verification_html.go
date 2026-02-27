package email

import (
	"fmt"
)

func VerificationEmailHTML(email string, token string, baseURL string) string {
	verificationURL := fmt.Sprintf("%s/v1/users/email_verification/%s/%s", baseURL, email, token)

	return fmt.Sprintf(`
		<html>
			<body>
				<h1>Email verification</h1>
				<p>click the link to verify:</p>
				<a href="%s">verify account</a>
			</body>
		</html>
	`, verificationURL)
}
