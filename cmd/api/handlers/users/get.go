package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *Handler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// clean the password for the return
	for i := range users {
		users[i].Password = ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request, id string) {
	id_int, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.GetUserFromID(uint(id_int))
	if err != nil {
		if errors.Is(err, dto.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
