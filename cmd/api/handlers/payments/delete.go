package payments

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
)

func (h *Handler) DeletePaymentHandler(w http.ResponseWriter, r *http.Request) {
	id_str := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_str)
	if err != nil {
		httphelpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid payment ID", err.Error())
		return
	}

	err = h.PaymentService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, dto.ErrPaymentNotFound) {
			httphelpers.WriteErrorResponse(w, http.StatusNotFound, "payment not found", "")
		} else {
			httphelpers.WriteErrorResponse(w, http.StatusInternalServerError, "could not delete payment", "")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	response := map[string]string{"message": "payment deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
