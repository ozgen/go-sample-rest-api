package camerametadata

import (
	"database/sql"
	"fmt"
	"go-sample-rest-api/types"
	"log"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateCameraMetadata(camera types.CameraMetadata) (*types.CameraMetadata, error) {
	query := `INSERT INTO camera_metadata (
        camera_name, firmware_version, created_at) VALUES ($1, $2, $3) 
        RETURNING cam_id, camera_name, firmware_version, created_at`

	var savedCamera types.CameraMetadata

	err := s.db.QueryRow(query, camera.CameraName, camera.FirmwareVersion, camera.CreatedAt).
		Scan(&savedCamera.CamID, &savedCamera.CameraName, &savedCamera.FirmwareVersion, &savedCamera.CreatedAt)
	if err != nil {
		return nil, err
	}
	log.Printf("camera metada saved successfully, cameraMetada: %v", savedCamera)
	return &savedCamera, nil
}

func (s *Store) UpdateCameraMetadata(camera types.CameraMetadata) (*types.CameraMetadata, error) {
	query := `
        UPDATE camera_metadata
        SET camera_name = $1, 
            firmware_version = $2, 
            container_name = $3, 
            name_of_stored_picture = $4, 
            created_at = $5, 
            onboarded_at = $6, 
            initialized_at = $7
        WHERE cam_id = $8;
    `

	_, err := s.db.Exec(query, camera.CameraName, camera.FirmwareVersion, camera.ContainerName, camera.NameOfStoredPicture, camera.CreatedAt, camera.OnboardedAt, camera.InitializedAt, camera.CamID)
	if err != nil {
		return nil, err
	}

	return &camera, nil
}

func (s *Store) GetCameraMetadataByID(camID string) (*types.CameraMetadata, error) {
	query := `SELECT cam_id, image_id, camera_name, firmware_version, container_name,
              name_of_stored_picture, created_at, onboarded_at, initialized_at 
              FROM camera_metadata WHERE cam_id = $1`

	row := s.db.QueryRow(query, camID)

	c := new(types.CameraMetadata)

	err := row.Scan(&c.CamID, &c.ImageId, &c.CameraName, &c.FirmwareVersion, &c.ContainerName,
		&c.NameOfStoredPicture, &c.CreatedAt, &c.OnboardedAt, &c.InitializedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no camera metadata found with ID: %s", camID)
		}
		return nil, err
	}

	return c, nil
}
