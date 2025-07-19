package httperr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWrite(t *testing.T) {
	type testItem struct {
		err             error
		expectedCode    int
		expectedDetails string
	}

	testItems := []testItem{
		{
			err:             errors.New("reason"),
			expectedCode:    ErrUnknownError.Code(),
			expectedDetails: ErrUnknownError.Error(),
		},
		{
			err:             ErrContentValidation,
			expectedCode:    ErrContentValidation.Code(),
			expectedDetails: ErrContentValidation.Error(),
		},
		{
			err:             fmt.Errorf("wrap: %w", ErrContentValidation),
			expectedCode:    ErrContentValidation.Code(),
			expectedDetails: ErrContentValidation.Error(),
		},
		{
			err:             New(ErrContentValidation, errors.New("reason")),
			expectedCode:    ErrContentValidation.Code(),
			expectedDetails: ErrContentValidation.Error(),
		},
		{
			err:             fmt.Errorf("wrap: %w", New(ErrContentValidation, errors.New("reason"))),
			expectedCode:    ErrContentValidation.Code(),
			expectedDetails: ErrContentValidation.Error(),
		},
	}

	log, err := zap.NewProduction()
	require.NoError(t, err)

	for _, tt := range testItems {
		t.Run("", func(t *testing.T) {
			rec := httptest.NewRecorder()
			Write(context.TODO(), tt.err, rec, log)

			// test error code
			require.Equal(t, tt.expectedCode, rec.Code)

			// test content
			var respContent map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &respContent)
			require.NoError(t, err)

			require.Equal(t, tt.expectedDetails, respContent["details"])
			require.Equal(t, http.StatusText(tt.expectedCode), respContent["error"])
		})
	}
}
