package httperr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// test errors.As/Is for wrapped response
func TestResponseError(t *testing.T) {
	// test Response
	resp := ErrContentValidation
	require.True(t, errors.As(resp, new(*Response)))
	// test fmt-wrapped Response
	wrappedResp := fmt.Errorf("wrapped resp: %w", resp)
	require.True(t, errors.As(wrappedResp, new(*Response)))

	// create Error with reason
	reason := errors.New("test error reason")
	reasoned := New(resp, reason)

	// test response extraction from reasoned error
	require.True(t, errors.As(reasoned, new(*Response)))
	// test fmt-wrapped reasoned error
	wrappedReasoned := fmt.Errorf("wrapped reasoned: %w", reasoned)
	require.True(t, errors.As(wrappedReasoned, new(*Response)))

	// test error text
	require.Equal(t,
		"wrapped reasoned: test error reason",
		wrappedReasoned.Error(),
	)
}
