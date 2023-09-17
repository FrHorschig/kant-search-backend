package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertErrorResponse(t *testing.T, res *httptest.ResponseRecorder) {
	assert.Contains(t, res.Body.String(), "code")
	assert.Contains(t, res.Body.String(), "message")
}
