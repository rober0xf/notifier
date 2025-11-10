package users

import (
	"context"
	"strconv"
	"testing"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/domain"
	"github.com/stretchr/testify/assert"
)

/*  UPDATE PROFILE */

func (m *MockUserRepository) UpdateUserProfile(ctx context.Context, id int, username, email string) error {
	if m.err != nil {
		return m.err
	}

	existing_user, err := m.GetUserByID(ctx, id)
	if err != nil {
		return dto.ErrNotFound
	}
	old_email := existing_user.Email

	if old_email != email {
		other_user, err := m.GetUserByEmail(ctx, email)
		if err == nil && other_user.ID != id {
			return dto.ErrAlreadyExists
		}
		delete(m.users, old_email)
	}
	existing_user.Username = username
	existing_user.Email = email
	m.users[email] = existing_user
	m.users[strconv.Itoa(id)] = existing_user

	return nil
}

func TestUpdateUserProfile(t *testing.T) {
	t.Run("updated user profile", func(t *testing.T) {
		service, _ := setup_test_user_service(t)

		// original data
		original_username := "original username"
		original_email := "original@gmail.com"
		password := "password123#!%"

		created_user, err := service.Create(t.Context(), original_username, original_email, password)
		assert.NoError(t, err)

		// updated data
		changed_username := "username changed"
		changed_email := "changed@test.com"
		updated_data := &domain.User{
			ID:       created_user.ID,
			Username: changed_username,
			Email:    changed_email,
			Password: created_user.Password,
		}

		_, err = service.Update(context.Background(), updated_data)
		assert.NoError(t, err)

		updated_user, err := service.GetByID(context.Background(), created_user.ID)
		assert.NoError(t, err)
		assert.Equal(t, changed_username, updated_user.Username)
		assert.Equal(t, changed_email, updated_user.Email)

		_, err = service.GetByEmail(context.Background(), original_email)
		assert.Error(t, err)
		assert.Equal(t, dto.ErrUserNotFound, err)
	})
}

/* END UPDATE PROFILE */

/* UPDATE PASSWORD */
func (m *MockUserRepository) UpdateUserPassword(ctx context.Context, id int, hashed_password string) error {
	if m.err != nil {
		return m.err
	}
	user, err := m.GetUserByID(ctx, id)
	if err != nil {
		return dto.ErrNotFound
	}
	user.Password = hashed_password

	return nil
}

func TestUpdateUserPassword(t *testing.T) {
	t.Run("update password", func(t *testing.T) {
		service, _ := setup_test_user_service(t)

		username := "testuser"
		email := "test@gmail.com"
		password := "password123#!%"

		created_user, err := service.Create(context.Background(), username, email, password)
		assert.NoError(t, err)
		assert.NotNil(t, created_user)

		new_password := "newpassword456$&*"
		update_data := &domain.User{
			ID:       created_user.ID,
			Username: created_user.Username,
			Email:    created_user.Email,
			Password: new_password,
		}
		updated_user, err := service.Update(context.Background(), update_data)
		assert.NoError(t, err)

		assert.NotEqual(t, created_user.Password, updated_user.Password)
	})
}

/* END UPDATE PASSWORD */

/*  UPDATE ACTIVE */
func (m *MockUserRepository) UpdateUserActive(ctx context.Context, id int, active bool) error {
	if m.err != nil {
		return m.err
	}

	existing_user, err := m.GetUserByID(ctx, id)
	if err != nil {
		return dto.ErrNotFound
	}
	existing_user.Active = active
	m.users[strconv.Itoa(id)] = existing_user
	m.users[existing_user.Email] = existing_user

	return nil
}

func TestGetVerificationEmail(t *testing.T) {
	t.Run("activates when a user verifies", func(t *testing.T) {
		service, _ := setup_test_user_service(t)

		username := "testuser"
		email := "test@gmail.com"
		password := "password123#!%"

		created_user, err := service.Create(context.Background(), username, email, password)
		assert.NoError(t, err)
		assert.NotNil(t, created_user)

		assert.False(t, created_user.Active)

		verified_user, err := service.GetVerificationEmail(context.Background(), email)
		assert.NoError(t, err)

		assert.True(t, verified_user.Active)

		fetched_user, err := service.GetByID(context.Background(), created_user.ID)
		assert.NoError(t, err)
		assert.True(t, fetched_user.Active)
		assert.Equal(t, email, fetched_user.Email)
		assert.Equal(t, username, fetched_user.Username)
	})
}

/* END UPDATE ACTIVE */
