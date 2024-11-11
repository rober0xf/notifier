package handlers

import (
	"encoding/json"
	"errors"
	"goapi/models"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func hashPassw(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(pass), err
}

func (s *Store) CreateUser(w http.ResponseWriter, r *http.Request) {
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

	// fix autoincrement id
	var existingUser models.User
	if err := s.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "credentials already exists", http.StatusConflict)
		return
	}

	// first we hash the password and then we store it
	password, err := hashPassw(user.Password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	user.Password = password

	if err := s.DB.Create(&user).Error; err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				http.Error(w, "user already exists", http.StatusConflict)
				return
			default:
				http.Error(w, "error while creating user", http.StatusInternalServerError)
				return
			}
		} else if errors.Is(err, gorm.ErrInvalidData) {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "error while connecting", http.StatusInternalServerError)
			return
		}
	}

	// clear the password before returning
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (s *Store) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"] // take the id from the request

	if id == "" {
		s.getAllUsers(w)
	} else {
		s.getUserFromId(id, w)
	}
}

func (s *Store) getAllUsers(w http.ResponseWriter) {
	var users []models.User

	// find retuns all found values
	if err := s.DB.Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "resource not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (s *Store) getUserFromId(id string, w http.ResponseWriter) {
	var user models.User

	// first returns the first value
	if err := s.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *Store) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User

	if err := s.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
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

	if err := s.DB.Save(&user).Error; err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				http.Error(w, "user already exists", http.StatusConflict)
				return
			default:
				http.Error(w, "error while updating", http.StatusInternalServerError)
				return
			}
		} else if errors.Is(err, gorm.ErrInvalidData) {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *Store) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User

	if id == "" {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := s.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := s.DB.Delete(&user).Error; err != nil {
		http.Error(w, "error deleting", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "user deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
