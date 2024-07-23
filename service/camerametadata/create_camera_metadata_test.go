package camerametadata

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"go-sample-rest-api/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_CreateCameraMetadata(t *testing.T) {
	t.Run("CreateCameraMetadata_withValidData_returnCreated", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		payload := types.CameraMetadataPayload{
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
		}
		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		expectedCamera := types.CameraMetadata{
			CamID:           uuid.New().String(),
			CameraName:      payload.CameraName,
			FirmwareVersion: payload.FirmwareVersion,
			CreatedAt:       nullTime,
		}

		var capturedArg types.CameraMetadata
		mockCameraStore.On("CreateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Run(func(args mock.Arguments) {
			capturedArg = args.Get(0).(types.CameraMetadata)
		}).Return(&expectedCamera, nil)

		cameraData, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		req, err := http.NewRequest(http.MethodPost, "/camera_metadata", bytes.NewBuffer(cameraData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata", handler.CreateCameraMetadata).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}

		// Use captured argument to assert values
		if capturedArg.CameraName != payload.CameraName {
			t.Errorf("expected CameraName %s, got %s", payload.CameraName, capturedArg.CameraName)
		}
		if capturedArg.FirmwareVersion != payload.FirmwareVersion {
			t.Errorf("expected FirmwareVersion %s, got %s", payload.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("CreateCameraMetadata_withMalformedData_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)
		// Act
		req, err := http.NewRequest(http.MethodPost, "/camera_metadata", bytes.NewBufferString("{invalid json"))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata", handler.CreateCameraMetadata).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("CreateCameraMetadata_withInvalidJsonData_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		payload := types.CameraMetadataPayload{
			CameraName:      "",
			FirmwareVersion: "",
		}

		cameraData, _ := json.Marshal(payload)

		// Act
		req, err := http.NewRequest(http.MethodPost, "/camera_metadata", bytes.NewBuffer(cameraData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata", handler.CreateCameraMetadata).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("CreateCameraMetadata_withDBError_returnInternalError", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		payload := types.CameraMetadataPayload{
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
		}

		mockCameraStore.On("CreateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Return(nil, fmt.Errorf("DB error"))

		cameraData, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		req, err := http.NewRequest(http.MethodPost, "/camera_metadata", bytes.NewBuffer(cameraData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata", handler.CreateCameraMetadata).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})
}
