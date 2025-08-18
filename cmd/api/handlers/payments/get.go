package payments

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *Handler) GetAllPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := h.AuthUtils.GetUserIDFromRequest(r)
	if err != nil {
		switch {
		case errors.Is(err, dto.ErrMissingAuthHeader):
			httphelpers.WriteErrorResponse(w, http.StatusUnauthorized, "missing Authorization header", "")
		case errors.Is(err, dto.ErrInvalidHeaderFormat):
			httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid Authorization header format", "")
		case errors.Is(err, dto.ErrInvalidToken):
			httphelpers.WriteErrorResponse(w, http.StatusUnauthorized, "invalid or expired token", "")
		default:
			log.Printf("unexpected auth error: %v", err)
			httphelpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal server error", "")
		}
		return
	}

	payments, err := h.PaymentService.GetAllPayments(user_id)
	if err != nil {
		if errors.Is(err, dto.ErrPaymentNotFound) {
			httphelpers.WriteErrorResponse(w, http.StatusNotFound, "no payments found", "")
			return
		}
		log.Printf("unexpected payment service error: %v", err)
		httphelpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal error", "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(payments); err != nil {
		log.Printf("error encoding payments: %v", err)
		httphelpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal server error", "")
		return
	}
}

func (h *Handler) GetPaymentByIDHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the url
	id_str := mux.Vars(r)["id"]
	if id_str == "" {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "must provide an id", "")
		return
	}
	id, err := strconv.Atoi(id_str)
	if err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid id value", "")
		return
	}
	if id <= 0 {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "id must be positive", "")
		return
	}

	// get the user_id
	user_id, err := h.AuthUtils.GetUserIDFromRequest(r)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "missing authorization header"):
			httphelpers.WriteErrorResponse(w, http.StatusUnauthorized, "missing authorization header", "")
		case strings.Contains(err.Error(), "invalid authorization header format"):
			httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid authorization header format", "")
		case strings.Contains(err.Error(), "invalid token"):
			httphelpers.WriteErrorResponse(w, http.StatusUnauthorized, "invalid or expired token", "")
		default:
			log.Printf("unexpected auth error: %v", err)
			httphelpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal error", "")
		}
		return
	}

	payment, err := h.PaymentService.GetPaymentFromID(uint(id), user_id)
	if err != nil {
		if errors.Is(err, dto.ErrPaymentNotFound) {
			httphelpers.WriteErrorResponse(w, http.StatusNotFound, "payment not found", "")
			return
		}
		log.Printf("Unexpected payment service error: %v", err)
		httphelpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal error", "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}
