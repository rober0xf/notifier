package auth

import (
	"context"

	"google.golang.org/api/idtoken"
)

type googleVerifier struct {
	clientID string
}

type GoogleVerifier interface {
	Verify(idToken string) (*GoogleUser, error)
}

func NewGoogleVerifier(clientID string) GoogleVerifier {
	return &googleVerifier{
		clientID: clientID,
	}
}

func (g *googleVerifier) Verify(idToken string) (*GoogleUser, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, g.clientID)
	if err != nil {
		return nil, err
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return nil, ErrInvalidEmailToken
	}

	emailVerified, _ := payload.Claims["email_verified"].(bool)
	if !emailVerified {
		return nil, ErrEmailNotVerified
	}

	name, _ := payload.Claims["name"].(string)

	return &GoogleUser{
		Sub:   payload.Subject,
		Email: email,
		Name:  name,
	}, nil
}
