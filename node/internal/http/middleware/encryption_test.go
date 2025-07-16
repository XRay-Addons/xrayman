package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func encryptRequestBody(t *testing.T, body bytes.Buffer, key []byte) io.Reader {
	if len(key) == 0 {
		return &body
	}
	encrypted, err := jwe.Encrypt(
		body.Bytes(),
		jwe.WithKey(jwa.A256GCMKW(), key),
		jwe.WithContentEncryption(jwa.A256GCM()),
	)
	require.NoError(t, err, "Failed to encrypt test data")
	return bytes.NewReader(encrypted)
}

// test handler - unmarshall json body and return status ok
func testEncryptionHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err, "Failed to read request body")

		var decoded map[string]string
		err = json.Unmarshal(body, &decoded)
		require.NoError(t, err, "Failed to unmarshal decrypted body")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func TestEncryptionMiddleware(t *testing.T) {
	const testKey = "0123456789abcdef0123456789abcdef"
	const testFakeKey = "0123456789abcdef0123456789abcde0"
	var testBodyContent = map[string]string{"key": "value"}

	type testItem struct {
		name           string
		key            string
		bodyContent    map[string]string
		expectedStatus int
	}

	testItems := []testItem{
		{"true key", testKey, testBodyContent, http.StatusOK},
		{"fake key", testFakeKey, testBodyContent, http.StatusUnauthorized},
		{"no key", "", testBodyContent, http.StatusUnauthorized},
	}

	log, err := zap.NewDevelopment()
	require.NoError(t, err)

	for _, testItem := range testItems {
		t.Run(testItem.name, func(t *testing.T) {
			// create encryption middleware
			middleware := Encryption([]byte(testKey), log)(testEncryptionHandler(t))

			// create request body
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(testItem.bodyContent)
			require.NoError(t, err, "Failed to encode request body")

			// encrypt request body
			var encBody = encryptRequestBody(t, body, []byte(testItem.key))

			// create request
			req := httptest.NewRequest(http.MethodPost, "/", encBody)

			// handle request
			rec := httptest.NewRecorder()
			middleware.ServeHTTP(rec, req)

			// check status
			require.Equal(t, testItem.expectedStatus, rec.Code)

			// check response
			if testItem.expectedStatus != http.StatusOK {
				return
			}

			// check response body
			encryptedRes := rec.Body.Bytes()
			_, err := jwe.Decrypt(
				encryptedRes,
				jwe.WithKey(jwa.A256GCMKW(), []byte(testItem.key)),
			)
			require.NoError(t, err, "Failed to decrypt response")
		})
	}
}
