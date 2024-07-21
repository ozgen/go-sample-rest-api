package types

import (
	"database/sql"
)

type CameraMetadata struct {
	CamID               string         `json:"cam_id"`
	ImageId             sql.NullString `json:"image_id"`
	CameraName          string         `json:"camera_name"`
	FirmwareVersion     string         `json:"firmware_version"`
	ContainerName       sql.NullString `json:"container_name"`
	NameOfStoredPicture sql.NullString `json:"name_of_stored_picture"`
	CreatedAt           sql.NullTime   `json:"createdAt"`
	OnboardedAt         sql.NullTime   `json:"onboarded_at"`
	InitializedAt       sql.NullTime   `json:"initialized_at"`
}

type CameraMetadataPayload struct {
	CameraName      string `json:"camera_name" validate:"required"`
	FirmwareVersion string `json:"firmware_version" validate:"required"`
}

type CameraMetadataResponse struct {
	CamID           string       `json:"cam_id"`
	CameraName      string       `json:"camera_name"`
	FirmwareVersion string       `json:"firmware_version"`
	CreatedAt       sql.NullTime `json:"createdAt"`
}
type ImageUploadedResponse struct {
	CamID           string         `json:"cam_id"`
	CameraName      string         `json:"camera_name"`
	FirmwareVersion string         `json:"firmware_version"`
	ImageId         sql.NullString `json:"image_id"`
}

type CameraMetadataStore interface {
	CreateCameraMetadata(camera CameraMetadata) (*CameraMetadata, error)
	GetCameraMetadataByID(camID string) (*CameraMetadata, error)
	UpdateCameraMetadata(camera CameraMetadata) (*CameraMetadata, error)
}
