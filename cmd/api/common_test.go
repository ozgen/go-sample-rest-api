package api

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil // Implement as needed for the tests
}

func (m *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil // Implement as needed for the tests
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil // Implement as needed for the tests
}

type MockAzureStorage struct {
	mock.Mock
}

func (m *MockAzureStorage) UploadImage(ctx context.Context, blobName string, imageData []byte) error {
	args := m.Called(ctx, blobName, imageData)
	return args.Error(0)
}

func (m *MockAzureStorage) DownloadImage(ctx context.Context, blobName string) ([]byte, error) {
	args := m.Called(ctx, blobName)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
