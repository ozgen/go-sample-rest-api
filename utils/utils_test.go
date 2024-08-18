package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	const envKey = "TEST_ENV_VAR"
	const defaultValue = "default"

	os.Setenv(envKey, "exists")
	assert.Equal(t, "exists", GetEnv(envKey, defaultValue))

	os.Unsetenv(envKey)
	assert.Equal(t, defaultValue, GetEnv(envKey, defaultValue))
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"hello": "world"}

	err := WriteJSON(w, http.StatusOK, data)
	assert.Nil(t, err)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "world", response["hello"])
}

func TestParseJSON(t *testing.T) {
	data := `{"key": "value"}`
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(data))

	var result map[string]string
	err := ParseJSON(r, &result)
	assert.Nil(t, err)
	assert.Equal(t, "value", result["key"])

	r = httptest.NewRequest(http.MethodPost, "/", nil)
	err = ParseJSON(r, &result)
	assert.NotNil(t, err)
}

func TestParseJSONResponse(t *testing.T) {
	data := `{"key": "value"}`
	// Create a new HTTP response with the JSON data
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBufferString(data)

	// Simulate the response by setting the content type and other necessary headers
	rr.Header().Set("Content-Type", "application/json")

	var result map[string]string
	// Use the Response from the httptest.Recorder
	err := ParseJSONResponse(rr.Result(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "value", result["key"])

	// Test with an empty body
	rrEmpty := httptest.NewRecorder()
	errEmpty := ParseJSONResponse(rrEmpty.Result(), &result)
	assert.NotNil(t, errEmpty)
}

func TestGetTokenFromRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?token=queryToken", nil)
	r.Header.Add("Authorization", "headerToken")

	assert.Equal(t, "headerToken", GetTokenFromRequest(r))

	// Test without the Authorization header
	r.Header.Del("Authorization")
	assert.Equal(t, "queryToken", GetTokenFromRequest(r))
}

func TestNormalizeBase64(t *testing.T) {
	input := "a\\// b"
	expected := "a//+b"
	assert.Equal(t, expected, NormalizeBase64(input))
}
