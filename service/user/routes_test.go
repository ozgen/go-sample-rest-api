package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"go-sample-rest-api/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserService(t *testing.T) {

	t.Run("should fail if the user ID is not a number", func(t *testing.T) {
		// arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)
		handler := NewHandler(mockUserStore, mockAuth)
		req, err := http.NewRequest(http.MethodGet, "/user/abc", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		//act
		router.HandleFunc("/user/{userID}", handler.handleGetUser).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		//assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should handle get user by ID", func(t *testing.T) {
		//arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)

		handler := NewHandler(mockUserStore, mockAuth)
		// Set up mock response
		expectedUser := &types.User{ID: 42, FirstName: "John Doe"}
		mockUserStore.On("GetUserByID", 42).Return(expectedUser, nil)

		//act
		req, err := http.NewRequest(http.MethodGet, "/user/42", nil)

		//assert
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}", handler.handleGetUser).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		// Additional check to assert the response content if necessary
		var user types.User
		if err := json.Unmarshal(rr.Body.Bytes(), &user); err != nil {
			t.Fatal("failed to unmarshal response")
		}
		if user.ID != expectedUser.ID {
			t.Errorf("expected user ID %d, got %d", expectedUser.ID, user.ID)
		}

		// Ensure mock expectations are met
		mockUserStore.AssertExpectations(t)
	})
	t.Run("should handle login with correct credentials", func(t *testing.T) {
		//arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)

		handler := NewHandler(mockUserStore, mockAuth)

		email := "test@test.com"
		user := types.LoginUserPayload{
			Email:    email,
			Password: "test123",
		}
		mockUser := &types.User{
			Email:    email,
			Password: "hashedPassword123", // assuming you are storing hashed passwords
			ID:       1,
		}

		expectedToken := "dummyToken123"

		// Correct mock setup
		mockAuth.On("ComparePasswords", mock.Anything, mock.Anything).Return(true)
		mockAuth.On("CreateJWT", mock.Anything, mock.Anything).Return(expectedToken, nil)
		mockUserStore.On("GetUserByEmail", email).Return(mockUser, nil)

		userData, err := json.Marshal(user)
		if err != nil {
			t.Fatal(err)
		}

		// act
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/login", handler.handleLogin).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		// Verify that all expectations were met
		mockAuth.AssertExpectations(t)
	})
	t.Run("invalid JSON payload", func(t *testing.T) {
		// arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)

		handler := NewHandler(mockUserStore, mockAuth)

		// act
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{invalid json"))
		rr := httptest.NewRecorder()
		handler.handleLogin(rr, req)

		// assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("validation errors on payload", func(t *testing.T) {
		// arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)

		handler := NewHandler(mockUserStore, mockAuth)
		mockUserStore.On("GetUserByEmail", mock.Anything).Return(nil, fmt.Errorf("not found"))

		user := types.LoginUserPayload{
			Email:    "", // Empty email to trigger validation error
			Password: "",
		}
		userData, _ := json.Marshal(user)
		// act
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userData))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.handleLogin(rr, req)

		// assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("user not found in database", func(t *testing.T) {
		// arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)

		handler := NewHandler(mockUserStore, mockAuth)

		testEmail := "nonexistent@test.com"
		mockUserStore.On("GetUserByEmail", mock.Anything).Return(types.User{}, fmt.Errorf("user not found"))

		// act
		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(fmt.Sprintf(`{"email":"%s","password":"password123"}`, testEmail)))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleLogin(rr, req)

		// assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("incorrect password", func(t *testing.T) {
		// arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)

		handler := NewHandler(mockUserStore, mockAuth)

		hashedPassword := "$2a$12$examplebcryptpasswordhash"
		mockUserStore.On("GetUserByEmail", "test@test.com").Return(&types.User{Password: hashedPassword}, nil)
		mockAuth.On("ComparePasswords", hashedPassword, []byte("wrongpassword")).Return(false)

		user := types.LoginUserPayload{
			Email:    "test@test.com",
			Password: "wrongpassword",
		}
		userData, _ := json.Marshal(user)

		// act
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userData))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleLogin(rr, req)

		// assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("failure to create JWT", func(t *testing.T) {
		//arrange
		mockUserStore := new(mockUserStore)
		mockAuth := new(MockAuthenticator)

		handler := NewHandler(mockUserStore, mockAuth)

		hashedPassword := "$2a$12$examplebcryptpasswordhash"
		mockUserStore.On("GetUserByEmail", "test@test.com").Return(&types.User{ID: 1, Password: hashedPassword}, nil)
		mockAuth.On("ComparePasswords", hashedPassword, []byte("password123")).Return(true)
		mockAuth.On("CreateJWT", mock.Anything, 1).Return("", fmt.Errorf("error creating token"))

		user := types.LoginUserPayload{
			Email:    "test@test.com",
			Password: "password123",
		}
		userData, _ := json.Marshal(user)

		// act
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userData))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleLogin(rr, req)

		// assert
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
	//TODO: 18-07-24 ozgen : add register tests
}

type mockUserStore struct {
	mock.Mock
}

func (m *mockUserStore) UpdateUser(u types.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	args := m.Called(email)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.User), args.Error(1)
}

func (m *mockUserStore) CreateUser(u types.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	args := m.Called(id)
	return args.Get(0).(*types.User), args.Error(1)
}

type MockAuthenticator struct {
	mock.Mock
}

func (m *MockAuthenticator) CreateJWT(secret []byte, userID int) (string, error) {
	args := m.Called(secret, userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthenticator) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthenticator) ComparePasswords(hashed string, plain []byte) bool {
	args := m.Called(hashed, plain)
	return args.Bool(0)
}
