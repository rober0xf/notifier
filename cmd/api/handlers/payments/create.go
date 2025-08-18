package payments

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/domain"
)

func (h *Handler) CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	var input input_payment

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // check all fields are filled

	if err := decoder.Decode(&input); err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid request", "")
		return
	}

	parsed_date, err := time.Parse("02-01-2006", input.Date)
	if err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid data form, expected: YY-MM-DD", "")
		return
	}

	// validate the user input
	if err := validate.Struct(input); err != nil {
		log.Printf("validation error: %v", err)
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid input", err.Error())
		return
	}

	userID, err := h.AuthUtils.GetUserIDFromRequest(r)
	if err != nil || userID == 0 {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid user id", err.Error())
		return
	}

	payment := &domain.Payment{
		UserID:      userID,
		NetAmount:   input.NetAmount,
		GrossAmount: input.GrossAmount,
		Deductible:  input.Deductible,
		Name:        input.Name,
		Type:        input.Type,
		Date:        parsed_date,
		Recurrent:   input.Recurrent,
		Paid:        input.Paid,
	}

	// use the business logic
	if err := h.PaymentService.Create(payment); err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusInternalServerError, "could not create payment", "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}
