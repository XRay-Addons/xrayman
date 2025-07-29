package client

import (
	"context"
	"net/http/httptest"
	"testing"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	mocks "github.com/XRay-Addons/xrayman/nodeman/internal/clients/node/mocks"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -destination=./mocks/mock_security_handler.go -package=mocks github.com/XRay-Addons/xrayman/node/pkg/api/http/gen SecurityHandler
//go:generate mockgen -destination=./mocks/mock_security_source.go -package=mocks github.com/XRay-Addons/xrayman/node/pkg/api/http/gen SecuritySource
//go:generate mockgen -destination=./mocks/mock_handler.go -package=mocks github.com/XRay-Addons/xrayman/node/pkg/api/http/gen Handler

func TestNodeClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSecuritySource := mocks.NewMockSecuritySource(ctrl)
	mockSecurityHandler := mocks.NewMockSecurityHandler(ctrl)
	mockHandler := mocks.NewMockHandler(ctrl)

	testAuthToken := api.BearerAuth{Token: "test token"}

	mockSecuritySource.
		EXPECT().
		BearerAuth(gomock.Any(), gomock.Any()).
		Return(testAuthToken, nil)

	mockSecurityHandler.
		EXPECT().
		HandleBearerAuth(gomock.Any(), "GetStatus", gomock.Any()).
		Do(func(ctx context.Context, op string, auth api.BearerAuth) {
			assert.Equal(t, testAuthToken.Token, auth.Token)
		}).
		Return(context.TODO(), nil)

	mockHandler.
		EXPECT().
		GetStatus(gomock.Any()).
		Return(&api.StatusResponse{ServiceStatus: api.ServiceStatusRunning}, nil)

	// run mock server
	handler, err := api.NewServer(mockHandler, mockSecurityHandler)
	require.NoError(t, err)
	testServer := httptest.NewServer(handler)

	// create client
	nodeAPI, err := NewNodeClient(testServer.URL, mockSecuritySource, nil)
	require.NoError(t, err)

	// call test method
	status, err := nodeAPI.CheckStatus(context.TODO())
	require.NoError(t, err)
	require.Equal(t, status, models.NodeStatusRunning)
}
