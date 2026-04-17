package postgres

import (
	"context"
	"database/sql"
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

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	createdUser, err := r.queries.CreateUser(ctx, database.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: pgtype.Text{String: user.PasswordHash, Valid: user.PasswordHash != ""},
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return nil, repoErr.ErrEmailAlreadyExists
			case "users_username_key":
				return nil, repoErr.ErrUsernameAlreadyExists
			}
			return nil, repoErr.ErrAlreadyExists
		}

		return nil, fmt.Errorf("create user query failed: %w", err)
	}

	return &entity.User{
		ID:           int(createdUser.ID),
		Username:     createdUser.Username,
		Email:        createdUser.Email,
		PasswordHash: createdUser.PasswordHash.String,
		IsActive:     createdUser.IsActive,
	}, nil
}

func (r *UserRepository) CreateOAuthUser(ctx context.Context, email, name, googleID string) (*entity.User, error) {
	username := generateUsername(email)

	createdUser, err := r.queries.CreateOAuthUser(ctx, database.CreateOAuthUserParams{
		Username: username,
		Email:    email,
		Name:     pgtype.Text{String: name, Valid: name != ""},
		GoogleID: pgtype.Text{String: googleID, Valid: true},
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, repoErr.ErrAlreadyExists
		}

		return nil, fmt.Errorf("create google user query failed: %w", err)
	}

	user := &entity.User{
		ID:       int(createdUser.ID),
		Username: createdUser.Username,
		Email:    createdUser.Email,
		IsActive: createdUser.IsActive,
	}

	if createdUser.CreatedAt.Valid {
		user.CreatedAt = createdUser.CreatedAt.Time
	}

	if createdUser.Name.Valid {
		user.Name = createdUser.Name.String
	}

	if createdUser.GoogleID.Valid {
		user.GoogleID = createdUser.GoogleID.String
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoErr.ErrNotFound
		}

		return nil, fmt.Errorf("get user by email query failed: %w", err)
	}

	return databaseToDomainUser(&user), nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	dbUsers, err := r.queries.GetAllUsers(ctx, database.GetAllUsersParams{
		Limit:  50,
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("get all users query failed: %w", err)
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

		return nil, fmt.Errorf("get user by id query failed: %w", err)
	}

	return databaseToDomainUser(&user), nil
}

func (r *UserRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	user, err := r.queries.GetUserByGoogleID(ctx, pgtype.Text{String: googleID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoErr.ErrNotFound
		}

		return nil, fmt.Errorf("get user by google id query failed: %w", err)
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
		return fmt.Errorf("update user profile query failed: %w", err)
	}

	if rows == 0 {
		return repoErr.ErrNotFound
	}

	return nil
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, id int, password string) error {
	rows, err := r.queries.UpdateUserPassword(ctx, database.UpdateUserPasswordParams{
		ID:           int32(id),
		PasswordHash: pgtype.Text{String: password, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("update user password query failed: %w", err)
	}

	if rows == 0 {
		return repoErr.ErrNotFound
	}

	return nil
}

func (r *UserRepository) UpdateUserIsActiveReturning(ctx context.Context, id int, isActive bool) (*entity.User, error) {
	rows, err := r.queries.UpdateUserIsActiveReturning(ctx, database.UpdateUserIsActiveReturningParams{
		ID:       int32(id),
		IsActive: isActive,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoErr.ErrNotFound
		}

		return nil, fmt.Errorf("update user is active query failed: %w", err)
	}

	return databaseToDomainUser(&rows), nil
}

func (r *UserRepository) UpdateUserGoogleID(ctx context.Context, userID int, googleID string) error {
	rows, err := r.queries.UpdateUserGoogleID(ctx, database.UpdateUserGoogleIDParams{
		ID:       int32(userID),
		GoogleID: pgtype.Text{String: googleID, Valid: true},
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repoErr.ErrGoogleExists
		}

		return fmt.Errorf("update user google id active query failed: %w", err)
	}

	if rows == 0 {
		return repoErr.ErrGoogleExists
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	rows, err := r.queries.DeleteUser(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("delete user query failed: %w", err)
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
		return nil, fmt.Errorf("create user token query failed: %w", err)
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

// TODO: for password reset
func (r *UserRepository) GetTokenByHash(ctx context.Context, tokenHash string, purpose entity.TokenPurpose) (*entity.UserToken, error) {
	row, err := r.queries.GetTokenByHash(ctx, database.GetTokenByHashParams{
		TokenHash: tokenHash,
		Purpose:   database.TokenPurpose(purpose),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoErr.ErrNotFound
		}

		return nil, fmt.Errorf("get token by hash query failed: %w", err)
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

// TODO: for resend verification token
func (r *UserRepository) DeleteByUserAndPurpose(ctx context.Context, userID int, purpose entity.TokenPurpose) error {
	return r.queries.DeleteByUserAndPurpose(ctx, database.DeleteByUserAndPurposeParams{
		UserID:  int32(userID),
		Purpose: database.TokenPurpose(purpose),
	})
}

func (r *UserRepository) DeleteOldTokens(ctx context.Context) (int64, error) {
	rows, err := r.queries.DeleteOldTokens(ctx)

	if err != nil {
		return 0, fmt.Errorf("delete old tokens query failed: %w", err)
	}

	return rows, nil
}
