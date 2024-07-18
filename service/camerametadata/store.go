package camerametadata

import (
	"database/sql"
	"go-sample-rest-api/types"
	"log"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateCameraMetaData(camera types.CameraMetadata) (*types.CameraMetadata, error) {
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

func scanRowsIntoUser(rows *sql.Rows) (*types.CameraMetadata, error) {
	cameraMetadata := new(types.CameraMetadata)

	err := rows.Scan(
		&cameraMetadata.CamID,
		&cameraMetadata.ImageId,
		&cameraMetadata.CameraName,
		&cameraMetadata.FirmwareVersion,
		&cameraMetadata.ContainerName,
		&cameraMetadata.NameOfStoredPicture,
		&cameraMetadata.CreatedAt,
		&cameraMetadata.OnboardedAt,
		&cameraMetadata.InitializedAt,
	)
	if err != nil {
		return nil, err
	}

	return cameraMetadata, nil
}
