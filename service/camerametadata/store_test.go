package camerametadata

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go-sample-rest-api/types"
	"testing"
	"time"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, mock, cleanup
}

func TestStore_CreateCameraMetadata(t *testing.T) {
	t.Run("CreateCameraMetaData_withValidData_toCreateCameraMetadata", func(t *testing.T) {
		// arrange
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		store := Store{db}

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camera := types.CameraMetadata{
			CameraName:      "Test Camera",
			FirmwareVersion: "v1.0",
			CreatedAt:       nullTime,
		}

		expectedID := uuid.New().String()

		mock.ExpectQuery(`INSERT INTO camera_metadata`).
			WithArgs(camera.CameraName, camera.FirmwareVersion, camera.CreatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"cam_id", "camera_name", "firmware_version", "created_at"}).
				AddRow(expectedID, camera.CameraName, camera.FirmwareVersion, camera.CreatedAt))

		// act
		savedCamera, err := store.CreateCameraMetadata(camera)

		// assert
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if err != nil {
			t.Errorf("expected no error, but got %s", err)
		}
		if savedCamera.CamID != expectedID {
			t.Errorf("expected returned camera ID to be %s, got %s", expectedID, savedCamera.CamID)
		}
	})
	t.Run("CreateCameraMetaData_withError_toReturnError", func(t *testing.T) {
		// arrange
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		store := Store{db}

		timeNow := time.Now()
		nullTime := sql.NullTime{
			Time:  timeNow,
			Valid: true,
		}
		camera := types.CameraMetadata{
			CameraName:      "Test Camera",
			FirmwareVersion: "v1.0",
			CreatedAt:       nullTime,
		}

		mock.ExpectQuery(`INSERT INTO camera_metadata`).
			WithArgs(camera.CameraName, camera.FirmwareVersion, camera.CreatedAt).
			WillReturnError(sql.ErrConnDone)

		// act
		_, err := store.CreateCameraMetadata(camera)

		// assert
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if err == nil {
			t.Errorf("expected an error, got none")
		}

		if err != sql.ErrConnDone {
			t.Errorf("expected a sql.ErrConnDone error, got %v", err)
		}
	})
}

