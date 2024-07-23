package api_test

import (
	"go-sample-rest-api/cmd/api"
	"go-sample-rest-api/logging"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIServer_Run(t *testing.T) {
	mockDB := new(api.MockDB)
	mockAzureStorage := new(api.MockAzureStorage)

	server := api.NewAPIServer("localhost:8080", mockDB, mockAzureStorage)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/v1/some-route", nil)

	// Start the server in a goroutine
	go func() {
		if err := server.Run(); err != nil {
			logging.GetLogger().Error("Server failed to run")
		}
	}()

	// Make requests to the server
	http.DefaultServeMux.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	mockDB.AssertExpectations(t)
	mockAzureStorage.AssertExpectations(t)
}
