package camerametadata

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go-sample-rest-api/config"
	"go-sample-rest-api/customerrors"
	"go-sample-rest-api/logging"
	"go-sample-rest-api/storage"
	"go-sample-rest-api/types"
	"go-sample-rest-api/utils"
	"net/http"
	"time"
)

type Handler struct {
	store        types.CameraMetadataStore
	azureStorage storage.ImageStore
}

func NewHandler(store types.CameraMetadataStore, azureStorage storage.ImageStore) *Handler {
	return &Handler{store: store, azureStorage: azureStorage}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/camera_metadata", h.CreateCameraMetadata).Methods(http.MethodPost)
	router.HandleFunc("/camera_metadata/{camID}/init", h.InitializeCameraMetaData).Methods(http.MethodPatch)
	router.HandleFunc("/camera_metadata/{camID}", h.GetCameraMetaData).Methods(http.MethodGet)
	router.HandleFunc("/camera_metadata/{camID}/upload_image", h.UploadImageHandler).Methods(http.MethodPost)
	router.HandleFunc("/camera_metadata/{camID}/download_image", h.DownloadImageHandler).Methods(http.MethodGet)
}

// CreateCameraMetadata godoc
// @Summary Create camera metadata
// @Description Creates a new camera metadata entry.
// @Tags camera
// @Accept json
// @Produce json
// @Param cameraMetadata body types.CameraMetadataPayload true "Camera Metadata Info"
// @Success 201 {object} types.CameraMetadataResponse "Camera metadata successfully created."
// @Failure 400 {object} types.HTTPError "Invalid request parameters."
// @Failure 500 {object} types.HTTPError "Internal server error."
// @Router /camera_metadata [post]
func (h *Handler) CreateCameraMetadata(writer http.ResponseWriter, request *http.Request) {
	var cameraMetadata types.CameraMetadataPayload
	log := logging.GetLogger()
	if err := utils.ParseJSON(request, &cameraMetadata); err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		log.WithFields(logrus.Fields{
			"error": err,
			"body":  request.Body,
		}).Error("Malformed cameraMetadata request")
		return
	}

	if err := utils.Validate.Struct(cameraMetadata); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		log.WithFields(logrus.Fields{
			"validationErrors": errors,
		}).Error("Validation failed for cameraMetadata request")
		return
	}
	timeNow := time.Now()
	nullTime := sql.NullTime{
		Time:  timeNow,
		Valid: true,
	}
	savedCamera, err := h.store.CreateCameraMetadata(types.CameraMetadata{
		CameraName:      cameraMetadata.CameraName,
		FirmwareVersion: cameraMetadata.FirmwareVersion,
		CreatedAt:       nullTime,
	})

	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to create camera metadata")
		return
	}
	cameraResponse := types.CameraMetadataResponse{
		CamID:           savedCamera.CamID,
		CameraName:      savedCamera.CameraName,
		FirmwareVersion: savedCamera.FirmwareVersion,
		CreatedAt:       savedCamera.CreatedAt.Time,
	}
	utils.WriteJSON(writer, http.StatusCreated, cameraResponse)
}

// InitializeCameraMetaData godoc
// @Summary Initialize camera metadata
// @Description Marks a camera metadata entry as initialized.
// @Tags camera
// @Accept json
// @Produce json
// @Param camID path string true "Camera ID"
// @Success 200 {object} nil "Camera metadata initialized successfully."
// @Failure 400 {object} types.HTTPError "Invalid camera ID."
// @Failure 404 {object} types.HTTPError "Camera not found."
// @Failure 409 {object} types.HTTPError "Camera already initialized."
// @Failure 500 {object} types.HTTPError "Internal server error."
// @Router /camera_metadata/{camID}/init [patch]
func (h *Handler) InitializeCameraMetaData(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	camID := vars["camID"]

	_, err := uuid.Parse(camID)
	if err != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid camID: %v", err))
		return
	}

	cameraMetadata, err := h.store.GetCameraMetadataByID(camID)
	if err != nil {
		utils.WriteError(writer, http.StatusNotFound, &customerrors.NotFoundError{ID: camID})
		return
	}

	if cameraMetadata.InitializedAt.Valid {
		utils.WriteError(writer, http.StatusConflict, &customerrors.AlreadyInitError{ID: camID})
		return
	}

	timeNow := time.Now()
	nullTime := sql.NullTime{
		Time:  timeNow,
		Valid: true,
	}
	cameraMetadata.InitializedAt = nullTime

	_, err = h.store.UpdateCameraMetadata(*cameraMetadata)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, fmt.Errorf("failed to update camera metadata: %v", err))
		return
	}

	utils.WriteJSON(writer, http.StatusOK, nil)
}

// GetCameraMetaData godoc
// @Summary Get camera metadata
// @Description Retrieves metadata for a specific camera.
// @Tags camera
// @Accept json
// @Produce json
// @Param camID path string true "Camera ID"
// @Success 200 {object} types.CameraMetadataResponse "Camera metadata found."
// @Failure 400 {object} types.HTTPError "Invalid camera ID."
// @Failure 404 {object} types.HTTPError "Camera metadata not found."
// @Router /camera_metadata/{camID} [get]
func (h *Handler) GetCameraMetaData(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	camID := vars["camID"]

	_, err := uuid.Parse(camID)
	if err != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid camID: %v", err))
		return
	}

	cameraMetadata, err := h.store.GetCameraMetadataByID(camID)
	if err != nil {
		utils.WriteError(writer, http.StatusNotFound, fmt.Errorf("camera given ID: %v not found", camID))
		return
	}

	cameraResponse := types.CameraMetadataResponse{
		CamID:           cameraMetadata.CamID,
		CameraName:      cameraMetadata.CameraName,
		FirmwareVersion: cameraMetadata.FirmwareVersion,
		CreatedAt:       cameraMetadata.CreatedAt.Time,
	}

	utils.WriteJSON(writer, http.StatusOK, cameraResponse)
}

