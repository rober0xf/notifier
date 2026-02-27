package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	database "github.com/rober0xf/notifier/internal/infraestructure/persistance/postgres/sqlc_generated"
)

type PaymentRepository struct {
	db      *pgxpool.Pool
	queries *database.Queries
}

var _ repository.PaymentRepository = (*PaymentRepository)(nil)

func NewPaymentRepository(db *pgxpool.Pool) repository.PaymentRepository {
	return &PaymentRepository{
		db:      db,
		queries: database.New(db),
	}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	var amountNumeric pgtype.Numeric
	if err := amountNumeric.Scan(fmt.Sprintf("%.2f", payment.Amount)); err != nil {
		return err
	}

	params := database.CreatePaymentParams{
		UserID:     int32(payment.UserID),
		Name:       payment.Name,
		Amount:     amountNumeric,
		Type:       mapTransactionTypeToDB(payment.Type),
		Category:   mapCategoryTypeToDB(payment.Category),
		Date:       payment.Date,
		DueDate:    toNullableText(payment.DueDate),
		PaidAt:     toNullableText(payment.PaidAt),
		ReceiptUrl: toNullableText(payment.ReceiptURL),
	}

	if payment.Frequency != nil {
		params.Frequency = database.NullFrequencyType{
			FrequencyType: mapFrequencyTypeToDB(*payment.Frequency),
			Valid:         true,
		}
	}

	createdPayment, err := r.queries.CreatePayment(ctx, params)

	if err != nil {
		var pgErr *pgconn.PgError // for unique
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return repoErr.ErrAlreadyExists
			}
		}
		return fmt.Errorf("error creating payment: %w", err)
	}

	payment.ID = createdPayment.ID
	return nil
}

func (r *PaymentRepository) GetAllPayments(ctx context.Context) ([]entity.Payment, error) {
	dbPayments, err := r.queries.GetAllPayments(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all payments: %w", err)
	}
	if len(dbPayments) == 0 {
		return []entity.Payment{}, nil
	}

	payments := make([]entity.Payment, 0, len(dbPayments))
	for _, p := range dbPayments {
		payments = append(payments, *databaseToDomainPayment(&p))
	}

	return payments, nil
}

func (r *PaymentRepository) GetPaymentByID(ctx context.Context, id int) (*entity.Payment, error) {
	payment, err := r.queries.GetPaymentByID(ctx, int32(id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoErr.ErrNotFound
		}
		return nil, fmt.Errorf("error getting payment by id: %w", err)
	}

	return databaseToDomainPayment(&payment), nil
}

func (r *PaymentRepository) GetAllPaymentsFromUser(ctx context.Context, userID int) ([]entity.Payment, error) {
	// first we need to check if the email exists in our db
	_, err := r.queries.GetUserByID(ctx, int32(userID))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoErr.ErrNotFound
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	// if the user doesnt have any payments it returns []
	dbPayments, err := r.queries.GetAllPaymentsFromUser(ctx, int32(userID))

	if err != nil {
		return nil, fmt.Errorf("error getting all payments from user: %w", err)
	}

	if len(dbPayments) == 0 {
		return []entity.Payment{}, nil
	}

	payments := make([]entity.Payment, 0, len(dbPayments))
	for _, p := range dbPayments {
		payments = append(payments, *databaseToDomainPayment(&p))
	}

	return payments, nil
}

func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment *entity.Payment) error {
	params := setNullableFieldsForUpdate(payment)
	rowsAffected, err := r.queries.UpdatePayment(ctx, params)

	if err != nil {
		return fmt.Errorf("error updating payment: %w", err)
	}

	if rowsAffected == 0 {
		return repoErr.ErrNotFound
	}

	return nil
}

func (r *PaymentRepository) DeletePayment(ctx context.Context, id int) error {
	rows, err := r.queries.DeletePayment(ctx, int32(id))

	if err != nil {
		return fmt.Errorf("error deleting payment: %w", err)
	}

	if rows == 0 {
		return repoErr.ErrNotFound
	}

	return nil
}
