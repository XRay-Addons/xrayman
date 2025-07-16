package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/http/errors"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"go.uber.org/zap"
)

func Encryption(jwekey []byte, log *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		jwtauth := func(w http.ResponseWriter, r *http.Request) {
			// reader with decryption
			dr, err := withDecryption(r, jwekey)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// writer with encryption
			ew, err := withEncryption(w, jwekey)
			if err != nil {
				errors.LogRequestError(log, r, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// process request
			next.ServeHTTP(ew, dr)

			// write responce
			if err := ew.FlushToClient(); err != nil {
				errors.LogRequestError(log, r, err)
			}
		}

		return http.HandlerFunc(jwtauth)
	}
}

func withDecryption(r *http.Request, key []byte) (*http.Request, error) {
	// decrypt request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("request body decryption: %w", err)
	}

	if len(body) == 0 {
		// nothing to read, nothing to decode
		return r, nil
	}

	decryptedBody, err := jwe.Decrypt(body,
		jwe.WithKey(jwa.A256GCMKW(), key))

	if err != nil {
		return nil, fmt.Errorf("request body decryption: %w", err)
	}

	// replace body to decoded
	decodedReq := r.Clone(r.Context())
	decodedReq.Body = io.NopCloser(bytes.NewReader(decryptedBody))
	decodedReq.ContentLength = int64(len(decryptedBody))
	decodedReq.TransferEncoding = nil

	return decodedReq, nil
}

type encodedWriter struct {
	http.ResponseWriter
	key []byte
	buf bytes.Buffer
}

var _ http.ResponseWriter = (*encodedWriter)(nil)

func withEncryption(baseWriter http.ResponseWriter, key []byte) (*encodedWriter, error) {
	return &encodedWriter{
		ResponseWriter: baseWriter,
		key:            key,
	}, nil
}

func (w *encodedWriter) Write(data []byte) (int, error) {
	return w.buf.Write(data)
}

func (w *encodedWriter) FlushToClient() error {
	if w.buf.Len() == 0 {
		return nil
	}

	// encode content
	encodedContent, err := jwe.Encrypt(
		w.buf.Bytes(),
		jwe.WithKey(jwa.A256GCMKW(), w.key),
		jwe.WithContentEncryption(jwa.A256GCM()))
	if err != nil {
		return fmt.Errorf("encoding content: %w", err)
	}

	// write encoded content
	if _, err = w.ResponseWriter.Write(encodedContent); err != nil {
		return fmt.Errorf("write encoded content: %w", err)
	}
	w.buf.Reset()

	return nil
}
