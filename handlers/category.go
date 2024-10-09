package handlers

import (
	"encoding/json"
	"fmt"
	"goapi/dbconnect"
	"goapi/models"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	db := dbconnect.DB
	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		http.Error(w, "no header authorization", http.StatusUnauthorized)
		return
	}

	// remove the prefix bearer
	if len(tokenString) > 7 && strings.ToUpper(tokenString[:7]) == "BEARER " {
		tokenString = tokenString[7:]
	} else {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}

	userIDstr := getIDfromToken(tokenString)

	if userIDstr == "" {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	fmt.Println(userIDstr)

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusInternalServerError)
		return
	}

	getAllCategories(db, userID, w)
}

func getAllCategories(db *gorm.DB, userID int, w http.ResponseWriter) {
	categories := []models.Category{}

	if err := db.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		http.Error(w, "error getting categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "error encoding categories", http.StatusInternalServerError)
		return
	}
}
