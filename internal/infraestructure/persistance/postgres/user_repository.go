package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		EmailVerificationHash: pgtype.Text{String: user.EmailVerificationHash, Valid: true},
		TokenExpiresAt:        pgtype.Timestamptz{Time: user.TokenExpiresAt, Valid: true},
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

func (r *UserRepository) CreateUserToken(ctx context.Context, token *entity.UserToken) (*entity.UserToken, error) {
	row, err := r.queries.CreateUserToken(ctx, database.CreateUserTokenParams{
		UserID:    int32(token.UserID),
		TokenHash: token.TokenHash,
		Purpose:   database.TokenPurpose(token.Purpose),
		ExpiresAt: pgtype.Timestamptz{Time: token.ExpiresAt, Valid: true},
	})

	if err != nil {
		log.Printf("DB error in CreateUserToken: %v", err)
		return nil, repoErr.ErrRepository
	}

	return &entity.UserToken{
		ID:        int(row.ID),
		UserID:    int(row.UserID),
		TokenHash: row.TokenHash,
		Purpose:   token.Purpose,
		Used:      row.Used,
		ExpiresAt: row.ExpiresAt.Time,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *UserRepository) VerifyAndConsumeToken(ctx context.Context, tokenHash string, purpose entity.TokenPurpose) (*entity.UserToken, error) {
	row, err := r.queries.VerifyAndConsumeToken(ctx, database.VerifyAndConsumeTokenParams{
		TokenHash: tokenHash,
		Purpose:   database.TokenPurpose(purpose),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoErr.ErrNotFound
		}
		return nil, fmt.Errorf("verify token query failed: %w", err)
	}

	return &entity.UserToken{
		ID:        int(row.ID),
		UserID:    int(row.UserID),
		TokenHash: row.TokenHash,
		Purpose:   entity.TokenPurpose(row.Purpose),
		Used:      row.Used,
		ExpiresAt: row.ExpiresAt.Time,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *UserRepository) GetTokenByHash(ctx context.Context, tokenHash string, purpose entity.TokenPurpose) (*entity.UserToken, error) {
	row, err := r.queries.GetTokenByHash(ctx, database.GetTokenByHashParams{
		TokenHash: tokenHash,
		Purpose:   database.TokenPurpose(purpose),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoErr.ErrNotFound
		}
		return nil, repoErr.ErrRepository
	}

	return &entity.UserToken{
		ID:        int(row.ID),
		UserID:    int(row.UserID),
		TokenHash: row.TokenHash,
		Purpose:   entity.TokenPurpose(row.Purpose),
		Used:      row.Used,
		ExpiresAt: row.ExpiresAt.Time,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *UserRepository) DeleteByUserAndPurpose(ctx context.Context, userID int, purpose entity.TokenPurpose) error {
	return r.queries.DeleteByUserAndPurpose(ctx, database.DeleteByUserAndPurposeParams{
		UserID:  int32(userID),
		Purpose: database.TokenPurpose(purpose),
	})
}

func (r *UserRepository) DeleteOldTokens(ctx context.Context) (int64, error) {
	rows, err := r.queries.DeleteOldTokens(ctx)

	if err != nil {
		return 0, repoErr.ErrRepository
	}

	return rows, nil
}
