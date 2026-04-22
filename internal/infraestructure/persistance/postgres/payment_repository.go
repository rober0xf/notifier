package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *entity.Payment) (*entity.Payment, error) {
	amountNumeric, err := floatToNumeric(payment.Amount)
	if err != nil {
		return nil, err
	}

	params := database.CreatePaymentParams{
		UserID:     int32(payment.UserID),
		Name:       payment.Name,
		Amount:     amountNumeric,
		Type:       mapTransactionTypeToDB(payment.Type),
		Category:   mapCategoryTypeToDB(payment.Category),
		Date:       payment.Date,
		DueDate:    toNullableText(payment.DueDate),
		Recurrent:  payment.Recurrent,
		Paid:       payment.Paid,
		PaidAt:     toNullableText(payment.PaidAt),
		ReceiptUrl: toNullableText(payment.ReceiptURL),
	}

	if payment.Frequency != nil {
		params.Frequency = database.NullFrequencyType{
			FrequencyType: mapFrequencyTypeToDB(*payment.Frequency),
			Valid:         true,
		}
	}

	created, err := r.queries.CreatePayment(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, repoErr.ErrAlreadyExists
			}
		}

		return nil, fmt.Errorf("creating payment query failed: %w", err)
	}

	out, err := databaseToDomainPayment(&created)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (r *PaymentRepository) GetAllPayments(ctx context.Context) ([]entity.Payment, error) {
	dbPayments, err := r.queries.GetAllPayments(ctx, database.GetAllPaymentsParams{
		Limit:  50,
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("get all payments query failed: %w", err)
	}

	payments := make([]entity.Payment, 0, len(dbPayments))
	for _, p := range dbPayments {
		payment, err := databaseToDomainPayment(&p)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) GetPaymentByID(ctx context.Context, id int) (*entity.Payment, error) {
	payment, err := r.queries.GetPaymentByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoErr.ErrNotFound
		}

		return nil, fmt.Errorf("get payment by id query failed: %w", err)
	}

	out, err := databaseToDomainPayment(&payment)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (r *PaymentRepository) GetMyPayments(ctx context.Context, userID int) ([]entity.Payment, error) {
	// if the user doesnt have any payments it returns []
	dbPayments, err := r.queries.GetMyPayments(ctx, int32(userID))
	if err != nil {
		return nil, fmt.Errorf("error getting my payments: %w", err)
	}

	payments := make([]entity.Payment, 0, len(dbPayments))
	for _, p := range dbPayments {
		payment, err := databaseToDomainPayment(&p)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment *entity.Payment) error {
	params, err := mapPaymentToUpdateParams(payment)
	if err != nil {
		return err
	}

	rowsAffected, err := r.queries.UpdatePayment(ctx, params)
	if err != nil {
		return fmt.Errorf("update payment query failed: %w", err)
	}

	if rowsAffected == 0 {
		return repoErr.ErrNotFound
	}

	return nil
}

func (r *PaymentRepository) DeletePayment(ctx context.Context, id int) error {
	rows, err := r.queries.DeletePayment(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("delete payment query failed: %w", err)
	}

	if rows == 0 {
		return repoErr.ErrNotFound
	}

	return nil
}
