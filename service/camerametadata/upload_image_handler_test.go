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

func TestHandler_UploadImageHandler(t *testing.T) {

	t.Run("UploadImageHandler_withValidData_returnOk", func(t *testing.T) {
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
		}

		var capturedArg types.CameraMetadata
		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockCameraStore.On("UpdateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Run(func(args mock.Arguments) {
			capturedArg = args.Get(0).(types.CameraMetadata)
		}).Return(&expectedCamera, nil)
		mockAzureStorage.On("UploadImage",
			mock.AnythingOfType("*context.valueCtx"), imageID+".png", mock.AnythingOfType("[]uint8")).Return(nil)
		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + Base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
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
		if capturedArg.NameOfStoredPicture.String != imageID {
			t.Errorf("expected NameOfStoredPicture %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}
		if capturedArg.ImageId.String != imageID {
			t.Errorf("expected ImageId %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
		mockAzureStorage.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withUploadError_returnInternalServerError", func(t *testing.T) {
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
		}

		var capturedArg types.CameraMetadata
		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockCameraStore.On("UpdateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Run(func(args mock.Arguments) {
			capturedArg = args.Get(0).(types.CameraMetadata)
		}).Return(&expectedCamera, nil)
		mockAzureStorage.On("UploadImage",
			mock.AnythingOfType("*context.valueCtx"), imageID+".png", mock.AnythingOfType("[]uint8")).Return(fmt.Errorf("upload err"))
		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + Base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		if capturedArg.CameraName != expectedCamera.CameraName {
			t.Errorf("expected CameraName %s, got %s", expectedCamera.CameraName, capturedArg.CameraName)
		}
		if capturedArg.FirmwareVersion != expectedCamera.FirmwareVersion {
			t.Errorf("expected FirmwareVersion %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}
		if capturedArg.NameOfStoredPicture.String != imageID {
			t.Errorf("expected NameOfStoredPicture %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}
		if capturedArg.ImageId.String != imageID {
			t.Errorf("expected ImageId %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
		mockAzureStorage.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withUpdateCameraMetadataError_returnInternalServerError", func(t *testing.T) {
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
		}

		var capturedArg types.CameraMetadata
		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		mockCameraStore.On("UpdateCameraMetadata", mock.AnythingOfType("types.CameraMetadata")).Run(func(args mock.Arguments) {
			capturedArg = args.Get(0).(types.CameraMetadata)
		}).Return(nil, fmt.Errorf("update error"))
		mockAzureStorage.On("UploadImage",
			mock.AnythingOfType("*context.valueCtx"), imageID+".png", mock.AnythingOfType("[]uint8")).Return(fmt.Errorf("upload err"))
		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + Base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		if capturedArg.CameraName != expectedCamera.CameraName {
			t.Errorf("expected CameraName %s, got %s", expectedCamera.CameraName, capturedArg.CameraName)
		}
		if capturedArg.FirmwareVersion != expectedCamera.FirmwareVersion {
			t.Errorf("expected FirmwareVersion %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}
		if capturedArg.NameOfStoredPicture.String != imageID {
			t.Errorf("expected NameOfStoredPicture %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}
		if capturedArg.ImageId.String != imageID {
			t.Errorf("expected ImageId %s, got %s", expectedCamera.FirmwareVersion, capturedArg.FirmwareVersion)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withNotInitCamera_returnBadRequest", func(t *testing.T) {
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
		}

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(&expectedCamera, nil)
		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + Base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withWrongCameraID_returnNotFound", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()
		imageID := uuid.New().String()

		mockCameraStore.On("GetCameraMetadataByID", camID).Return(nil, fmt.Errorf("not found"))
		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + Base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withMalformedBase64_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()
		imageID := uuid.New().String()

		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + "123123"

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withMalformedImageID_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()
		imageID := "12"

		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + Base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withNullImageID_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		camID := uuid.New().String()
		imageID := ""

		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + Base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})

	t.Run("UploadImageHandler_withInvalidCamID_returnBadRequest", func(t *testing.T) {
		//arrange
		mockCameraStore := new(MockCameraStore)
		mockAzureStorage := new(MockAzureStorage)
		handler := NewHandler(mockCameraStore, mockAzureStorage)

		imageID := uuid.New().String()
		camID := "123"

		url := "/camera_metadata/" + camID + "/upload_image?imageID=" + imageID + "&image_as_bytes=" + Base64Data

		// Act
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/camera_metadata/{camID}/upload_image", handler.UploadImageHandler).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		mockCameraStore.AssertExpectations(t)
	})
}
