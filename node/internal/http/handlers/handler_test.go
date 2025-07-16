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
	apimodels "github.com/XRay-Addons/xrayman/node/pkg/api/models"
	apirequests "github.com/XRay-Addons/xrayman/node/pkg/api/requests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)

	testRequest := apirequests.StartRequest{
		Users: []apimodels.User{
			{
				ID:        1,
				Name:      "user1",
				VlessUUID: "vless_uuid1",
			},
			{
				ID:        2,
				Name:      "user2",
				VlessUUID: "vless_uuid2",
			},
		},
	}

	expectedRequest := []models.User{
		{"user1", "vless_uuid1"},
		{"user2", "vless_uuid2"},
	}

	testNodeProps := models.NodeProperties{
		ClientCfgTemplate: models.ClientCfgTemplate{
			Template:       `{ "Config": "Template" }`,
			UserNameField:  "UserName",
			VlessUUIDField: "VlessUUID",
		},
	}

	expectedResponse := apirequests.StartResponse{
		NodeConfig: apimodels.NodeConfig{
			UserConfigTemplate: apimodels.ClientCfgTemplate{
				Template:       `{ "Config": "Template" }`,
				UserNameField:  "UserName",
				VlessUUIDField: "VlessUUID",
			},
		},
	}

	mockService.EXPECT().
		Start(gomock.Any(), expectedRequest).
		Return(testNodeProps, nil).
		Times(1)

	handler, err := New(mockService)
	require.NoError(t, err)

	log, err := logging.New()
	require.NoError(t, err)

	requestBody, err := json.Marshal(testRequest)
	require.NoError(t, err)
	req := httptest.NewRequest("POST", "/start", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Start(log).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response apirequests.StartResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}
