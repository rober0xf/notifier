package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rober0xf/notifier/internal/core/shared"
	"github.com/rober0xf/notifier/internal/core/utils"
	"github.com/rober0xf/notifier/internal/models"
)

type AuthHandler struct {
	authService shared.AuthServiceInterface
	authUtils   utils.AuthUtils
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	credentials, err := h.authUtils.Parse_login_request(r)
	if err != nil {
		utils.Write_error_response(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	user, err := h.authService.ExistsUser(ctx, credentials)
	if err != nil {
		log.Printf("%s authentication failed: %v", credentials.Email, err)
		utils.Write_error_response(w, http.StatusUnauthorized, "authentication failed", "")
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Printf("error generation token: %v", err)
		utils.Write_error_response(w, http.StatusInternalServerError, "error while token generation", "")
		return
	}

	utils.Set_auth_cookie(w, token)

	// if it comes from json
	if utils.Is_json_request(r) {
		utils.Write_json_response(w, http.StatusOK, shared.LoginResponse{
			Token: token,
			User: shared.UserInfo{
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

func (h *AuthHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // close the body to prevent resource leak

	// struct used for decode the input
	var input_user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input_user); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// check if empty fields
	if input_user.Email == "" || input_user.Password == "" {
		http.Error(w, "empty fields", http.StatusBadRequest)
		return
	}

	// here we use the service logic
	err := h.authService.CreateUserService(input_user.Email, input_user.Password)
	if err != nil {
		switch {
		case errors.Is(err, shared.ErrUserAlreadyExists):
			http.Error(w, "user already exists", http.StatusConflict)
		case errors.Is(err, shared.ErrInvalidUserData):
			http.Error(w, "invalid data", http.StatusBadRequest)
		case errors.Is(err, shared.ErrPasswordHashing):
			http.Error(w, "error hashing password", http.StatusInternalServerError)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input_user)
}

func (h *AuthHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.authService.GetAllUsersService()
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

func (h *AuthHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request, id string) {
	id_int, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	user, err := h.authService.GetUserFromIDService(uint(id_int))
	if err != nil {
		if errors.Is(err, shared.ErrUserNotFound) {
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

func (h *AuthHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input_user struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input_user); err != nil {
		utils.Write_error_response(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	userID, err := h.authUtils.Get_userID_from_request(r)
	if err != nil {
		utils.Write_error_response(w, http.StatusUnauthorized, "unauthorized", err.Error())
		return
	}

	// create the user with the new data
	user := &models.User{
		ID:       userID,
		Name:     input_user.Name,
		Username: input_user.Username,
		Email:    input_user.Email,
	}

	updated_user, err := h.authService.UpdateUserService(user)
	if err != nil {
		utils.Write_error_response(w, http.StatusInternalServerError, "error during update", err.Error())
		return
	}

	// clean to show in response
	updated_user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated_user)
}

// TODO: delete_user services and refactor the handler
func (h *AuthHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = mux.Vars(r)["id"]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "user deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