func TestStore_GetCameraMetadataByID(t *testing.T) {
	t.Run("GetCameraMetadataByID_withValidData_toGetCameraMetadata", func(t *testing.T) {
		// arrange
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		store := Store{db}

		camID := uuid.New().String()

		rows := sqlmock.NewRows([]string{"cam_id", "image_id", "camera_name", "firmware_version", "container_name", "name_of_stored_picture", "created_at", "onboarded_at", "initialized_at"}).
			AddRow(camID, nil, "Test Camera", "v1.0", nil, nil, time.Now(), time.Now(), time.Now())
		mock.ExpectQuery(`^SELECT cam_id, image_id, camera_name, firmware_version, container_name, name_of_stored_picture, created_at, onboarded_at, initialized_at FROM camera_metadata WHERE cam_id = \$1$`).
			WithArgs(camID).
			WillReturnRows(rows)

		// act
		cameraMetadata, err := store.GetCameraMetadataByID(camID)

		// assert
		assert.NoError(t, mock.ExpectationsWereMet())

		assert.NoError(t, err)
		assert.NotNil(t, cameraMetadata)
		assert.Equal(t, camID, cameraMetadata.CamID)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if err != nil {
			t.Errorf("expected no error, but got %s", err)
		}
	})
	t.Run("GetCameraMetadataByID_withNoRows_toReturnError", func(t *testing.T) {
		// arrange
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()
		store := Store{db}

		camID := "non-existent-id"
		mock.ExpectQuery(`^SELECT.*FROM camera_metadata WHERE cam_id = \$1$`).
			WithArgs(camID).
			WillReturnError(sql.ErrNoRows)

		// act
		cameraMetadata, err := store.GetCameraMetadataByID(camID)

		// assert
		assert.Error(t, err)
		assert.Nil(t, cameraMetadata)
		assert.Equal(t, "no camera metadata found with ID: non-existent-id", err.Error())

	})
	t.Run("CreateCameraMetaData_withError_toReturnError", func(t *testing.T) {
		// arrange
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		store := Store{db}

		camID := "non-existent-id"
		mock.ExpectQuery(`^SELECT.*FROM camera_metadata WHERE cam_id = \$1$`).
			WithArgs(camID).
			WillReturnError(sql.ErrConnDone)

		// act
		_, err := store.GetCameraMetadataByID(camID)

		// assert
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if err == nil {
			t.Errorf("expected an error, got none")
		}

		if err != sql.ErrConnDone {
			t.Errorf("expected a sql.ErrConnDone error, got %v", err)
		}
	})
}
func TestStore_UpdateCameraMetadata(t *testing.T) {
	t.Run("UpdateCameraMetadata_withValidData_toGetCameraMetadata", func(t *testing.T) {
		// arrange
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		store := Store{db}

		cam := types.CameraMetadata{
			CamID:               "123",
			CameraName:          "Test Camera",
			FirmwareVersion:     "v1.0",
			ContainerName:       sql.NullString{String: "Test Container", Valid: true},
			NameOfStoredPicture: sql.NullString{String: "test_picture.jpg", Valid: true},
			CreatedAt:           sql.NullTime{Time: time.Now(), Valid: true},
			OnboardedAt:         sql.NullTime{Time: time.Now(), Valid: true},
			InitializedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		}

		mock.ExpectExec("UPDATE camera_metadata").WithArgs(
			cam.CameraName, cam.FirmwareVersion, cam.ContainerName, cam.NameOfStoredPicture,
			cam.CreatedAt, cam.OnboardedAt, cam.InitializedAt, cam.CamID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		updatedCamera, err := store.UpdateCameraMetadata(cam)

		// assert
		assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")

		assert.NoError(t, err, "Error was not expected when updating camera metadata")
		assert.NotNil(t, updatedCamera, "Updated camera metadata should not be nil")
		assert.Equal(t, cam.CamID, updatedCamera.CamID, "Updated camera metadata should have the same CamID")
	})
	t.Run("UpdateCameraMetadata_withNoRows_toReturnError", func(t *testing.T) {
		// arrange
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()
		store := Store{db}

		cam := types.CameraMetadata{
			CamID:               "123",
			CameraName:          "Test Camera",
			FirmwareVersion:     "v1.0",
			ContainerName:       sql.NullString{String: "Test Container", Valid: true},
			NameOfStoredPicture: sql.NullString{String: "test_picture.jpg", Valid: true},
			CreatedAt:           sql.NullTime{Time: time.Now(), Valid: true},
			OnboardedAt:         sql.NullTime{Time: time.Now(), Valid: true},
			InitializedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		}

		mock.ExpectExec("UPDATE camera_metadata").WithArgs(
			cam.CameraName, cam.FirmwareVersion, cam.ContainerName, cam.NameOfStoredPicture,
			cam.CreatedAt, cam.OnboardedAt, cam.InitializedAt, cam.CamID,
		).WillReturnError(sql.ErrNoRows)

		// act
		updatedCamera, err := store.UpdateCameraMetadata(cam)

		// assert
		assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")

		assert.Error(t, err, "Expected an error when no rows are affected")
		assert.Nil(t, updatedCamera, "No camera metadata should be returned when no rows are affected")
	})
	t.Run("UpdateCameraMetadata_withError_toReturnError", func(t *testing.T) {
		// arrange
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		store := Store{db}

		cam := types.CameraMetadata{
			CamID:               "123",
			CameraName:          "Test Camera",
			FirmwareVersion:     "v1.0",
			ContainerName:       sql.NullString{String: "Test Container", Valid: true},
			NameOfStoredPicture: sql.NullString{String: "test_picture.jpg", Valid: true},
			CreatedAt:           sql.NullTime{Time: time.Now(), Valid: true},
			OnboardedAt:         sql.NullTime{Time: time.Now(), Valid: true},
			InitializedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		}

		mock.ExpectExec("UPDATE camera_metadata").WithArgs(
			cam.CameraName, cam.FirmwareVersion, cam.ContainerName, cam.NameOfStoredPicture,
			cam.CreatedAt, cam.OnboardedAt, cam.InitializedAt, cam.CamID,
		).WillReturnError(sql.ErrConnDone)

		// act
		_, err := store.UpdateCameraMetadata(cam)

		// assert
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if err == nil {
			t.Errorf("expected an error, got none")
		}

		if err != sql.ErrConnDone {
			t.Errorf("expected a sql.ErrConnDone error, got %v", err)
		}
	})
}