// UploadImageHandler godoc
// @Summary Upload an image to a camera
// @Description Uploads an image file associated with a camera.
// @Tags camera
// @Accept multipart/form-data
// @Produce json
// @Param camID path string true "Camera ID"
// @Param imageID query string true "Image ID"
// @Param image_as_bytes body string true "Base64 encoded image data"
// @Success 200 {object} types.ImageUploadedResponse "Image uploaded successfully."
// @Failure 400 {object} types.HTTPError "Bad request parameters."
// @Failure 404 {object} types.HTTPError "Camera metadata not found."
// @Failure 500 {object} types.HTTPError "Failed to upload image."
// @Router /camera_metadata/{camID}/upload_image [post]
func (h *Handler) UploadImageHandler(writer http.ResponseWriter, request *http.Request) {
	log := logging.GetLogger()
	vars := mux.Vars(request)
	camID := vars["camID"]

	_, err := uuid.Parse(camID)
	if err != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid camID: %v", err))
		return
	}

	// Extracting query parameters
	query := request.URL.Query()
	imageID := query.Get("imageID")
	imageAsBytes := utils.NormalizeBase64(query.Get("image_as_bytes"))

	if imageID == "" || imageAsBytes == "" {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("Missing required query parameters"))
		return
	}

	_, err = uuid.Parse(imageID)
	if err != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid imageID: %v", err))
		return
	}

	// Assuming imageAsBytes is base64 encoded data
	imageData, err := base64.StdEncoding.DecodeString(imageAsBytes)
	if err != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("failed to decode image data: %v", err))
		return
	}

	cameraMetadata, err := h.store.GetCameraMetadataByID(camID)
	if err != nil {
		utils.WriteError(writer, http.StatusNotFound, &customerrors.NotFoundError{ID: camID})
		return
	}

	if !cameraMetadata.InitializedAt.Valid {
		utils.WriteError(writer, http.StatusBadRequest, &customerrors.NotInitError{ID: camID})
		return
	}

	cameraMetadata.ImageId = sql.NullString{String: imageID, Valid: true}
	cameraMetadata.NameOfStoredPicture = sql.NullString{String: imageID, Valid: true}
	cameraMetadata.ContainerName = sql.NullString{String: config.Envs.AzureContainerName, Valid: true}

	_, err = h.store.UpdateCameraMetadata(*cameraMetadata)
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, fmt.Errorf("failed to update camera metadata: %v", err))
		return
	}

	err = h.azureStorage.UploadImage(request.Context(), imageID+".png", imageData)

	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, fmt.Errorf("failed to upload image: %v", err))
		return
	}

	log.WithFields(logrus.Fields{
		"camera": cameraMetadata,
	}).Info("Image uploaded successfully")

	response := types.ImageUploadedResponse{
		CamID:           cameraMetadata.CamID,
		CameraName:      cameraMetadata.CameraName,
		FirmwareVersion: cameraMetadata.FirmwareVersion,
		ImageId:         cameraMetadata.ImageId.String,
	}
	utils.WriteJSON(writer, http.StatusOK, response)
}

// DownloadImageHandler godoc
// @Summary Download an image from a camera
// @Description Downloads an image file associated with a camera.
// @Tags camera
// @Produce octet-stream
// @Param camID path string true "Camera ID"
// @Success 200 {file} file "Image file downloaded successfully."
// @Failure 400 {object} types.HTTPError "Invalid camera ID."
// @Failure 404 {object} types.HTTPError "Image not found."
// @Failure 500 {object} types.HTTPError "Failed to download image."
// @Router /camera_metadata/{camID}/download_image [get]
func (h *Handler) DownloadImageHandler(writer http.ResponseWriter, request *http.Request) {
	log := logging.GetLogger()
	vars := mux.Vars(request)
	camID := vars["camID"]

	_, err := uuid.Parse(camID)
	if err != nil {
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid camID: %v", err))
		return
	}

	cameraMetadata, err := h.store.GetCameraMetadataByID(camID)
	if err != nil {
		utils.WriteError(writer, http.StatusNotFound, &customerrors.NotFoundError{ID: camID})
		return
	}

	if !cameraMetadata.ImageId.Valid {
		utils.WriteError(writer, http.StatusNotFound, &customerrors.NotFoundError{ID: "ImageId"})
		return
	}

	image, err := h.azureStorage.DownloadImage(request.Context(), cameraMetadata.ImageId.String+".png")
	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, fmt.Errorf("failed to download image: %v", err))
		return
	}

	writer.Header().Set("Content-Type", "image/png") // Assuming the image is in PNG format
	if _, err := writer.Write(image); err != nil {
		log.WithFields(logrus.Fields{
			"camID": camID,
			"error": err,
		}).Error("Failed to write image to response")
		utils.WriteError(writer, http.StatusInternalServerError, fmt.Errorf("failed to write image to response: %v", err))
		return
	}
	writer.WriteHeader(http.StatusOK)
	log.Infof("Successfully sent image for camera ID: %s", camID)
}
