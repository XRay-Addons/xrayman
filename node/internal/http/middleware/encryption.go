package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/http/httperr"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"go.uber.org/zap"
)

func Encryption(jwekey []byte, log *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		jweenc := func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			var err error
			defer func() { httperr.Write(r.Context(), err, w, log) }()

			// dectypt req
			decR, err := decodeReq(r, jwekey)
			if err != nil {
				err = httperr.New(httperr.ErrContentEncryption, err)
				return
			}

			// write result to writer with allowed encoding
			encW := newEncodableWriter(w)
			next.ServeHTTP(encW, decR)

			// don't encode failed requests, flush as is
			if encW.StatusCode() != http.StatusOK {
				encW.Flush()
				return
			}

			// encode request
			if err = encW.Encode(jwekey); err != nil {
				return
			}

			encW.Flush()
		}

		return http.HandlerFunc(jweenc)
	}
}

func decodeReq(r *http.Request, key []byte) (*http.Request, error) {
	// read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("request body read: %w", err)
	}
	if len(body) == 0 {
		// nothing to read, nothing to decode
		return r, nil
	}
	defer r.Body.Close()

	// decode request body
	decryptedBody, err := jwe.Decrypt(body,
		jwe.WithKey(jwa.A256GCMKW(), key))
	if err != nil {
		return nil, fmt.Errorf("request body decryption: %w", err)
	}

	// create request clone with decoded body
	decodedReq := r.Clone(r.Context())
	decodedReq.Body = io.NopCloser(bytes.NewReader(decryptedBody))
	decodedReq.ContentLength = int64(len(decryptedBody))
	decodedReq.TransferEncoding = nil

	return decodedReq, nil
}

type encodableWriter struct {
	baseWriter http.ResponseWriter
	buf        bytes.Buffer
	statusCode int
}

func newEncodableWriter(baseWriter http.ResponseWriter) *encodableWriter {
	return &encodableWriter{
		baseWriter: baseWriter,
		statusCode: http.StatusOK,
	}
}

func (e *encodableWriter) Header() http.Header {
	return e.baseWriter.Header()
}

func (e *encodableWriter) Write(p []byte) (int, error) {
	return e.buf.Write(p)
}

func (e *encodableWriter) WriteHeader(statusCode int) {
	e.statusCode = statusCode
}

func (e *encodableWriter) StatusCode() int {
	return e.statusCode
}

// encode buffer
func (e *encodableWriter) Encode(key []byte) error {
	if e.buf.Len() == 0 {
		return nil
	}
	// try to encode buffer
	encodedContent, err := jwe.Encrypt(
		e.buf.Bytes(),
		jwe.WithKey(jwa.A256GCMKW(), key),
		jwe.WithContentEncryption(jwa.A256GCM()))
	if err != nil {
		return err
	}
	e.buf = *bytes.NewBuffer(encodedContent)
	return nil
}

// flush to base writer
func (e *encodableWriter) Flush() {
	// ignore flushing errors, we can't do anything with it
	e.baseWriter.WriteHeader(e.statusCode)
	_, _ = e.baseWriter.Write(e.buf.Bytes())
}

var _ http.ResponseWriter = (*encodableWriter)(nil)
