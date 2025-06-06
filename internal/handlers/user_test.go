package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rober0xf/notifier/internal/handlers"
	"github.com/rober0xf/notifier/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setup_mock_db() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func Test_create_user_with_mock(t *testing.T) {
	db, mock, err := setup_mock_db()
	if err != nil {
		t.Fatalf("could not setup mock: %v", err)
	}

	store := &handlers.Store{DB: db}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").WithArgs(
		sqlmock.AnyArg(),
		"test user",
		"test@example.com",
		sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := models.User{
		Name:     "test user",
		Email:    "test@example.com",
		Password: "password123",
	}

	jsonUser, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonUser))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	store.CreateUser(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status created. got: %v", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
