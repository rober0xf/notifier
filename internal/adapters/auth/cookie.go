package auth

import "net/http"

const (
	TokenExpirationHours = 6
	BearerPrefix         = "BEARER "
	SessionCookieName    = "session_token"
	AuthHeaderName       = "Authorization"
)

func Set_auth_cookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true, // prevent XSS
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(TokenExpirationHours),
	})
}
