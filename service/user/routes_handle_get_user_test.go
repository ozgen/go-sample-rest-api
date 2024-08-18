package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"go-sample-rest-api/types"
	"go-sample-rest-api/utils"
)

func TestUserService_Handle_GetUser(t *testing.T) {
	t.Run("missing user ID in URL", func(t *testing.T) {
		// arrange
		handler := NewHandler(new(mockUserStore), new(MockAuthenticator))
		req, _ := http.NewRequest(http.MethodGet, "/users/", nil) // Missing userID in the URL
		rr := httptest.NewRecorder()

		// act
		handler.handleGetUser(rr, req)

		// assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("invalid user ID", func(t *testing.T) {
		// arrange
		handler := NewHandler(new(mockUserStore), new(MockAuthenticator))
		router := mux.NewRouter()
		router.HandleFunc("/users/{userID}", handler.handleGetUser)
		req, _ := http.NewRequest(http.MethodGet, "/users/abc", nil) // 'abc' is not a valid user ID
		rr := httptest.NewRecorder()

		// act
		router.ServeHTTP(rr, req)

		// assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		// arrange
		mockUserStore := new(mockUserStore)
		handler := NewHandler(mockUserStore, new(MockAuthenticator))
		router := mux.NewRouter()
		router.HandleFunc("/users/{userID}", handler.handleGetUser)

		mockUserStore.On("GetUserByID", 1).Return(nil, fmt.Errorf("user not found"))
		req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
		rr := httptest.NewRecorder()

		// act
		router.ServeHTTP(rr, req)

		// assert
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})

	t.Run("successful retrieval", func(t *testing.T) {
		// arrange
		mockUserStore := new(mockUserStore)
		handler := NewHandler(mockUserStore, new(MockAuthenticator))
		router := mux.NewRouter()
		router.HandleFunc("/users/{userID}", handler.handleGetUser)

		expectedUser := &types.User{ID: 1, FirstName: "John", LastName: "Doe"}
		mockUserStore.On("GetUserByID", 1).Return(expectedUser, nil)
		req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
		rr := httptest.NewRecorder()

		// act
		router.ServeHTTP(rr, req)

		// assert
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		var user types.User
		if err := utils.ParseJSONResponse(rr.Result(), &user); err != nil {
			t.Fatal("failed to parse response body")
		}
		if user.ID != expectedUser.ID {
			t.Errorf("expected user ID %d, got %d", expectedUser.ID, user.ID)
		}
	})
}
