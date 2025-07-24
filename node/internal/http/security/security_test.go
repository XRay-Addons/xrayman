package security

import (
	"mime"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/http/httperr"
	"github.com/XRay-Addons/xrayman/node/internal/http/security/mocks"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=../../../pkg/api/http/gen/oas_server_gen.go -destination=./mocks/mock_handler.go -package=mocks

func TestSecurity(t *testing.T) {

	//authErr := api.ErrorStatusCode(*httperr.ErrAuthToken)

	tests := []struct {
		name          string
		path          string
		requestSecret []byte
		serverSecret  []byte
		mockSetup     func(*mocks.MockHandler)
		expectedCode  int
		expectedBody  string
	}{
		{
			name:          "Security OK",
			path:          "/status",
			requestSecret: []byte("very-secret-access-key"),
			serverSecret:  []byte("very-secret-access-key"),
			mockSetup: func(m *mocks.MockHandler) {
				m.EXPECT().
					GetStatus(gomock.Any()).
					Return(&api.StatusResponse{ServiceStatus: api.ServiceStatusRunning}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"serviceStatus":"running"}`,
		},
		{
			name:          "Security Failed",
			path:          "/status",
			requestSecret: []byte("very-secret-access-key"),
			serverSecret:  []byte("very-secret-access-lay"),
			mockSetup: func(m *mocks.MockHandler) {
				m.EXPECT().
					NewError(gomock.Any(), gomock.Any()).
					DoAndReturn(func(any, any) error {
						authErr := api.ErrorStatusCode(*httperr.ErrAuthToken)
						return &authErr
					})
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"details":"try another one", "message":"Invalid auth token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHandler := mocks.NewMockHandler(ctrl)
			tt.mockSetup(mockHandler)

			srv, err := api.NewServer(mockHandler, New(tt.serverSecret))
			require.NoError(t, err)

			claims := jwt.MapClaims{
				"sub": "test user",
				"exp": time.Now().Add(time.Hour).Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenSign, err := token.SignedString(tt.requestSecret)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tokenSign)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			srv.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			expectedContentType := "application/json"
			mt, _, err := mime.ParseMediaType(rr.Header().Get("Content-Type"))
			require.NoError(t, err)
			assert.Equal(t, expectedContentType, mt)
		})
	}
}
