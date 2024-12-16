package handlers_test

import (
	"bytes"
	"encoding/json"
	"goapi/handlers"
	"goapi/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, nil
}

func TestCreateUserWichMock(t *testing.T) {
	db, mock, err := setupMockDB()
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
