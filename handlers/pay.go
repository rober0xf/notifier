package handlers

import (
	"encoding/json"
	"errors"
	"goapi/dbconnect"
	"goapi/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type inputPayment struct {
	NetAmount   float64 `json:"net_amount" validate:"required"`
	GrossAmount float64 `json:"gross_amount"`
	Deductible  float64 `json:"deductible"`
	Name        string  `gorm:"not null" json:"name" validate:"required"`
	Type        string  `gorm:"not null" json:"type" validate:"required"`
	Date        string  `json:"date" validate:"required"`
	Recurrent   bool    `gorm:"not null" json:"recurrent" validate:"required"`
	Paid        bool    `gorm:"not null" json:"paid" validate:"required"`
}

func CreatePayment(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	inputPayment := new(inputPayment)
	payment := new(models.Payment)

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // check all fields are filled

	if err := decoder.Decode(&inputPayment); err != nil {
		log.Printf("el error es: %v: \n", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	parsedDate, err := time.Parse("02-01-2006", inputPayment.Date)
	if err != nil {
		log.Printf("error al parsear la fecha: %v", err)
		http.Error(w, "invalid data form, expected: YY-MM-DD", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(inputPayment); err != nil {
		log.Printf("validation error: %v", err)
		http.Error(w, "invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	message, userID := getUserId(w, r)
	if userID == -1 {
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	payment.UserID = uint(userID)
	payment.NetAmount = inputPayment.NetAmount
	payment.GrossAmount = inputPayment.GrossAmount
	payment.Deductible = inputPayment.Deductible
	payment.Name = inputPayment.Name
	payment.Type = inputPayment.Type
	payment.Date = parsedDate
	payment.Recurrent = inputPayment.Recurrent
	payment.Paid = inputPayment.Paid

	if err := db.Create(&payment).Error; err != nil {
		http.Error(w, "error during payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

func GetPayment(w http.ResponseWriter, r *http.Request) {
	db := dbconnect.DB
	id := mux.Vars(r)["id"]
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "no header authorization", http.StatusUnauthorized)
		return
	}
	if len(tokenString) > 7 && strings.ToUpper(tokenString[:7]) == "BEARER " {
		tokenString = tokenString[7:]
	} else {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(getIDfromToken(tokenString))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusInternalServerError)
		return
	}

	if id != "" {
		getPaymentFromId(db, id, userID, w, r)
	}

	getAllPayments(db, userID, w)

}

func getAllPayments(db *gorm.DB, userID int, w http.ResponseWriter) {
	payments := []models.Payment{}

	if err := db.Where("user_id = ?", userID).Find(&payments).Error; err != nil {
		http.Error(w, "error getting payments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payments); err != nil {
		http.Error(w, "error encoding payments", http.StatusInternalServerError)
		return
	}
}

func getPaymentFromId(db *gorm.DB, id string, userID int, w http.ResponseWriter, r *http.Request) {
	var payment models.Payment

	message, userID := getUserId(w, r)
	if userID == -1 {
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	if err := db.First(&payment, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "payment not found", http.StatusNotFound)
			return
		}
		http.Error(w, "could not find payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

func UpdatePayment(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updatedPayment inputPayment
	var payment models.Payment

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&updatedPayment); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	parsedDate, err := time.Parse("02-01-2006", updatedPayment.Date)
	if err != nil {
		http.Error(w, "invalid date format, expected DD-MM-YYYY", http.StatusBadRequest)
		return
	}
	if err := db.First(&payment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "payment not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	payment.NetAmount = updatedPayment.NetAmount
	payment.GrossAmount = updatedPayment.GrossAmount
	payment.Deductible = updatedPayment.Deductible
	payment.Name = updatedPayment.Name
	payment.Type = updatedPayment.Type
	payment.Date = parsedDate
	payment.Recurrent = updatedPayment.Recurrent
	payment.Paid = updatedPayment.Paid

	if err := db.Save(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "payment updated successfully"}
	json.NewEncoder(w).Encode(response)
}

func DeletePayment(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var payment models.Payment

	message, userID := getUserId(w, r)
	if userID == -1 {
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	if err := db.First(&payment, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "payment not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := db.Delete(&payment).Error; err != nil {
		http.Error(w, "error deleting", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "payment deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
