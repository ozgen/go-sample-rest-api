package camerametadata

import (
	"database/sql"
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

func TestHandler_InitializeCameraMetaData(t *testing.T) {
	t.Run("InitializeCameraMetaData_withValidData_returnOk", func(t *testing.T) {
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

		var capturedArg types.CameraMetadata
		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockCameraStore.On("UpdateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Run(func(args mock.Arguments) {
			capturedArg = args.Get(0).(types.CameraMetadata)
		}).Return(&expectedCamera, nil)
		// Act
		req, err := http.NewRequest(http.MethodPatch, "/camera_metadata/"+camID+"/init", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/init", handler.InitializeCameraMetaData).Methods(http.MethodPatch)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		if capturedArg.CameraName != expectedCamera.CameraName {
			t.Errorf("expected CameraName %s, got %s", expectedCamera.CameraName, capturedArg.CameraName)
		}
		if capturedArg.FirmwareVersion != expectedCamera.FirmwareVersion {
			t.Errorf("expected FirmwareVersion %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("InitializeCameraMetaData_withAlreadyInitializedData_returnConflict", func(t *testing.T) {
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
			InitializedAt:   nullTime,
		}

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)

		// Act
		req, err := http.NewRequest(http.MethodPatch, "/camera_metadata/"+camID+"/init", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/init", handler.InitializeCameraMetaData).Methods(http.MethodPatch)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusConflict {
			t.Errorf("expected status code %d, got %d", http.StatusConflict, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("InitializeCameraMetaData_withErrOnUpdate_returnInternalError", func(t *testing.T) {
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
		mockCameraStore.On("UpdateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Return(nil, fmt.Errorf("Failed"))
		// Act
		req, err := http.NewRequest(http.MethodPatch, "/camera_metadata/"+camID+"/init", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/init", handler.InitializeCameraMetaData).Methods(http.MethodPatch)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("InitializeCameraMetaData_withWrongCamID_returnNotFound", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(nil, fmt.Errorf("Not Found"))
		// Act
		req, err := http.NewRequest(http.MethodPatch, "/camera_metadata/"+camID+"/init", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/init", handler.InitializeCameraMetaData).Methods(http.MethodPatch)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("InitializeCameraMetaData_withInvalidCamID_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := "123"

		// Act
		req, err := http.NewRequest(http.MethodPatch, "/camera_metadata/"+camID+"/init", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/init", handler.InitializeCameraMetaData).Methods(http.MethodPatch)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("InitializeCameraMetaData_withEmptyCamID_returnStatusMovedPermanently", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		// Act
		req, err := http.NewRequest(http.MethodPatch, "/camera_metadata//init", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/init", handler.InitializeCameraMetaData).Methods(http.MethodPatch)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusMovedPermanently {
			t.Errorf("expected status code %d, got %d", http.StatusMovedPermanently, rr.Code)
		}
	})
}
