package payments

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/domain"
)

func (h *Handler) UpdatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	var input_payment input_payment

	defer r.Body.Close()

	id_str := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_str)
	if err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid payment ID", err.Error())
		return
	}

	userID, err := h.AuthUtils.GetUserIDFromRequest(r)
	if err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input_payment); err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	parsed_date, err := time.Parse("02-01-2006", input_payment.Date)
	if err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid date format, expected: DD-MM-YYYY", err.Error())
		return
	}

	domainPayment := &domain.Payment{
		ID:          uint(id),
		UserID:      userID,
		NetAmount:   input_payment.NetAmount,
		GrossAmount: input_payment.GrossAmount,
		Deductible:  input_payment.Deductible,
		Name:        input_payment.Name,
		Type:        input_payment.Type,
		Date:        parsed_date,
		Recurrent:   input_payment.Recurrent,
		Paid:        input_payment.Paid,
	}

	updated_payment, err := h.PaymentService.Update(domainPayment)
	if err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusInternalServerError, "error during update", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(updated_payment)
}
