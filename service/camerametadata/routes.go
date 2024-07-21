package camerametadata

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go-sample-rest-api/customerrors"
	"go-sample-rest-api/logging"
	"go-sample-rest-api/types"
	"go-sample-rest-api/utils"
	"net/http"
	"time"
)

type Handler struct {
	store types.CameraMetadataStore
}

func NewHandler(store types.CameraMetadataStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/camera_metadata", h.CreateCameraMetadata).Methods(http.MethodPost)
	router.HandleFunc("/camera_metadata/{camID}/init", h.InitializeCameraMetaData).Methods(http.MethodPatch)
	router.HandleFunc("/camera_metadata/{camID}", h.GetCameraMetaData).Methods(http.MethodGet)
}

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
		CreatedAt:       savedCamera.CreatedAt,
	}
	utils.WriteJSON(writer, http.StatusCreated, cameraResponse)
}

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
		CreatedAt:       cameraMetadata.CreatedAt,
	}

	utils.WriteJSON(writer, http.StatusOK, cameraResponse)
}
