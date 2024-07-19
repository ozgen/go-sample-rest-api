package camerametadata

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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

func TestStore_CreateCameraMetaDataCreateCameraMetaData(t *testing.T) {
	t.Run("CreateCameraMetaData_withValidData_toCreateCameraMetadata", func(t *testing.T) {
		// arrange
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		store := Store{db}

		camera := types.CameraMetadata{
			CameraName:      "Test Camera",
			FirmwareVersion: "v1.0",
			CreatedAt:       time.Now(),
		}

		expectedID := uuid.New().String()

		mock.ExpectQuery(`INSERT INTO camera_metadata`).
			WithArgs(camera.CameraName, camera.FirmwareVersion, camera.CreatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"cam_id", "camera_name", "firmware_version", "created_at"}).
				AddRow(expectedID, camera.CameraName, camera.FirmwareVersion, camera.CreatedAt))

		// act
		savedCamera, err := store.CreateCameraMetaData(camera)

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

		camera := types.CameraMetadata{
			CameraName:      "Test Camera",
			FirmwareVersion: "v1.0",
			CreatedAt:       time.Now(),
		}

		mock.ExpectQuery(`INSERT INTO camera_metadata`).
			WithArgs(camera.CameraName, camera.FirmwareVersion, camera.CreatedAt).
			WillReturnError(sql.ErrConnDone)

		// act
		_, err := store.CreateCameraMetaData(camera)

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
