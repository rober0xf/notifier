package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id_str := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_str)
	if err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid user id", "")
		return
	}

	err = h.UserService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, dto.ErrUserNotFound) {
			httphelpers.WriteErrorResponse(w, http.StatusNotFound, "user not found", "")
		} else {
			httphelpers.WriteErrorResponse(w, http.StatusInternalServerError, "could not delete usr", "")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "user deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
