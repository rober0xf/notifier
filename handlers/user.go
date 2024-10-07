package handlers

import (
	"encoding/json"
	"errors"
	"goapi/dbconnect"
	"goapi/models"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func hashPassw(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(pass), err
}

func CreateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	/// when we use new to decode the json input we are using the complete model including the pass, so we have to encrypt it
	user := new(models.User)

	defer r.Body.Close() // close the body to prevent resource leak
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// check if empty fields
	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "empty fields", http.StatusBadRequest)
		return
	}

	// first we hash the password and then we store it
	password, err := hashPassw(user.Password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	user.Password = password

	if err := db.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			http.Error(w, "user already exists", http.StatusConflict)
		} else if errors.Is(err, gorm.ErrInvalidData) {
			http.Error(w, "invalid data", http.StatusBadRequest)
		} else {
			http.Error(w, "error while connecting", http.StatusInternalServerError)
		}
		return
	}

	// clear the password before returning
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	db := dbconnect.DB
	id := mux.Vars(r)["id"] // take the id from the request

	if id == "" {
		getAllUsers(db, w)
	} else {
		getUserFromId(db, id, w)
	}
}

func getAllUsers(db *gorm.DB, w http.ResponseWriter) {
	var users []models.User

	// find retuns all found values
	if err := db.Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "resource not found", http.StatusNotFound)
		} else {
			http.Error(w, "error while connecting", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func getUserFromId(db *gorm.DB, id string, w http.ResponseWriter) {
	var user models.User

	// first returns the first value
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, "error while connecting", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, "error while connecting", http.StatusInternalServerError)
		}
		return
	}

	// we create a custom structure to use it as an intermediary
	type updateUser struct {
		Name     string `gorm:"not null" json:"name"`
		Username string `json:"username"`
		Email    string `gorm:"not null" json:"email"`
		Password string `gorm:"not null" json:"password,omitempty"`
	}

	var inputUser updateUser

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&inputUser); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if inputUser.Name == "" || inputUser.Email == "" || inputUser.Password == "" {
		http.Error(w, "empty fields", http.StatusBadRequest)
		return
	}

	// update the password
	password, err := hashPassw(inputUser.Password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	user.Name = inputUser.Name
	user.Username = inputUser.Username
	user.Email = inputUser.Email
	user.Password = password

	if err := db.Save(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			http.Error(w, "invalid data", http.StatusBadRequest)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User

	if id == "" {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		http.Error(w, "error deleting", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "user deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
