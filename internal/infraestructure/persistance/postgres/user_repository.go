package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/domain/repository"

	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	database "github.com/rober0xf/notifier/internal/infraestructure/persistance/postgres/sqlc_generated"
)

// postgresql.org/docs/current/errcodes-appendix.html
// go-database-sql.org/errors.html

type UserRepository struct {
	db      *pgxpool.Pool
	queries *database.Queries
}

var _ repository.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &UserRepository{
		db:      db,
		queries: database.New(db),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	createdUser, err := r.queries.CreateUser(ctx, database.CreateUserParams{
		Username:              user.Username,
		Email:                 user.Email,
		Password:              user.Password,
		EmailVerificationHash: pgtype.Text{String: user.EmailVerificationHash},
		TokenExpiresAt:        pgtype.Timestamptz{Time: (user.TokenExpiresAt)},
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "unique constraint") {
			return repoErr.ErrAlreadyExists
		}
		return repoErr.ErrRepository
	}

	user.ID = int(createdUser.ID)
	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoErr.ErrNotFound
		}
		return nil, repoErr.ErrRepository
	}

	return databaseToDomainUser(&user), nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	dbUsers, err := r.queries.GetAllUsers(ctx)

	if err != nil {
		return nil, repoErr.ErrRepository
	}

	if len(dbUsers) == 0 {
		return []entity.User{}, nil
	}

	users := make([]entity.User, 0, len(dbUsers))
	for _, u := range dbUsers {
		// best append than users[i] for out of bound
		users = append(users, *databaseToDomainUser(&u))
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	user, err := r.queries.GetUserByID(ctx, int32(id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoErr.ErrNotFound
		}
		return nil, repoErr.ErrRepository
	}

	return databaseToDomainUser(&user), nil
}

func (r *UserRepository) UpdateUserProfile(ctx context.Context, id int, username, email string) error {
	rows, err := r.queries.UpdateUserProfile(ctx, database.UpdateUserProfileParams{
		ID:       int32(id),
		Username: username,
		Email:    email,
	})

	if err != nil {
		return repoErr.ErrRepository
	}

	if rows == 0 {
		return repoErr.ErrNotFound
	}

	return nil
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, id int, password string) error {
	rows, err := r.queries.UpdateUserPassword(ctx, database.UpdateUserPasswordParams{
		ID:       int32(id),
		Password: password,
	})

	if err != nil {
		return repoErr.ErrRepository
	}

	if rows == 0 {
		return repoErr.ErrNotFound
	}

	return nil
}

func (r *UserRepository) UpdateUserActive(ctx context.Context, id int, active bool) error {
	rows, err := r.queries.UpdateUserActive(ctx, database.UpdateUserActiveParams{
		ID:     int32(id),
		Active: active,
	})

	if err != nil {
		return repoErr.ErrRepository
	}

	if rows == 0 {
		return repoErr.ErrNotFound
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	rows, err := r.queries.DeleteUser(ctx, int32(id))

	if err != nil {
		return repoErr.ErrRepository
	}

	if rows == 0 {
		return repoErr.ErrNotFound
	}

	return nil
}
