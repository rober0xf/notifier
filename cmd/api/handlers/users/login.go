package users

import (
	"log"
	"net/http"

	"github.com/rober0xf/notifier/internal/adapters/authentication"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	credentials, err := h.AuthUtils.ParseLoginRequest(r)
	if err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	user, err := h.AuthUtils.ExistsUser(ctx, credentials)
	if err != nil {
		log.Printf("%s authentication failed: %v", credentials.Email, err)
		httphelpers.WriteErrorResponse(w, http.StatusUnauthorized, "authentication failed", "")
		return
	}

	token, err := h.AuthUtils.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Printf("error generation token: %v", err)
		httphelpers.WriteErrorResponse(w, http.StatusInternalServerError, "error while token generation", "")
		return
	}

	authentication.SetAuthCookie(w, token)

	// if it comes from json
	if httphelpers.IsJSONRequest(r) {
		httphelpers.WriteJSONResponse(w, http.StatusOK, dto.LoginResponse{
			Token: token,
			User: dto.UserInfo{
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
