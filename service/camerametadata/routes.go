package camerametadata

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go-sample-rest-api/types"
	"go-sample-rest-api/utils"
	"log"
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
	router.HandleFunc("/camera_metadata", h.CreateCameraMetaData).Methods(http.MethodPost)
}

func (h *Handler) CreateCameraMetaData(writer http.ResponseWriter, request *http.Request) {
	var cameraMetadata types.CameraMetadataPayload
	if err := utils.ParseJSON(request, &cameraMetadata); err != nil {
		utils.WriteError(writer, http.StatusBadRequest, err)
		log.Printf("Malformed cameraMetadata request: %v", request.Body)
		return
	}

	if err := utils.Validate.Struct(cameraMetadata); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		log.Printf("Malformed cameraMetadata request: %v", errors)
		return
	}

	savedCamera, err := h.store.CreateCameraMetaData(types.CameraMetadata{
		CameraName:      cameraMetadata.CameraName,
		FirmwareVersion: cameraMetadata.FirmwareVersion,
		CreatedAt:       time.Now(),
	})

	if err != nil {
		utils.WriteError(writer, http.StatusInternalServerError, err)
		log.Printf("Malformed cameraMetadata request: %v", err)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, savedCamera)
}
