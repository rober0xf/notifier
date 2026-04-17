package user_test

import (
	"context"
	"time"

	"github.com/rober0xf/notifier/internal/domain/entity"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type MockUserRepository struct {
	users  map[int]*entity.User    // key: id
	emails map[string]*entity.User // key: email
	tokens map[string]*entity.UserToken
	err    error
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	if _, exists := m.emails[user.Email]; exists {
		return nil, repoErr.ErrAlreadyExists // return the repo error
	}

	if user.ID == 0 {
		user.ID = len(m.users) + 1
	}

	m.emails[user.Email] = user
	m.users[user.ID] = user

	return user, nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	user, exists := m.users[id]
	if !exists {
		return nil, repoErr.ErrNotFound // we return the repo error, not the service error
	}

	user_copy := *user
	return &user_copy, nil
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	user, exists := m.emails[email]
	if !exists {
		return nil, repoErr.ErrNotFound
	}

	user_copy := *user
	return &user_copy, nil
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	users := make([]entity.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, *user)
	}

	return users, nil
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int) error {
	if m.err != nil {
		return m.err
	}

	user, ok := m.users[id]
	if !ok {
		return repoErr.ErrNotFound
	}

	delete(m.emails, user.Email)
	delete(m.users, id)

	return nil
}

func (m *MockUserRepository) UpdateUserProfile(ctx context.Context, id int, username string, email string) error {
	if m.err != nil {
		return m.err
	}

	user, ok := m.users[id]
	if !ok {
		return repoErr.ErrNotFound
	}

	if user.Email != email {
		existing, ok := m.emails[email]
		if ok && existing.ID != id {
			return repoErr.ErrAlreadyExists
		}
		delete(m.emails, user.Email) // clean old email
		m.emails[email] = user       // set the new
	}

	user.Email = email
	user.Username = username

	return nil
}

func (m *MockUserRepository) UpdateUserPassword(ctx context.Context, id int, hashedPassword string) error {
	if m.err != nil {
		return m.err
	}

	user, ok := m.users[id]
	if !ok {
		return repoErr.ErrNotFound
	}

	user.PasswordHash = hashedPassword

	return nil
}

func (m *MockUserRepository) UpdateUserIsActiveReturning(ctx context.Context, id int, active bool) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	user, ok := m.users[id]
	if !ok {
		return nil, repoErr.ErrNotFound
	}

	user.IsActive = active

	return user, nil
}

func (m *MockUserRepository) VerifyToken(ctx context.Context, tokenHash string, purpose entity.TokenPurpose) (*entity.UserToken, error) {
	if m.err != nil {
		return nil, m.err
	}

	token, ok := m.tokens[tokenHash]
	if !ok {
		return nil, repoErr.ErrNotFound
	}

	if token.Purpose != purpose {
		return nil, repoErr.ErrNotFound
	}

	if token.Used || time.Now().After(token.ExpiresAt) {
		return nil, repoErr.ErrNotFound
	}
	token.Used = true

	return token, nil
}

func (m *MockUserRepository) CreateUserToken(ctx context.Context, token *entity.UserToken) (*entity.UserToken, error) {
	if m.err != nil {
		return nil, m.err
	}

	if m.tokens == nil {
		m.tokens = make(map[string]*entity.UserToken)
	}

	m.tokens[token.TokenHash] = &entity.UserToken{
		ID:        token.ID,
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		Purpose:   token.Purpose,
		Used:      token.Used,
		ExpiresAt: token.ExpiresAt,
		CreatedAt: token.CreatedAt,
	}

	return m.tokens[token.TokenHash], nil
}

func (*MockUserRepository) CreateOAuthUser(ctx context.Context, email, name, googleID string) (*entity.User, error) {
	return nil, nil
}

func (*MockUserRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	return nil, nil
}

func (*MockUserRepository) UpdateUserGoogleID(ctx context.Context, userID int, googleID string) error {
	return nil
}

func (m *MockUserRepository) DeleteOldTokens(ctx context.Context) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}

	var deleted int64
	for k, t := range m.tokens {
		if t.Used || time.Now().After(t.ExpiresAt) {
			delete(m.tokens, k)
			deleted++
		}
	}

	return deleted, nil
}

func (m *MockUserRepository) GetTokenByHash(ctx context.Context, tokenHash string, purpose entity.TokenPurpose) (*entity.UserToken, error) {
	if m.err != nil {
		return nil, m.err
	}

	token, ok := m.tokens[tokenHash]
	if !ok {
		return nil, repoErr.ErrNotFound
	}

	if token.Used || time.Now().After(token.ExpiresAt) {
		return nil, repoErr.ErrNotFound
	}

	return token, nil
}
