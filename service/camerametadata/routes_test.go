package camerametadata

import (
	"bytes"
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

func TestCameraMetadataService_CreateCameraMetaData(t *testing.T) {

	t.Run("CreateCameraMetaData_withValidData_returnCreated", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		handler := NewHandler(mockCameraStore)

		payload := types.CameraMetadataPayload{
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
		}

		expectedCamera := types.CameraMetadata{
			CamID:           uuid.New().String(),
			CameraName:      payload.CameraName,
			FirmwareVersion: payload.FirmwareVersion,
			CreatedAt:       time.Now(),
		}

		var capturedArg types.CameraMetadata
		mockCameraStore.On("CreateCameraMetaData", mock.AnythingOfType("types.CameraMetadata")).Run(func(args mock.Arguments) {
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
		router.HandleFunc("/camera_metadata", handler.CreateCameraMetaData).Methods(http.MethodPost)
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
			t.Errorf("expected CameraName %s, got %s", payload.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
	})
	t.Run("CreateCameraMetaData_withMalformedData_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		handler := NewHandler(mockCameraStore)

		// Act
		req, err := http.NewRequest(http.MethodPost, "/camera_metadata", bytes.NewBufferString("{invalid json"))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata", handler.CreateCameraMetaData).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
	t.Run("CreateCameraMetaData_withInvalidJsonData_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		handler := NewHandler(mockCameraStore)

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
		router.HandleFunc("/camera_metadata", handler.CreateCameraMetaData).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
	t.Run("CreateCameraMetaData_withDBError_returnInternalError", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		handler := NewHandler(mockCameraStore)

		payload := types.CameraMetadataPayload{
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
		}

		mockCameraStore.On("CreateCameraMetaData", mock.AnythingOfType("types.CameraMetadata")).Return(nil, fmt.Errorf("DB error"))

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
		router.HandleFunc("/camera_metadata", handler.CreateCameraMetaData).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})
}

type mockCameraStore struct {
	mock.Mock
}

func (m *mockCameraStore) CreateCameraMetaData(c types.CameraMetadata) (*types.CameraMetadata, error) {
	args := m.Called(c)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.CameraMetadata), args.Error(1)
}
