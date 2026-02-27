package email

type MockSender struct {
	SentEmails []SentEmail
	Err        error
}

type SentEmail struct {
	To      []string
	Subject string
	Body    string
}

func NewMockSender() *MockSender {
	return &MockSender{
		SentEmails: []SentEmail{},
	}
}

func (m *MockSender) Send(to []string, subject, htmlBody string) error {
	if m.Err != nil {
		return m.Err
	}

	m.SentEmails = append(m.SentEmails, SentEmail{
		To:      to,
		Subject: subject,
		Body:    htmlBody,
	})

	return nil
}
