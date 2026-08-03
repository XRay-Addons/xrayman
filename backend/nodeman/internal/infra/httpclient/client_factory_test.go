package httpclient

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestBlackholeTimeout(t *testing.T) {
	logger := zaptest.NewLogger(t)
	cf, err := NewClientFactory(5*time.Second, logger)
	require.NoError(t, err)
	client, err := cf.GetNodeClient(models.CertHash{})
	require.NoError(t, err)

	// --- HTTP request ---
	req, err := http.NewRequest("GET", "https://10.255.255.1:443", nil)
	require.NoError(t, err)

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	t.Logf("elapsed: %v", elapsed)
	if err != nil {
		t.Logf("error: %v", err)
	} else {
		t.Logf("response: %+v", resp)
		resp.Body.Close()
	}
}

func TestHandshakeTimeout(t *testing.T) {
	logger := zaptest.NewLogger(t)
	cf, err := NewClientFactory(5*time.Second, logger)
	require.NoError(t, err)
	client, err := cf.GetNodeClient(models.CertHash{})
	require.NoError(t, err)

	// --- TCP server that accepts but never responds ---
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_ = conn
		}
	}()

	// --- HTTP request ---
	req, err := http.NewRequest("GET", "https://"+ln.Addr().String(), nil)
	require.NoError(t, err)

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	t.Logf("elapsed: %v", elapsed)
	if err != nil {
		t.Logf("error: %v", err)
	} else {
		t.Logf("response: %+v", resp)
		resp.Body.Close()
	}
}

func TestResponseTimeout(t *testing.T) {

	logger := zaptest.NewLogger(t)

	cf, err := NewClientFactory(7*time.Second, logger)
	require.NoError(t, err)

	client, err := cf.GetNodeClient(models.CertHash{})
	require.NoError(t, err)

	// --- server: TLS + hang after accept ---
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(10 * time.Second):
			return
		case <-r.Context().Done():
			return
		}
	}))
	defer ts.Close()

	client.Transport = ts.Client().Transport

	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	t.Logf("elapsed: %v", elapsed)

	if err != nil {
		t.Logf("error: %v", err)
		return
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)

	t.Logf("read error: %v", err)
}
