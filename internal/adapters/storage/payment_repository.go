package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/rober0xf/notifier/internal/ports"
	database "github.com/rober0xf/notifier/internal/ports/db"
)

var _ ports.PaymentRepository = (*Repository)(nil)

func (r *Repository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	var amount_numeric pgtype.Numeric
	if err := amount_numeric.Scan(fmt.Sprintf("%.2f", payment.Amount)); err != nil {
		return err
	}

	params := database.CreatePaymentParams{
		UserID:     int32(payment.UserID),
		Name:       payment.Name,
		Amount:     amount_numeric,
		Type:       database.TransactionType(payment.Type),
		Category:   database.CategoryType(payment.Category),
		Date:       payment.Date,
		DueDate:    to_nullable_text(payment.DueDate),
		PaidAt:     to_nullable_text(payment.PaidAt),
		ReceiptUrl: to_nullable_text(payment.ReceiptURL),
	}
	if payment.Frequency != nil {
		params.Frequency = database.NullFrequencyType{
			FrequencyType: database.FrequencyType(*payment.Frequency),
			Valid:         true,
		}
	}

	created_payment, err := r.queries.CreatePayment(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError // for unique
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return dto.ErrAlreadyExists
			}
		}
		return fmt.Errorf("error creating payment: %w", err)
	}

	payment.ID = created_payment.ID
	return nil
}

func (r *Repository) GetAllPayments(ctx context.Context) ([]domain.Payment, error) {
	db_payments, err := r.queries.GetAllPayments(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all payments: %w", err)
	}
	if len(db_payments) == 0 {
		return []domain.Payment{}, nil
	}
	payments := make([]domain.Payment, 0, len(db_payments))
	for _, p := range db_payments {
		payments = append(payments, *database_to_domain_payment(&p))
	}

	return payments, nil
}

func (r *Repository) GetPaymentByID(ctx context.Context, id int) (*domain.Payment, error) {
	payment, err := r.queries.GetPaymentByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dto.ErrNotFound
		}
		return nil, fmt.Errorf("error getting payment by id: %w", err)
	}

	return database_to_domain_payment(&payment), nil
}

func (r *Repository) GetAllPaymentsFromUser(ctx context.Context, email string) ([]domain.Payment, error) {
	// first we need to check if the email exists in our db
	_, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dto.ErrNotFound
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	// if the user doesnt have any payments it returns []
	db_payments, err := r.queries.GetAllPaymentsFromUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error getting all payments from user: %w", err)
	}
	if len(db_payments) == 0 {
		return []domain.Payment{}, nil
	}

	payments := make([]domain.Payment, 0, len(db_payments))
	for _, p := range db_payments {
		payments = append(payments, *database_to_domain_payment(&p))
	}

	return payments, nil
}

func (r *Repository) UpdatePayment(ctx context.Context, payment *domain.Payment) error {
	params := set_nullable_fields_for_update(payment)
	rows_affected, err := r.queries.UpdatePayment(ctx, params)
	if err != nil {
		return fmt.Errorf("error updating payment: %w", err)
	}
	if rows_affected == 0 {
		return dto.ErrNotFound
	}

	return nil
}

func (r *Repository) DeletePayment(ctx context.Context, id int) error {
	rows, err := r.queries.DeletePayment(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("error deleting payment: %w", err)
	}
	if rows == 0 {
		return dto.ErrNotFound
	}

	return nil
}

func database_to_domain_payment(db_payment *database.Payment) *domain.Payment {
	payment := &domain.Payment{
		ID:        db_payment.ID,
		UserID:    int(db_payment.UserID),
		Name:      db_payment.Name,
		Amount:    numeric_to_float(db_payment.Amount),
		Type:      domain.TransactionType(db_payment.Type),
		Category:  domain.CategoryType(db_payment.Category),
		Date:      db_payment.Date,
		Paid:      db_payment.Paid,
		Recurrent: db_payment.Recurrent,
	}
	if db_payment.Frequency.Valid {
		freq := domain.FrequencyType(db_payment.Frequency.FrequencyType)
		payment.Frequency = &freq
	}
	if db_payment.DueDate.Valid {
		payment.DueDate = &db_payment.DueDate.String
	}
	if db_payment.PaidAt.Valid {
		payment.PaidAt = &db_payment.PaidAt.String
	}
	if db_payment.ReceiptUrl.Valid {
		payment.ReceiptURL = &db_payment.ReceiptUrl.String
	}

	return payment
}

func numeric_to_float(num pgtype.Numeric) float64 {
	if !num.Valid {
		return 0.0
	}
	f64, _ := num.Float64Value()
	return f64.Float64
}

func float_to_numeric(num float64) pgtype.Numeric {
	var numeric pgtype.Numeric
	_ = numeric.Scan(fmt.Sprintf("%.2f", num))

	return numeric
}

func set_nullable_fields_for_update(payment *domain.Payment) database.UpdatePaymentParams {
	amount := float_to_numeric(payment.Amount)
	params := &database.UpdatePaymentParams{
		ID:         int32(payment.ID),
		Name:       payment.Name,
		Amount:     amount,
		Type:       database.TransactionType(payment.Type),
		Category:   database.CategoryType(payment.Category),
		Date:       payment.Date,
		DueDate:    to_nullable_text(payment.DueDate),
		Paid:       payment.Paid,
		PaidAt:     to_nullable_text(payment.PaidAt),
		Recurrent:  payment.Recurrent,
		ReceiptUrl: to_nullable_text(payment.ReceiptURL),
	}
	if payment.Frequency != nil {
		params.Frequency = database.NullFrequencyType{
			FrequencyType: database.FrequencyType(*payment.Frequency),
			Valid:         true,
		}
	}

	return *params
}

func to_nullable_text(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}

	return pgtype.Text{String: *s, Valid: true}
}
