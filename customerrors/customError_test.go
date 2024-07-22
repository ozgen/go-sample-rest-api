package customerrors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNotFoundError(t *testing.T) {
	err := &NotFoundError{ID: "123"}
	expectedMessage := "item with ID 123 not found"
	assert.Equal(t, expectedMessage, err.Error(), "Error message should match expected output")
}

func TestNotInitError(t *testing.T) {
	err := &NotInitError{ID: "456"}
	expectedMessage := "camera with ID 456 not initialized"
	assert.Equal(t, expectedMessage, err.Error(), "Error message should match expected output")
}

func TestAlreadyInitError(t *testing.T) {
	err := &AlreadyInitError{ID: "789"}
	expectedMessage := "camera with ID 789 is already initialized"
	assert.Equal(t, expectedMessage, err.Error(), "Error message should match expected output")
}

func TestAzureStorageError(t *testing.T) {
	err := &AzureStorageError{Message: "test message"}
	expectedMessage := "Azure blob storage err: test message"
	assert.Equal(t, expectedMessage, err.Error(), "Error message should match expected output")
}
