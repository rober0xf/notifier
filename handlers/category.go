package handlers

import (
	"encoding/json"
	"errors"
	"goapi/dbconnect"
	"goapi/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type inputCategory struct {
	Name      string `json:"name"`
	Priority  uint   `json:"priority"`
	Recurrent bool   `json:"recurrent"`
	Notify    bool   `json:"notify"`
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

func CreateCategory(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	inputCate := new(inputCategory)
	category := new(models.Category)

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&inputCate); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if inputCate.Name == "" {
		http.Error(w, "a name is required", http.StatusBadRequest)
		return
	}

	message, userID := getUserId(w, r)
	if userID == -1 {
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	category.UserID = uint(userID)
	category.Name = inputCate.Name
	category.Priority = inputCate.Priority
	category.Recurrent = inputCate.Recurrent
	category.Notify = inputCate.Notify

	if err := db.Create(&category).Error; err != nil {
		http.Error(w, "error creating category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	db := dbconnect.DB
	message, userID := getUserId(w, r)

	if userID == -1 {
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	if id != "" {
		getUserFromId(db, id, w)
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

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "error encoding categories", http.StatusInternalServerError)
		return
	}
}

func getCategoryFromId(db *gorm.DB, id string, w http.ResponseWriter) {
	var category models.Category

	if err := db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "category not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func UpdateCategory(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updatedCategory inputCategory
	var category models.Category

	/*
		message, userID := getUserId(w, r)
		if userID == -1 {
			http.Error(w, message, http.StatusBadRequest)
			return
		}
	*/

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "category not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if updatedCategory.Name == "" {
		http.Error(w, "empty fields", http.StatusBadRequest)
		return
	}

	category.Name = updatedCategory.Name
	category.Priority = updatedCategory.Priority
	category.Recurrent = updatedCategory.Recurrent
	category.Notify = updatedCategory.Notify

	if err := db.Save(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func DeleteCategory(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var category models.Category

	if err := db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "category not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := db.Delete(&category).Error; err != nil {
		http.Error(w, "error deleting", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "category deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
