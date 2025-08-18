package users

import (
	"encoding/json"
	"net/http"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/domain"
)

func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input_user struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input_user); err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	userID, err := h.AuthUtils.GetUserIDFromRequest(r)
	if err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", err.Error())
		return
	}

	// create the user with the new data
	user := &domain.User{
		ID:    userID,
		Name:  input_user.Name,
		Email: input_user.Email,
	}

	updated_user, err := h.UserService.Update(user)
	if err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusInternalServerError, "error during update", err.Error())
		return
	}

	// clean to show in response
	updated_user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated_user)
}
