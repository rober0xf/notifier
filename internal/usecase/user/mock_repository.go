package user

import (
	"context"
	"strconv"

	"github.com/rober0xf/notifier/internal/domain/entity"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type MockUserRepository struct {
	users map[string]*entity.User // key: id or email
	err   error
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if m.err != nil {
		return m.err
	}

	if _, exists := m.users[user.Email]; exists {
		return repoErr.ErrAlreadyExists // return the repo error
	}

	if user.ID == 0 {
		user.ID = len(m.users)/2 + 1
	}

	m.users[user.Email] = user
	m.users[strconv.Itoa(user.ID)] = user

	return nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	user, exists := m.users[strconv.Itoa(id)]
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

	user, exists := m.users[email]
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

	seen := make(map[int]bool)
	users := make([]entity.User, 0)

	for key, user := range m.users {
		if _, err := strconv.Atoi(key); err == nil {
			if !seen[user.ID] {
				users = append(users, *user)
				seen[user.ID] = true
			}
		}
	}

	return users, nil
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int) error {
	if m.err != nil {
		return m.err
	}

	user, exists := m.users[strconv.Itoa(id)]
	if !exists {
		return repoErr.ErrNotFound
	}

	delete(m.users, user.Email)
	delete(m.users, strconv.Itoa(id))

	return nil
}

func (m *MockUserRepository) UpdateUserProfile(ctx context.Context, id int, username string, email string) error {
	if m.err != nil {
		return m.err
	}

	existingUser, err := m.GetUserByID(ctx, id)
	if err != nil {
		return repoErr.ErrNotFound
	}
	oldMail := existingUser.Email

	if oldMail != email {
		otherUser, err := m.GetUserByEmail(ctx, email)
		if err == nil && otherUser.ID != id {
			return repoErr.ErrAlreadyExists
		}
		delete(m.users, oldMail)
	}

	existingUser.Username = username
	existingUser.Email = email

	m.users[email] = existingUser
	m.users[strconv.Itoa(id)] = existingUser

	return nil
}

func (m *MockUserRepository) UpdateUserPassword(ctx context.Context, id int, hashedPassword string) error {
	if m.err != nil {
		return m.err
	}

	user, err := m.GetUserByID(ctx, id)
	if err != nil {
		return repoErr.ErrNotFound
	}
	user.Password = hashedPassword

	m.users[user.Email] = user
	m.users[strconv.Itoa(id)] = user

	return nil
}

func (m *MockUserRepository) UpdateUserActive(ctx context.Context, id int, active bool) error {
	if m.err != nil {
		return m.err
	}

	existingUser, err := m.GetUserByID(ctx, id)
	if err != nil {
		return repoErr.ErrNotFound
	}

	existingUser.Active = active

	m.users[strconv.Itoa(id)] = existingUser
	m.users[existingUser.Email] = existingUser

	return nil
}
