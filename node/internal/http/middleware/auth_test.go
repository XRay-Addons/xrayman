package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func signJWT(t *testing.T, key []byte) string {
	t.Helper()
	token, err := jwt.NewBuilder().
		Issuer("test").
		IssuedAt(time.Now()).
		Build()
	require.NoError(t, err)

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256(), []byte(key)))
	require.NoError(t, err)
	return string(signed)
}

// test handler - do nothing
func testAuthHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	}
}

func TestAuth(t *testing.T) {
	const testKey = "0123456789abcdef0123456789f"
	const testFakeKey = "0123456789abcdef01234567890"

	type testItem struct {
		name           string
		key            string
		expectedStatus int
	}

	testItems := []testItem{
		{"true key", testKey, http.StatusOK},
		{"fake key", testFakeKey, http.StatusUnauthorized},
		{"no key", "", http.StatusUnauthorized},
	}

	log, err := zap.NewDevelopment()
	require.NoError(t, err)

	for _, testItem := range testItems {
		t.Run(testItem.name, func(t *testing.T) {
			// request with encrypted body and signed header
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			if len(testItem.key) > 0 {
				signedJWT := signJWT(t, []byte(testItem.key))
				req.Header.Set("Authorization", signedJWT)
			}

			// test handler with middleware
			handler := Auth([]byte(testKey), log)(testAuthHandler(t))

			// handle request
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			// check response status
			require.Equal(t, testItem.expectedStatus, rec.Result().StatusCode)

			// check response header for success requests
			if testItem.expectedStatus != http.StatusOK {
				return
			}

			// check response header
			require.NotEmpty(t, rec.Result().Header.Get("Authorization"))

			validationToken, err := jwt.ParseHeader(rec.Result().Header,
				authHeader,
				jwt.WithKey(jwa.HS256(), []byte(testItem.key)))
			_ = validationToken
			require.NoError(t, err)
		})
	}
}
