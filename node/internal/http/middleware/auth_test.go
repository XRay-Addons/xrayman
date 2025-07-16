package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/logging"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAuthHeader = "Authorization"

func signAuthJWT(t *testing.T, header http.Header, key []byte) {
	if len(key) == 0 {
		return
	}
	token, err := jwt.NewBuilder().Issuer("test").IssuedAt(time.Now()).Build()
	require.NoError(t, err)

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256(), []byte(key)))
	require.NoError(t, err)
	header.Set(testAuthHeader, string(signed))
}

func testAuthSendRequest(
	t *testing.T,
	handler http.Handler,
	key string, correctKey bool,
	rec *httptest.ResponseRecorder,
) {
	// create request
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// sigh request header
	signAuthJWT(t, req.Header, []byte(key))

	// handle request
	handler.ServeHTTP(rec, req)

	if !correctKey {
		return
	}

	// if key is correct - check auth header and decode body
	validationToken, err := jwt.ParseHeader(rec.Result().Header, testAuthHeader,
		jwt.WithKey(jwa.HS256(), []byte(key)))
	_ = validationToken
	require.NoError(t, err)
}

// test auth
func TestAuth(t *testing.T) {
	const testKey = "0123456789abcdef0123456789abcdef"
	const testFakeKey = "0123456789abcdef0123456789abcde0"

	type testItem struct {
		name           string
		key            string
		expectedStatus int
	}

	testItems := []testItem{
		{
			"true key POST",
			testKey,
			http.StatusOK,
		},
		{
			"fake key POST",
			testFakeKey,
			http.StatusUnauthorized,
		},
		{
			"no key POST",
			"",
			http.StatusUnauthorized,
		},
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "Failed to read request body")
		w.WriteHeader(http.StatusOK)
	}

	log, err := logging.New()
	require.NoError(t, err)

	for _, tt := range testItems {
		t.Run(tt.name, func(t *testing.T) {
			handler := Auth([]byte(testKey), log)(http.HandlerFunc(testHandler))

			// send request, record response
			rec := httptest.NewRecorder()
			testAuthSendRequest(t, handler,
				tt.key, tt.key == testKey, rec)

			// check response status
			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
