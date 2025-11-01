package storage

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/rober0xf/notifier/internal/ports"
	database "github.com/rober0xf/notifier/internal/ports/db"
)

// postgresql.org/docs/current/errcodes-appendix.html
// go-database-sql.org/errors.html

var _ ports.UserRepository = (*Repository)(nil)

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	created_user, err := r.queries.CreateUser(ctx, database.CreateUserParams{
		Username:              user.Username,
		Email:                 user.Email,
		Password:              user.Password,
		EmailVerificationHash: pgtype.Text{String: user.EmailVerificationHash},
		Timeout:               pgtype.Interval{Days: int32(user.Timeout)},
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "unique constraint") {
			return dto.ErrUserAlreadyExists
		}
		return dto.ErrRepository
	}

	user.ID = int(created_user.ID)
	return nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dto.ErrUserNotFound
		}
		return nil, dto.ErrRepository
	}

	return database_to_domain_user(&user), nil
}

func (r *Repository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	db_users, err := r.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, dto.ErrRepository
	}
	if len(db_users) == 0 {
		return []domain.User{}, nil
	}
	users := make([]domain.User, 0, len(db_users))
	for _, u := range db_users {
		// best append than users[i] for out of bound
		users = append(users, *database_to_domain_user(&u))
	}

	return users, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := r.queries.GetUserByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dto.ErrNotFound
		}
		return nil, dto.ErrRepository
	}

	return database_to_domain_user(&user), nil
}

func (r *Repository) UpdateUserProfile(ctx context.Context, id int, username, email string) error {
	rows, err := r.queries.UpdateUserProfile(ctx, database.UpdateUserProfileParams{
		ID:       int32(id),
		Username: username,
		Email:    email,
	})
	if err != nil {
		return dto.ErrRepository
	}
	if rows == 0 {
		return dto.ErrUserNotFound
	}

	return nil
}

func (r *Repository) UpdateUserPassword(ctx context.Context, id int, password string) error {
	rows, err := r.queries.UpdateUserPassword(ctx, database.UpdateUserPasswordParams{
		ID:       int32(id),
		Password: password,
	})
	if err != nil {
		return dto.ErrRepository
	}
	if rows == 0 {
		return dto.ErrUserNotFound
	}

	return nil
}

func (r *Repository) UpdateUserActive(ctx context.Context, id int, active bool) error {
	rows, err := r.queries.UpdateUserActive(ctx, database.UpdateUserActiveParams{
		ID:     int32(id),
		Active: active,
	})
	if err != nil {
		return dto.ErrRepository
	}
	if rows == 0 {
		return dto.ErrUserNotFound
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int) error {
	rows, err := r.queries.DeleteUser(ctx, int32(id))
	if err != nil {
		return dto.ErrRepository
	}
	if rows == 0 {
		return dto.ErrNotFound
	}

	return nil
}

func database_to_domain_user(db_user *database.User) *domain.User {
	var hash string
	if db_user.EmailVerificationHash.Valid {
		hash = db_user.EmailVerificationHash.String
	}
	var timeout time.Duration
	if db_user.Timeout.Valid {
		timeout = time.Duration(db_user.Timeout.Microseconds) * time.Microsecond
	}

	return &domain.User{
		ID:                    int(db_user.ID),
		Username:              db_user.Username,
		Email:                 db_user.Email,
		Password:              db_user.Password,
		Active:                db_user.Active,
		EmailVerificationHash: hash,
		CreatedAt:             db_user.CreatedAt.Time,
		Timeout:               timeout,
	}
}
