package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupServer(t *testing.T) {
	server, err := SetupServer()
	assert.NoError(t, err)
	assert.NotNil(t, server, "The server should be initialized")
}
