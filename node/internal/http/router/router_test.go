package router

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/http/handler"
	"github.com/XRay-Addons/xrayman/node/internal/http/handler/mocks"
	"github.com/XRay-Addons/xrayman/node/internal/http/security"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

func TestRouter(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		path            string
		body            []byte
		mockSetup       func(*mocks.MockService)
		expectedCode    int
		expectedBody    string
		checkOnlyStatus bool
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
			name:            "Post Validation Error",
			method:          http.MethodPost,
			path:            "/start",
			body:            []byte(`{"fusers":[]}`),
			mockSetup:       func(m *mocks.MockService) {},
			expectedCode:    http.StatusBadRequest,
			checkOnlyStatus: true,
		},
		{
			name:   "Get Panic",
			method: http.MethodGet,
			path:   "/status",
			body:   nil,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Status(gomock.Any(), gomock.Any()).
					Do(func(any, any) { panic("test panic") })
			},
			expectedCode:    http.StatusInternalServerError,
			checkOnlyStatus: true,
		},
		{
			name:   "Get Timeout",
			method: http.MethodGet,
			path:   "/status",
			body:   nil,
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Status(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, _ any) (any, any) {
						timer := time.NewTimer(4 * time.Second)
						select {
						case <-timer.C:
							return nil, nil
						case <-ctx.Done():
							return nil, fmt.Errorf("test timeout")
						}
					})
			},
			expectedCode:    http.StatusInternalServerError,
			checkOnlyStatus: true,
		},
	}

	log := zaptest.NewLogger(t)
	secutiry := security.NewBackdoor()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockService(ctrl)
			tt.mockSetup(mockService)

			handler, err := handler.New(mockService, log)
			require.NoError(t, err)

			h, err := New(handler, secutiry,
				WithLogger(log),
				WithTimeout(2*time.Second),
			)
			require.NoError(t, err)

			req, err := http.NewRequest(tt.method, tt.path, bytes.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test")
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code)

			if tt.checkOnlyStatus {
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
