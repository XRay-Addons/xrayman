package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XRay-Addons/xrayman/node/internal/logging"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func encryptRequestBody(t *testing.T, body []byte, key []byte) io.Reader {
	if len(key) == 0 || len(body) == 0 {
		return bytes.NewReader(body)
	}
	encrypted, err := jwe.Encrypt(
		body,
		jwe.WithKey(jwa.A256GCMKW(), key),
		jwe.WithContentEncryption(jwa.A256GCM()),
	)
	require.NoError(t, err, "Failed to encrypt test data")
	return bytes.NewReader(encrypted)
}

func testEncryptionSendRequest(
	t *testing.T,
	handler http.Handler,
	key string, correctKey bool,
	body []byte,
	rec *httptest.ResponseRecorder,
) {
	// create request with encrypted body
	encBody := encryptRequestBody(t, body, []byte(key))
	req := httptest.NewRequest(http.MethodPost, "/", encBody)

	// handle request
	handler.ServeHTTP(rec, req)

	if !correctKey || rec.Body.Len() == 0 {
		return
	}
	if rec.Code != http.StatusOK {
		return
	}

	// if key is correct - decode response body
	decryptedBody, err := jwe.Decrypt(
		rec.Body.Bytes(),
		jwe.WithKey(jwa.A256GCMKW(), []byte(key)),
	)
	require.NoError(t, err, "Failed to decrypt response")
	rec.Body = bytes.NewBuffer(decryptedBody)
}

// test encryption and auth
func TestEncryption(t *testing.T) {
	const testKey = "0123456789abcdef0123456789abcdef"
	const testFakeKey = "0123456789abcdef0123456789abcde0"

	type testItem struct {
		name           string
		key            string
		reqBody        []byte
		respBody       []byte
		expectedStatus int
		expectedBody   []byte
	}

	testItems := []testItem{
		{
			"true key",
			testKey,
			[]byte("request body"),
			[]byte("response body"),
			http.StatusOK,
			[]byte("response body"),
		},
		{
			"fake key",
			testFakeKey,
			[]byte("request body"),
			[]byte("response body"),
			http.StatusUnauthorized,
			[]byte(`{"error":"Unauthorized","details":"invalid content JWE"}`),
		},
		{
			"no key",
			"",
			[]byte("request body"),
			[]byte("response body"),
			http.StatusUnauthorized,
			[]byte(`{"error":"Unauthorized","details":"invalid content JWE"}`),
		},
		{
			"empty request",
			testKey,
			[]byte(""),
			[]byte("response body"),
			http.StatusOK,
			[]byte("response body"),
		},
		{
			"empty response",
			testKey,
			[]byte("request body"),
			nil,
			http.StatusOK,
			nil,
		},
	}

	log, err := logging.New()
	require.NoError(t, err)

	for _, tt := range testItems {
		t.Run(tt.name, func(t *testing.T) {
			// handler for test request and response
			testHandler := func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err, "Failed to read request body")
				assert.Equal(t, tt.reqBody, body)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.respBody))
			}

			handler := Encryption([]byte(testKey), log)(http.HandlerFunc(testHandler))

			// send request, record response
			rec := httptest.NewRecorder()
			testEncryptionSendRequest(t, handler,
				tt.key, tt.key == testKey,
				tt.reqBody, rec)

			// check response status
			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				// check response body
				require.Equal(t, []byte(tt.respBody), rec.Body.Bytes())
			} else {
				// check response error body
				require.JSONEq(t, string(tt.expectedBody), rec.Body.String())
			}
		})
	}
}

// test errors not encrypted
func TestEncryptionError(t *testing.T) {
	const testKey = "0123456789abcdef0123456789abcdef"

	type testItem struct {
		name       string
		reqBody    []byte
		respBody   []byte
		respStatus int
	}

	testItems := []testItem{
		{
			"err 451",
			[]byte("request body"),
			[]byte(`{"error":"451 error","details":"no details"}`),
			451,
		},
	}

	log, err := logging.New()
	require.NoError(t, err)

	for _, tt := range testItems {
		t.Run(tt.name, func(t *testing.T) {
			// handler for test request and response
			testHandler := func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err, "Failed to read request body")
				assert.Equal(t, tt.reqBody, body)
				w.WriteHeader(tt.respStatus)
				w.Write([]byte(tt.respBody))
			}

			handler := Encryption([]byte(testKey), log)(http.HandlerFunc(testHandler))

			// send request, record response
			rec := httptest.NewRecorder()
			testEncryptionSendRequest(t, handler,
				testKey, true, tt.reqBody, rec)

			// check response status
			require.Equal(t, tt.respStatus, rec.Code)

			// check response error body
			require.JSONEq(t, string(tt.respBody), rec.Body.String())
		})
	}
}
