package camerametadata

import (
	"bytes"
	"context"
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
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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
			t.Errorf("expected CameraName %s, got %s", payload.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("CreateCameraMetadata_withMalformedData_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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

func TestHandler_InitializeCameraMetaData(t *testing.T) {
	t.Run("InitializeCameraMetaData_withValidData_returnOk", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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
			t.Errorf("expected CameraName %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("InitializeCameraMetaData_withAlreadyInitializedData_returnConflict", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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

func TestHandler_GetCameraMetaData(t *testing.T) {
	t.Run("GetCameraMetaData_withValidData_returnOk", func(t *testing.T) {
		//arrange
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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
		mockCameraStore := new(mockCameraStore)
		mockAzureStorage := new(mockAzureStorage)
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

// TODO 22-07-24 - ozgen : write UploadImage unit tests

type mockCameraStore struct {
	mock.Mock
}

func (m *mockCameraStore) CreateCameraMetadata(c types.CameraMetadata) (*types.CameraMetadata, error) {
	args := m.Called(c)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.CameraMetadata), args.Error(1)
}

func (m *mockCameraStore) UpdateCameraMetadata(c types.CameraMetadata) (*types.CameraMetadata, error) {
	args := m.Called(c)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.CameraMetadata), args.Error(1)
}
func (m *mockCameraStore) GetCameraMetadataByID(c string) (*types.CameraMetadata, error) {
	args := m.Called(c)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.CameraMetadata), args.Error(1)
}

type mockAzureStorage struct {
	mock.Mock
}

func (m *mockAzureStorage) UploadImage(ctx context.Context, blobName string, imageData []byte) error {
	args := m.Called(ctx, blobName, imageData)
	return args.Error(0)
}

func (m *mockAzureStorage) DownloadImage(ctx context.Context, blobName string) ([]byte, error) {
	args := m.Called(ctx, blobName)
	return args.Get(0).([]byte), args.Error(1)
}
