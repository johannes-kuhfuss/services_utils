package api_error

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	msg := "new error test"
	err := NewError(msg, 55, nil)
	assert.NotNil(t, err)
	assert.EqualValues(t, msg, err.Message())
	assert.EqualValues(t, 55, err.StatusCode())
	assert.Nil(t, err.Causes())
}

func TestError(t *testing.T) {
	err := NewError("error test", 54, nil)
	errString := err.Error()
	assert.EqualValues(t, errString, "Message: error test - Status: 54 - Causes: []")
}

func TestNewErrorFromBytesError(t *testing.T) {
	var bytes []byte
	restErr, err := NewErrorFromBytes(bytes)
	assert.NotNil(t, err)
	assert.Nil(t, restErr)
	assert.EqualValues(t, "invalid json", err.Error())
}

func TestNewErrorFromBytesNoError(t *testing.T) {
	bytes := []byte("{\"message\":\"bytesTest\",\"statuscode\":56,\"causes\":null}")
	restErr, err := NewErrorFromBytes(bytes)
	assert.NotNil(t, restErr)
	assert.Nil(t, err)
	assert.EqualValues(t, "bytesTest", restErr.Message())
	assert.EqualValues(t, 56, restErr.StatusCode())
	assert.Nil(t, restErr.Causes())
}

func TestNewBadRequestError(t *testing.T) {
	err := NewBadRequestError("new bad request error test")
	assert.NotNil(t, err)
	assert.EqualValues(t, "new bad request error test", err.Message())
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.Nil(t, err.Causes())
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("new not found error test")
	assert.NotNil(t, err)
	assert.EqualValues(t, "new not found error test", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.Nil(t, err.Causes())
}

func TestNewInternalServerErrorNoExtraError(t *testing.T) {
	msg := "new internal server error test, no extra error"
	err := NewInternalServerError(msg, nil)
	assert.NotNil(t, err)
	assert.EqualValues(t, msg, err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.Nil(t, err.Causes())
}

func TestNewUnauthorizedError(t *testing.T) {
	err := NewUnauthorizedError("new unauthorized server error test - user cannot access resource - 403 permanent")
	assert.NotNil(t, err)
	assert.EqualValues(t, "new unauthorized server error test - user cannot access resource - 403 permanent", err.Message())
	assert.EqualValues(t, http.StatusForbidden, err.StatusCode())
	assert.Nil(t, err.Causes())
}

func TestNewUnauthenticatedError(t *testing.T) {
	err := NewUnauthenticatedError("new unauthenticated error test - wrong credentials - 401 try again")
	assert.NotNil(t, err)
	assert.EqualValues(t, "new unauthenticated error test - wrong credentials - 401 try again", err.Message())
	assert.EqualValues(t, http.StatusUnauthorized, err.StatusCode())
	assert.Nil(t, err.Causes())
}

func TestNewInternalServerErrorWithExtraError(t *testing.T) {
	msg := "new internal server error test"
	err := NewInternalServerError(msg, errors.New("test error"))
	assert.NotNil(t, err)
	assert.EqualValues(t, msg, err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.NotNil(t, err.Causes())
	assert.EqualValues(t, 1, len(err.Causes()))
	assert.EqualValues(t, "test error", err.Causes()[0])
}

func TestNewProcessingConflictError(t *testing.T) {
	err := NewProcessingConflictError("new processing conflict server error test")
	assert.NotNil(t, err)
	assert.EqualValues(t, "new processing conflict server error test", err.Message())
	assert.EqualValues(t, http.StatusConflict, err.StatusCode())
	assert.Nil(t, err.Causes())
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("new validation error test")
	assert.NotNil(t, err)
	assert.EqualValues(t, "new validation error test", err.Message())
	assert.EqualValues(t, http.StatusUnprocessableEntity, err.StatusCode())
	assert.Nil(t, err.Causes())
}
