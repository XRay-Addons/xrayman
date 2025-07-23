package handler

import (
	"bytes"
	"fmt"
	"mime"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XRay-Addons/xrayman/node/internal/http/handler/mocks"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		body         []byte
		mockSetup    func(*mocks.MockService)
		expectedCode int
		expectedBody string
	}{
		{
			name:   "Get OK",
			method: http.MethodGet,
			path:   "/status",
			body:   nil,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Status(gomock.Any(), gomock.Any()).
					Return(&models.StatusResult{ServiceStatus: models.ServiceRunning}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"serviceStatus":"running"}`,
		},
		{
			name:   "Get InternalError",
			method: http.MethodGet,
			path:   "/status",
			body:   nil,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Status(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("test error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"message":"Internal server error"}`,
		},
		{
			name:   "Post OK",
			method: http.MethodPost,
			path:   "/start",
			body:   []byte(`{"users":[]}`),
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Start(gomock.Any(), gomock.Any()).
					Return(&models.StartResult{}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"clientCfg": {"template":"", "userNameField":"", "vlessUUIDField":""}}`,
		},
		{
			name:         "Post Error",
			method:       http.MethodPost,
			path:         "/start",
			body:         []byte(`{"fusers":[]}`),
			mockSetup:    func(m *mocks.MockService) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockService(ctrl)
			tt.mockSetup(mockService)

			h, err := NewHandlerImpl(mockService)
			require.NoError(t, err)

			srv, err := api.NewServer(h)
			require.NoError(t, err)

			req, err := http.NewRequest(tt.method, tt.path, bytes.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			srv.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedCode != http.StatusOK {
				return
			}

			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			expectedContentType := "application/json"
			mt, _, err := mime.ParseMediaType(rr.Header().Get("Content-Type"))
			require.NoError(t, err)
			assert.Equal(t, expectedContentType, mt)
		})
	}
}
