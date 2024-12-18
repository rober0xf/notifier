package handlers

import (
	"encoding/json"
	"errors"
	"goapi/config"
	"goapi/dbconnect"
	"goapi/models"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	ID       string `json:"ID"`
	jwt.RegisteredClaims
}

type errorResponse struct {
	Message string
	Status  int
}

func getUserId(w http.ResponseWriter, r *http.Request) (string, int) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "no header authorization", http.StatusUnauthorized)
		return "no header authorization", -1
	}

	// remove the prefix bearer
	if len(tokenString) > 7 && strings.ToUpper(tokenString[:7]) == "BEARER " {
		tokenString = tokenString[7:]
	} else {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
		return "invalid authorization header", -1
	}

	userIDstr := getIDfromToken(tokenString)
	if userIDstr == "" {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return "invalid token", -1
	}

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusInternalServerError)
		return "invalid user id", -1
	}

	return "", userID
}

func getIDfromToken(tokenString string) string {
	secret := []byte(config.JwtKey)

	token, err := jwt.ParseWithClaims(tokenString, &jwtclaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return ""
	}

	if claims, ok := token.Claims.(*jwtclaims); ok {
		return claims.ID
	}

	return ""
}

func generateToken(email string, userID string) (string, error) {
	expiration := time.Now().Add(6 * time.Hour)
	claims := &jwtclaims{
		Identity: email,
		ID:       userID,
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

func LoginHandler(w http.ResponseWriter, r *http.Request) *errorResponse {
	var db = dbconnect.DB
	var user models.User
	var credentials fieldsInput

	// checking if the request comes from a form or json
	defer r.Body.Close()
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			return &errorResponse{"invalid json request", http.StatusBadRequest}
		}
	} else {
		// check if the fields are valid
		if err := r.ParseForm(); err != nil {
			return &errorResponse{"invalid request", http.StatusBadRequest}
		}
		credentials.Identity = r.FormValue("identity")
		credentials.Password = r.FormValue("password")
	}

	log.Printf("buscando user con email: %v", credentials.Identity)

	// search for the user in the database
	if err := db.Where("email = ?", credentials.Identity).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("llegaste hasta aca?")
			return &errorResponse{"user not found", http.StatusNotFound}
		}
		log.Printf("error al buscar el usuario: %v", err)
		return &errorResponse{"error while connecting", http.StatusInternalServerError}
	}

	// validate the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return &errorResponse{"invalid credentials", http.StatusUnauthorized}
	}

	// take from the token the user.ID
	var userID string = strconv.FormatUint(uint64(user.ID), 10)
	tokenResponse, err := generateToken(user.Email, userID)
	if err != nil {
		return &errorResponse{"error while generating token", http.StatusInternalServerError}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session:id",
		Value:    tokenResponse,
		Path:     "/",
		HttpOnly: true, // prevent XSS
	})

	// json response
	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": tokenResponse})
	} else {
		// frontend response
		log.Printf("todo salio bien y vas a ser redireccionado")
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
	return nil
}
