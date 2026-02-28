package user_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	t.Run("succesfully found all users", func(t *testing.T) {
		uc, mockRepo := setupGetAllUsersTest(t)

		users := []*entity.User{
			{ID: 1, Email: "user1@test.com", Username: "user1", Active: true},
			{ID: 2, Email: "user2@test.com", Username: "user2", Active: true},
			{ID: 3, Email: "user3@test.com", Username: "user3", Active: false},
		}

		for _, u := range users {
			mockRepo.users[u.Email] = u
			mockRepo.users[strconv.Itoa(u.ID)] = u
		}

		result, err := uc.Execute(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)

		userIDs := make(map[int]bool)
		for _, u := range result {
			userIDs[u.ID] = true
		}

		assert.True(t, userIDs[1])
		assert.True(t, userIDs[2])
		assert.True(t, userIDs[3])
	})

	t.Run("returns empty list when no users exist", func(t *testing.T) {
		uc, _ := setupGetAllUsersTest(t)

		users, err := uc.Execute(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Empty(t, users) // empty, not nil
	})

	t.Run("returns users with different status", func(t *testing.T) {
		uc, mockRepo := setupGetAllUsersTest(t)

		activeUser := &entity.User{ID: 1, Email: "active@test.com", Username: "active", Active: true}
		inactiveUser := &entity.User{ID: 2, Email: "inactive@test.com", Username: "inactive", Active: false}

		mockRepo.users["active@test.com"] = activeUser
		mockRepo.users["1"] = activeUser
		mockRepo.users["inactive@test.com"] = inactiveUser
		mockRepo.users["2"] = inactiveUser

		users, err := uc.Execute(context.Background())

		assert.NoError(t, err)
		assert.Len(t, users, 2)

		var hasActive, hasInactive bool
		for _, u := range users {
			if u.Active {
				hasActive = true
			} else {
				hasInactive = true
			}
		}
		assert.True(t, hasActive)
		assert.True(t, hasInactive)
	})

	t.Run("handles repository errors", func(t *testing.T) {
		uc, mockRepo := setupGetAllUsersTest(t)

		mockRepo.err = errors.New("database connection failed")

		users, err := uc.Execute(context.Background())

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Contains(t, err.Error(), "database connection failed")
	})
}
