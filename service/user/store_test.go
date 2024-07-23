package user

import (
	_ "database/sql"
	_ "errors"
	"github.com/DATA-DOG/go-sqlmock"
	db2 "go-sample-rest-api/db"
	"go-sample-rest-api/types"
	"reflect"
	"testing"
	"time"
)

func setupMockDB(t *testing.T) (*db2.SQLDB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	cleanup := func() {
		db.Close()
	}

	sqldb := db2.NewSQLDB(db)
	return sqldb, mock, cleanup
}

func TestStore_CreateUser(t *testing.T) {
	// arrange
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	store := NewStore(db)
	user := types.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
		Password:  "securepassword",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.FirstName, user.LastName, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// act
	err := store.CreateUser(user)

	// assert
	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStore_GetUserByEmail(t *testing.T) {
	// arrange
	db, mock, cleanup := setupMockDB(t)

	defer cleanup()

	store := NewStore(db)
	email := "johndoe@example.com"
	expectedUser := &types.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     email,
		Password:  "securepassword",
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password", "createdAt"}).
		AddRow(expectedUser.ID, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email, expectedUser.Password, expectedUser.CreatedAt)

	mock.ExpectQuery("SELECT \\* FROM users WHERE email =").
		WithArgs(email).
		WillReturnRows(rows)

	// act
	user, err := store.GetUserByEmail(email)

	// assert
	if err != nil {
		t.Errorf("error was not expected while getting user: %s", err)
	}

	if !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("expected user does not match the returned user. Expected: %v, Got: %v", expectedUser, user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestStore_GetUserByID(t *testing.T) {
	// arrange
	db, mock, cleanup := setupMockDB(t)

	defer cleanup()

	store := NewStore(db)
	ID := 1
	expectedUser := &types.User{
		ID:        ID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
		Password:  "securepassword",
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password", "createdAt"}).
		AddRow(expectedUser.ID, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email, expectedUser.Password, expectedUser.CreatedAt)

	mock.ExpectQuery("SELECT \\* FROM users WHERE id =").
		WithArgs(ID).
		WillReturnRows(rows)

	// act
	user, err := store.GetUserByID(ID)

	// assert
	if err != nil {
		t.Errorf("error was not expected while getting user: %s", err)
	}

	if !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("expected user does not match the returned user. Expected: %v, Got: %v", expectedUser, user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
