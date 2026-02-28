package postgres

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rober0xf/notifier/internal/domain/entity"
	database "github.com/rober0xf/notifier/internal/infraestructure/persistance/postgres/sqlc_generated"
)

func databaseToDomainUser(dbUser *database.User) *entity.User {
	var hash string
	if dbUser.EmailVerificationHash.Valid {
		hash = dbUser.EmailVerificationHash.String
	}

	var expiresAt time.Time
	if dbUser.TokenExpiresAt.Valid {
		expiresAt = dbUser.TokenExpiresAt.Time
	}

	return &entity.User{
		ID:                    int(dbUser.ID),
		Username:              dbUser.Username,
		Email:                 dbUser.Email,
		Password:              dbUser.Password,
		Active:                dbUser.Active,
		EmailVerificationHash: hash,
		CreatedAt:             dbUser.CreatedAt.Time,
		TokenExpiresAt:        expiresAt,
	}
}

func databaseToDomainPayment(dbPayment *database.Payment) *entity.Payment {
	payment := &entity.Payment{
		ID:        dbPayment.ID,
		UserID:    int(dbPayment.UserID),
		Name:      dbPayment.Name,
		Amount:    numericToFloat(dbPayment.Amount),
		Type:      entity.TransactionType(dbPayment.Type),
		Category:  entity.CategoryType(dbPayment.Category),
		Date:      dbPayment.Date,
		Paid:      dbPayment.Paid,
		Recurrent: dbPayment.Recurrent,
	}

	if dbPayment.Frequency.Valid {
		freq := entity.FrequencyType(dbPayment.Frequency.FrequencyType)
		payment.Frequency = &freq
	}
	if dbPayment.DueDate.Valid {
		payment.DueDate = &dbPayment.DueDate.String
	}
	if dbPayment.PaidAt.Valid {
		payment.PaidAt = &dbPayment.PaidAt.String
	}
	if dbPayment.ReceiptUrl.Valid {
		payment.ReceiptURL = &dbPayment.ReceiptUrl.String
	}

	return payment
}

func numericToFloat(num pgtype.Numeric) float64 {
	if !num.Valid {
		return 0.0
	}

	f64, _ := num.Float64Value()
	return f64.Float64
}

func floatToNumeric(num float64) pgtype.Numeric {
	var numeric pgtype.Numeric
	_ = numeric.Scan(fmt.Sprintf("%.2f", num))

	return numeric
}

func setNullableFieldsForUpdate(payment *entity.Payment) database.UpdatePaymentParams {
	amount := floatToNumeric(payment.Amount)
	params := &database.UpdatePaymentParams{
		ID:         int32(payment.ID),
		Name:       payment.Name,
		Amount:     amount,
		Type:       mapTransactionTypeToDB(payment.Type),
		Category:   mapCategoryTypeToDB(payment.Category),
		Date:       payment.Date,
		DueDate:    toNullableText(payment.DueDate),
		Paid:       payment.Paid,
		PaidAt:     toNullableText(payment.PaidAt),
		Recurrent:  payment.Recurrent,
		ReceiptUrl: toNullableText(payment.ReceiptURL),
	}

	if payment.Frequency != nil {
		params.Frequency = database.NullFrequencyType{
			FrequencyType: mapFrequencyTypeToDB(*payment.Frequency),
			Valid:         true,
		}
	}

	return *params
}

func toNullableText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}

	return pgtype.Text{String: *s, Valid: true}
}

func mapTransactionTypeToDB(t entity.TransactionType) database.TransactionType {
	return database.TransactionType(string(t))
}

func mapCategoryTypeToDB(t entity.CategoryType) database.CategoryType {
	return database.CategoryType(string(t))
}

func mapFrequencyTypeToDB(t entity.FrequencyType) database.FrequencyType {
	return database.FrequencyType(string(t))
}
