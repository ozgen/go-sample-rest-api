package types

import "time"

type CameraMetadata struct {
	CamID               string    `json:"cam_id"`
	ImageId             string    `json:"image_id"`
	CameraName          string    `json:"camera_name"`
	FirmwareVersion     string    `json:"firmware_version"`
	ContainerName       string    `json:"container_name"`
	NameOfStoredPicture string    `json:"name_of_stored_picture"`
	CreatedAt           time.Time `json:"createdAt"`
	OnboardedAt         time.Time `json:"onboarded_at"`
	InitializedAt       time.Time `json:"initialized_at"`
}

type CameraMetadataPayload struct {
	CameraName      string `json:"camera_name" validate:"required"`
	FirmwareVersion string `json:"firmware_version" validate:"required"`
}

type CameraMetadataStore interface {
	CreateCameraMetaData(camera CameraMetadata) (*CameraMetadata, error)
}
