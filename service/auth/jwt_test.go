package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-sample-rest-api/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	panic("implement me")
}

func (m *mockUserStore) CreateUser(user types.User) error {
	panic("implement me")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return &types.User{ID: id}, nil
}

func TestCreateJWT(t *testing.T) {
	secret := []byte("secret")

	token, err := CreateJWT(secret, 1)
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if token == "" {
		t.Error("expected token to be not empty")
	}
}

func TestWithJWTAuth(t *testing.T) {
	mockStore := &mockUserStore{}

	// Handler that will be wrapped by the middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserIDFromContext(r.Context())
		if userID != 1 {
			t.Errorf("Expected user ID 1, got %d", userID)
		}
		w.WriteHeader(http.StatusOK) // Ensure a response code is set for successful requests
	})

	// Define test scenarios
	tests := []struct {
		name         string
		token        string
		setupFunc    func() jwtValidatorFunc
		expectedCode int
	}{
		{
			name:  "Valid Token",
			token: "valid_token",
			setupFunc: func() jwtValidatorFunc {
				return func(tokenString string) (*jwt.Token, error) {
					return &jwt.Token{
						Valid: true,
						Claims: jwt.MapClaims{
							"userID": "1",
						},
					}, nil
				}
			},
			expectedCode: http.StatusOK,
		},
		{
			name:  "Invalid Token",
			token: "invalid_token",
			setupFunc: func() jwtValidatorFunc {
				return func(tokenString string) (*jwt.Token, error) {
					return nil, fmt.Errorf("unexpected signing method: alg")
				}
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name:  "No Token",
			token: "",
			setupFunc: func() jwtValidatorFunc {
				return func(tokenString string) (*jwt.Token, error) {
					return nil, fmt.Errorf("no token")
				}
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name:  "Token With Invalid UserID",
			token: "bad_userid_token",
			setupFunc: func() jwtValidatorFunc {
				return func(tokenString string) (*jwt.Token, error) {
					return &jwt.Token{
						Valid: true,
						Claims: jwt.MapClaims{
							"userID": "abc", // Non-integer userID
						},
					}, nil
				}
			},
			expectedCode: http.StatusForbidden,
		},
	}

	// Run the test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			validateJWT = tc.setupFunc() // Setup the JWT validation function

			req, _ := http.NewRequest("GET", "/some-path", nil)
			if tc.token != "" {
				req.Header.Add("Authorization", "Bearer "+tc.token)
			}
			rr := httptest.NewRecorder()
			handler := WithJWTAuth(testHandler, mockStore)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("%s: expected HTTP status %v, got %v", tc.name, tc.expectedCode, status)
			}
		})
	}

	// Reset the validateJWT function to prevent side effects in other tests
	validateJWT = validateJWTDefault
}
