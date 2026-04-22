package postgres

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rober0xf/notifier/internal/domain/entity"
	database "github.com/rober0xf/notifier/internal/infraestructure/persistance/postgres/sqlc_generated"
)

func generateUsername(email string) string {
	base, _, _ := strings.Cut(email, "@")
	n := rand.Intn(1000)
	return fmt.Sprintf("%s_%d", base, n)
}

// TODO: database user_token to entity user_token
// TODO: return error
func databaseToDomainUser(dbUser *database.User) *entity.User {
	return &entity.User{
		ID:           int(dbUser.ID),
		Username:     dbUser.Username,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash.String,
		IsActive:     dbUser.IsActive,
		Role:         entity.Role(dbUser.Role),
		CreatedAt:    dbUser.CreatedAt.Time,
	}
}

func databaseToDomainPayment(dbPayment *database.Payment) (entity.Payment, error) {
	amount, err := numericToFloat(dbPayment.Amount)
	if err != nil {
		return entity.Payment{}, fmt.Errorf("mapping payment %d to domain: %w", dbPayment.ID, err)
	}

	payment := &entity.Payment{
		ID:        dbPayment.ID,
		UserID:    int(dbPayment.UserID),
		Name:      dbPayment.Name,
		Amount:    amount,
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

	return *payment, nil
}

func numericToFloat(num pgtype.Numeric) (float64, error) {
	if !num.Valid {
		return 0.0, fmt.Errorf("numeric value is invalid")
	}

	f64, err := num.Float64Value()
	if err != nil {
		return 0.0, fmt.Errorf("casting numeric float64: %w", err)
	}

	return f64.Float64, nil
}

func floatToNumeric(amount float64) (pgtype.Numeric, error) {
	var num pgtype.Numeric
	if err := num.Scan(fmt.Sprintf("%.2f", amount)); err != nil {
		return num, fmt.Errorf("invalid amount value %.2f: %w", amount, err)
	}

	return num, nil
}

func mapPaymentToUpdateParams(payment *entity.Payment) (database.UpdatePaymentParams, error) {
	amount, err := floatToNumeric(payment.Amount)
	if err != nil {
		return database.UpdatePaymentParams{}, fmt.Errorf("invalid amount: %w", err)
	}

	params := database.UpdatePaymentParams{
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

	return params, nil
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
