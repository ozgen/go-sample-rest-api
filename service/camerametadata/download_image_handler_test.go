package camerametadata

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"go-sample-rest-api/types"
	"go-sample-rest-api/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_DownloadImageHandler(t *testing.T) {

	t.Run("DownloadImageHandler_withValidData_returnOk", func(t *testing.T) {
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
		imageID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
			InitializedAt:   nullTime,
			ImageId:         sql.NullString{String: imageID, Valid: true},
		}

		imageData, err := base64.StdEncoding.DecodeString(utils.NormalizeBase64(Base64Data))
		if err != nil {
			t.Fatal(err)
		}
		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockAzureStorage.On("DownloadImage",
			mock.AnythingOfType("*context.valueCtx"), imageID+".png").Return(imageData, nil)

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID+"/download_image", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "image/png")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/download_image", handler.DownloadImageHandler).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
		mockAzureStorage.AssertExpectations(t)
	})

	t.Run("DownloadImageHandler_withWriteFailure_returnsInternalServerError", func(t *testing.T) {
		// Arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()
		imageID := uuid.New().String()
		imageData := []byte("image data")

		expectedCamera := types.CameraMetadata{
			CamID:   camID,
			ImageId: sql.NullString{String: imageID, Valid: true},
		}

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockAzureStorage.On("DownloadImage", mock.Anything, imageID+".png").Return(imageData, nil)

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID+"/download_image", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()                      // Normal recorder
		fw := &FailWriter{ResponseWriter: rr, fail: true} // Wrap the recorder to simulate failure

		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/download_image", handler.DownloadImageHandler).Methods(http.MethodGet)
		router.ServeHTTP(fw, req) // Use the failing writer here

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("DownloadImageHandler_withDownLoadError_returnsInternalServerError", func(t *testing.T) {
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
		imageID := uuid.New().String()
		expectedCamera := types.CameraMetadata{
			CamID:           camID,
			CameraName:      "camera-name",
			FirmwareVersion: "v123",
			CreatedAt:       nullTime,
			InitializedAt:   nullTime,
			ImageId:         sql.NullString{String: imageID, Valid: true},
		}

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockAzureStorage.On("DownloadImage",
			mock.AnythingOfType("*context.valueCtx"), imageID+".png").Return(nil, fmt.Errorf("download err"))

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID+"/download_image", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "image/png")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/download_image", handler.DownloadImageHandler).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("DownloadImageHandler_withNullImageID_returnNotFound", func(t *testing.T) {
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
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID+"/download_image", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "image/png")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/download_image", handler.DownloadImageHandler).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("DownloadImageHandler_withWrongCameraID_returnNotFound", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(nil, fmt.Errorf("not found"))

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID+"/download_image", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "image/png")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/download_image", handler.DownloadImageHandler).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("DownloadImageHandler_withMalformedCameraID_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := "123"

		// Act
		req, err := http.NewRequest(http.MethodGet, "/camera_metadata/"+camID+"/download_image", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "image/png")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/download_image", handler.DownloadImageHandler).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})
}
