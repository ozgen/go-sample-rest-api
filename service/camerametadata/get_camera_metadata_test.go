package camerametadata

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go-sample-rest-api/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_GetCameraMetaData(t *testing.T) {
	t.Run("GetCameraMetaData_withValidData_returnOk", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
		}

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}", handler.GetCameraMetaData).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("GetCameraMetaData_withCamId_returnNotFound", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(nil, fmt.Errorf("Not Found"))

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}", handler.GetCameraMetaData).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("GetCameraMetaData_withMalformedCamId_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := "123"

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}", handler.GetCameraMetaData).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}
