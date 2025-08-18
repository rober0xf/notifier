package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

// for routes
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // close the body to prevent resource leak

	// struct used for decode the input
	var input_user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input_user); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// check if empty fields
	if input_user.Name == "" || input_user.Email == "" || input_user.Password == "" {
		http.Error(w, "empty fields", http.StatusBadRequest)
		return
	}

	// here we use the service logic
	err := h.UserService.Create(input_user.Name, input_user.Email, input_user.Password)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrUserAlreadyExists):
			http.Error(w, "user already exists", http.StatusConflict)
		case errors.Is(err, dto.ErrInvalidUserData):
			http.Error(w, "invalid data", http.StatusBadRequest)
		case errors.Is(err, dto.ErrPasswordHashing):
			http.Error(w, "error hashing password", http.StatusInternalServerError)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// we return a custum structure without the password
	json.NewEncoder(w).Encode(struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Name:  input_user.Name,
		Email: input_user.Email,
	})
}
