package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rober0xf/notifier/internal/core/shared"
	"github.com/rober0xf/notifier/internal/core/utils"
	"github.com/rober0xf/notifier/internal/models"
)

type PaymentsHandler struct {
	paymentService shared.PaymentServiceInterface
	authUtils      utils.AuthUtils
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type input_payment struct {
	NetAmount   float64 `json:"net_amount" validate:"required"`
	GrossAmount float64 `json:"gross_amount"`
	Deductible  float64 `json:"deductible"`
	Name        string  `gorm:"not null" json:"name" validate:"required"`
	Type        string  `gorm:"not null" json:"type" validate:"required"`
	Date        string  `json:"date" validate:"required"`
	Recurrent   bool    `gorm:"not null" json:"recurrent" validate:"required"`
	Paid        bool    `gorm:"not null" json:"paid" validate:"required"`
}

func (h *PaymentsHandler) CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	var input input_payment

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // check all fields are filled

	if err := decoder.Decode(&input); err != nil {
		utils.Write_error_response(w, http.StatusBadRequest, "invalid request", "")
		return
	}

	parsed_date, err := time.Parse("02-01-2006", input.Date)
	if err != nil {
		utils.Write_error_response(w, http.StatusBadRequest, "invalid data form, expected: YY-MM-DD", "")
		return
	}

	// validate the user input
	if err := validate.Struct(input); err != nil {
		log.Printf("validation error: %v", err)
		utils.Write_error_response(w, http.StatusBadRequest, "invalid input", err.Error())
		return
	}

	userID, err := h.authUtils.Get_userID_from_request(r)
	if err != nil || userID == 0 {
		utils.Write_error_response(w, http.StatusBadRequest, "invalid user id", err.Error())
		return
	}

	payment := &models.Payment{
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
	if err := h.paymentService.CreatePaymentService(payment); err != nil {
		utils.Write_error_response(w, http.StatusInternalServerError, "could not create payment", "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

func (h *PaymentsHandler) GetAllPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := h.authUtils.Get_userID_from_request(r)
	if err != nil {
		switch {
		case errors.Is(err, shared.ErrMissingAuthHeader):
			utils.Write_error_response(w, http.StatusUnauthorized, "missing Authorization header", "")
		case errors.Is(err, shared.ErrInvalidHeaderFormat):
			utils.Write_error_response(w, http.StatusBadRequest, "invalid Authorization header format", "")
		case errors.Is(err, shared.ErrInvalidToken):
			utils.Write_error_response(w, http.StatusUnauthorized, "invalid or expired token", "")
		default:
			log.Printf("unexpected auth error: %v", err)
			utils.Write_error_response(w, http.StatusInternalServerError, "internal server error", "")
		}
		return
	}

	payments, err := h.paymentService.GetAllPaymentsService(user_id)
	if err != nil {
		if errors.Is(err, shared.ErrPaymentNotFound) {
			utils.Write_error_response(w, http.StatusNotFound, "no payments found", "")
			return
		}
		log.Printf("unexpected payment service error: %v", err)
		utils.Write_error_response(w, http.StatusInternalServerError, "internal error", "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(payments); err != nil {
		log.Printf("error encoding payments: %v", err)
		utils.Write_error_response(w, http.StatusInternalServerError, "internal server error", "")
		return
	}
}

func (h *PaymentsHandler) GetPaymentByIDHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the url
	id_str := mux.Vars(r)["id"]
	if id_str == "" {
		utils.Write_error_response(w, http.StatusBadRequest, "must provide an id", "")
		return
	}
	id, err := strconv.Atoi(id_str)
	if err != nil {
		utils.Write_error_response(w, http.StatusBadRequest, "invalid id value", "")
		return
	}
	if id <= 0 {
		utils.Write_error_response(w, http.StatusBadRequest, "id must be positive", "")
		return
	}

	// get the user_id
	user_id, err := h.authUtils.Get_userID_from_request(r)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "missing authorization header"):
			utils.Write_error_response(w, http.StatusUnauthorized, "missing authorization header", "")
		case strings.Contains(err.Error(), "invalid authorization header format"):
			utils.Write_error_response(w, http.StatusBadRequest, "invalid authorization header format", "")
		case strings.Contains(err.Error(), "invalid token"):
			utils.Write_error_response(w, http.StatusUnauthorized, "invalid or expired token", "")
		default:
			log.Printf("unexpected auth error: %v", err)
			utils.Write_error_response(w, http.StatusInternalServerError, "internal error", "")
		}
		return
	}

	payment, err := h.paymentService.GetPaymentFromIDService(uint(id), user_id)
	if err != nil {
		if errors.Is(err, shared.ErrPaymentNotFound) {
			utils.Write_error_response(w, http.StatusNotFound, "payment not found", "")
			return
		}
		log.Printf("Unexpected payment service error: %v", err)
		utils.Write_error_response(w, http.StatusInternalServerError, "internal error", "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

// func (s *Store) UpdatePayment(w http.ResponseWriter, r *http.Request) {
// 	id := mux.Vars(r)["id"]
// 	var updatedPayment inputPayment
// 	var payment models.Payment

// 	defer r.Body.Close()
// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()

// 	if err := decoder.Decode(&updatedPayment); err != nil {
// 		http.Error(w, "invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	parsedDate, err := time.Parse("02-01-2006", updatedPayment.Date)
// 	if err != nil {
// 		http.Error(w, "invalid date format, expected DD-MM-YYYY", http.StatusBadRequest)
// 		return
// 	}
// 	if err := s.DB.First(&payment, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			http.Error(w, "payment not found", http.StatusNotFound)
// 			return
// 		}
// 		http.Error(w, "internal error", http.StatusInternalServerError)
// 		return
// 	}

// 	payment.NetAmount = updatedPayment.NetAmount
// 	payment.GrossAmount = updatedPayment.GrossAmount
// 	payment.Deductible = updatedPayment.Deductible
// 	payment.Name = updatedPayment.Name
// 	payment.Type = updatedPayment.Type
// 	payment.Date = parsedDate
// 	payment.Recurrent = updatedPayment.Recurrent
// 	payment.Paid = updatedPayment.Paid

// 	if err := s.DB.Save(&payment).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			http.Error(w, "invalid data", http.StatusBadRequest)
// 			return
// 		} else {
// 			http.Error(w, "internal error", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	response := map[string]string{"message": "payment updated successfully"}
// 	json.NewEncoder(w).Encode(response)
// }

// func (s *Store) DeletePayment(w http.ResponseWriter, r *http.Request) {
// 	id := mux.Vars(r)["id"]
// 	var payment models.Payment

// 	message, userID := getUserId(w, r)
// 	if userID == -1 {
// 		http.Error(w, message, http.StatusBadRequest)
// 		return
// 	}

// 	if err := s.DB.First(&payment, "id = ? AND user_id = ?", id, userID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			http.Error(w, "payment not found", http.StatusNotFound)
// 			return
// 		}
// 		http.Error(w, "internal error", http.StatusInternalServerError)
// 		return
// 	}

// 	if err := s.DB.Delete(&payment).Error; err != nil {
// 		http.Error(w, "error deleting", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	response := map[string]string{"message": "payment deleted successfully"}
// 	json.NewEncoder(w).Encode(response)
// }
