package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XRay-Addons/xrayman/node/internal/http/handlers/mocks"
	"github.com/XRay-Addons/xrayman/node/internal/logging"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/XRay-Addons/xrayman/node/pkg/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)

	testRequest := api.StartRequest{
		Users: []api.User{
			{
				Name:      "user1",
				VlessUUID: "vless_uuid1",
			},
			{
				Name:      "user2",
				VlessUUID: "vless_uuid2",
			},
		},
	}

	expectedRequest := models.StartParams{
		Users: []models.User{
			{
				Name:      "user1",
				VlessUUID: "vless_uuid1",
			},
			{
				Name:      "user2",
				VlessUUID: "vless_uuid2",
			},
		},
	}

	testResult := models.StartResult{
		ClientCfg: models.ClientCfg{
			Template:       `{ "Config": "Template" }`,
			UserNameField:  "UserName",
			VlessUUIDField: "VlessUUID",
		},
	}

	expectedResponse := api.StartResponse{
		ClientCfg: api.ClientCfg{
			Template:       `{ "Config": "Template" }`,
			UserNameField:  "UserName",
			VlessUUIDField: "VlessUUID",
		},
	}

	mockService.EXPECT().
		Start(gomock.Any(), expectedRequest).
		Return(&testResult, nil).
		Times(1)

	handler, err := New(mockService)
	require.NoError(t, err)

	log, err := logging.New()
	require.NoError(t, err)

	requestBody, err := json.Marshal(testRequest)
	require.NoError(t, err)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Start(log).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response api.StartResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)

	expectedRequest := models.StatusParams{}

	testResult := models.StatusResult{
		ServiceStatus: models.ServiceRunning,
	}

	expectedResponse := api.StatusResponse{
		ServiceStatus: api.ServiceRunning,
	}

	mockService.EXPECT().
		Status(gomock.Any(), expectedRequest).
		Return(&testResult, nil).
		Times(1)

	handler, err := New(mockService)
	require.NoError(t, err)

	log, err := logging.New()
	require.NoError(t, err)

	require.NoError(t, err)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.Status(log).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response api.StatusResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}
