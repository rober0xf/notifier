package handlers

import (
	"encoding/json"
	"errors"
	"goapi/config"
	"goapi/dbconnect"
	"goapi/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type fieldsInput struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type jwtclaims struct {
	Identity string `json:"identity"`
	jwt.RegisteredClaims
}

func generateToken(email string) (string, error) {
	expiration := time.Now().Add(6 * time.Hour)
	claims := &jwtclaims{
		Identity: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var db = dbconnect.DB
	var user models.User
	var credentials fieldsInput

	// check if the fields are valid
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// search for the user in the database
	if err := db.Where("email = ?", credentials.Identity).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "invalid credentials", http.StatusNotFound)
		} else {
			http.Error(w, "error while connecting", http.StatusInternalServerError)
		}
		return
	}

	// validate the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	tokenResponse, err := generateToken(user.Email)
	if err != nil {
		http.Error(w, "error while generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenResponse})
}
