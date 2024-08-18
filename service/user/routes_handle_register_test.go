package user

import (
	"bytes"
	"fmt"
	"go-sample-rest-api/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestUserService_Handle_Register(t *testing.T) {
	t.Run("invalid JSON payload", func(t *testing.T) {
		// arrange
		handler := NewHandler(new(mockUserStore), new(MockAuthenticator))

		// act
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("{invalid json"))
		rr := httptest.NewRecorder()
		handler.handleRegister(rr, req)

		// assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("invalid payload data", func(t *testing.T) {
		// arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)
		handler := NewHandler(mockUserStore, mockAuth)

		userData := `{"email": "test@test.com", "password": ""}` // Invalid because password is empty
		req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(userData))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// act
		handler.handleRegister(rr, req)

		// assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("user already exists", func(t *testing.T) {
		// arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)
		handler := NewHandler(mockUserStore, mockAuth)

		userData := `{"firstName": "John", "lastName": "Doe", "email": "exists@test.com", "password": "password123"}`
		mockUserStore.On("GetUserByEmail", "exists@test.com").Return(&types.User{}, nil) // User already exists

		req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(userData))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// act
		handler.handleRegister(rr, req)

		// assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		mockUserStore.AssertExpectations(t)
	})

	t.Run("error during user creation", func(t *testing.T) {
		// arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)
		handler := NewHandler(mockUserStore, mockAuth)

		userData := `{"firstName": "John", "lastName": "Doe", "email": "john@doe.com", "password": "password123"}`
		mockUserStore.On("GetUserByEmail", "john@doe.com").Return(nil, fmt.Errorf("not found"))
		mockAuth.On("HashPassword", "password123").Return("hashedPassword123", nil)
		mockUserStore.On("CreateUser", mock.Anything).Return(fmt.Errorf("db error"))

		req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(userData))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// act
		handler.handleRegister(rr, req)

		// assert
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		mockUserStore.AssertExpectations(t)
		mockAuth.AssertExpectations(t)

	})

	t.Run("error during hashing password", func(t *testing.T) {
		// arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)
		handler := NewHandler(mockUserStore, mockAuth)

		userData := `{"firstName": "John", "lastName": "Doe", "email": "john@doe.com", "password": "password123"}`
		mockUserStore.On("GetUserByEmail", "john@doe.com").Return(nil, fmt.Errorf("not found"))
		mockAuth.On("HashPassword", "password123").Return(nil, fmt.Errorf("hash error"))

		req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(userData))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// act
		handler.handleRegister(rr, req)

		// assert
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		mockUserStore.AssertExpectations(t)
		mockAuth.AssertExpectations(t)

	})

	t.Run("successful registration", func(t *testing.T) {
		// arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)
		handler := NewHandler(mockUserStore, mockAuth)

		userData := `{"firstName": "John", "lastName": "Doe", "email": "new@user.com", "password": "securePass123"}`
		mockUserStore.On("GetUserByEmail", "new@user.com").Return(nil, fmt.Errorf("not found")) // No existing user
		mockAuth.On("HashPassword", "securePass123").Return("hashedPassword123", nil)
		mockUserStore.On("CreateUser", mock.Anything).Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(userData))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// act
		handler.handleRegister(rr, req)

		// assert
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		mockUserStore.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})
}
