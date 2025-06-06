package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	database "github.com/rober0xf/notifier/internal/db"
	"github.com/rober0xf/notifier/internal/models"
	"github.com/rober0xf/notifier/internal/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	jwtKey []byte
}

func NewAuthService(db *gorm.DB, jwtKey []byte) *AuthService {
	return &AuthService{
		db:     db,
		jwtKey: jwtKey,
	}
}

func (as *AuthService) get_userID_from_request(r *http.Request) (uint, error) {
	tokenString := r.Header.Get(types.AuthHeaderName)
	if tokenString == "" {
		return 0, errors.New("missing authorization header")
	}

	// remove the bearer prefix
	if len(tokenString) > 7 && strings.ToUpper(tokenString[:7]) == types.BearerPrefix {
		tokenString = tokenString[7:]
	} else {
		return 0, errors.New("invalid authorization header format")
	}

	userID, err := as.validate_token(tokenString)
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	return userID, nil
}

func (as *AuthService) validate_token(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return as.jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*types.JWTClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token claims")
	}

	return claims.UserID, nil
}

// creates a new jwt token for the user
func (as *AuthService) generate_token(userID uint, email string) (string, error) {
	expiration := time.Now().Add(types.TokenExpirationHours)
	claims := &types.JWTClaims{
		Email:  email,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			Audience:  []string{"notifier"},
			Issuer:    "notifier-service",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(as.jwtKey)
}

// validate the user using context instead of using the db directly
func (as *AuthService) exists_user(ctx context.Context, credentials types.LoginRequest) (*models.User, error) {
	var user models.User

	// check if the user exists
	if err := as.db.WithContext(ctx).Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("internal error: %w", err.Error())
	}

	// compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func parse_login_request(r *http.Request) (types.LoginRequest, error) {
	var credentials types.LoginRequest
	defer r.Body.Close() // idk if this is necessary

	// if it comes from json
	if types.Is_json_request(r) {
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			return credentials, fmt.Errorf("invalid json: %w", err)
		}
	} else {
		// if it comes from a form
		if err := r.ParseForm(); err != nil {
			return credentials, fmt.Errorf("invalid form: %w", err)
		}
		// set the values from the form
		credentials.Email = r.FormValue("email")
		credentials.Password = r.FormValue("password")
	}

	if credentials.Email == "" || credentials.Password == "" {
		return credentials, errors.New("email and password are empty")
	}

	return credentials, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth_service := NewAuthService(database.DB, database.JwtKey) // get the info

	credentials, err := parse_login_request(r)
	if err != nil {
		types.Write_error_response(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	user, err := auth_service.exists_user(ctx, credentials)
	if err != nil {
		log.Printf("%s authentication failed: %v", credentials.Email, err)
		types.Write_error_response(w, http.StatusUnauthorized, "authentication failed", "")
		return
	}

	token, err := auth_service.generate_token(user.ID, user.Email)
	if err != nil {
		log.Printf("error generation token: %v", err)
		types.Write_error_response(w, http.StatusInternalServerError, "error while token generation", "")
		return
	}

	types.Set_auth_cookie(w, token)

	// if it comes from json
	if types.Is_json_request(r) {
		types.Write_json_response(w, http.StatusOK, types.LoginResponse{
			Token: token,
			User: types.UserInfo{
				ID:    user.ID,
				Email: user.Email,
			},
		})
	} else {
		// from frontend
		w.Header().Set("HX-Redirect", "/dashboard")
		w.WriteHeader(http.StatusOK)
	}
}
