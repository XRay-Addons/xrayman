package httpclient

import (
	"net/http"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"go.uber.org/zap"
)

type httpTransport struct {
	base *http.Transport
	log  *zap.Logger
}

func (t *httpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	t.log.Info("http response",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Int("status", resp.StatusCode),
	)

	return resp, nil
}
