package email

import "context"

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

func (m *MockSender) Send(ctx context.Context, to []string, subject, htmlBody string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

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
