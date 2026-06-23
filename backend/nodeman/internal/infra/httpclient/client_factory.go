package httpclient

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type config struct {
	dialerTimeout       time.Duration
	tlsHandshakeTimeout time.Duration
	clientTimeout       time.Duration

	keepAlive       time.Duration
	maxIdleConns    int
	idleConnTimeout time.Duration
	log             *zap.Logger
}

type ClientFactory struct {
	cfg         config
	dialer      net.Dialer
	clientsPool map[models.CertHash]*http.Client
	mu          sync.RWMutex
}

const (
	defaultKeepAlive       = 30 * time.Second
	defaultMaxIdleConns    = 24
	defaultIdleConnTimeout = 90 * time.Second
)

func NewClientFactory(timeout time.Duration, log *zap.Logger) (*ClientFactory, error) {
	if timeout == 0 {
		return nil, xerr.NilArg("timeout")
	}
	if log == nil {
		return nil, xerr.NilArg("log")
	}
	cfg := config{
		dialerTimeout:       max(2*time.Second, timeout/4),
		tlsHandshakeTimeout: max(2*time.Second, timeout/4),
		clientTimeout:       timeout,

		keepAlive:       defaultKeepAlive,
		maxIdleConns:    defaultMaxIdleConns,
		idleConnTimeout: defaultIdleConnTimeout,
	}

	return &ClientFactory{
		cfg: cfg,
		dialer: net.Dialer{
			Timeout:   cfg.dialerTimeout,
			KeepAlive: cfg.keepAlive,
		},
		clientsPool: make(map[models.CertHash]*http.Client),
	}, nil
}

func (cf *ClientFactory) Close() {
	if cf == nil {
		return
	}
	for _, c := range cf.clientsPool {
		c.CloseIdleConnections()
	}
}

func (cf *ClientFactory) GetNodeClient(certHash models.CertHash) (*http.Client, error) {
	if cf == nil {
		return nil, errdefs.NilCall()
	}

	// fast check with RLock
	cf.mu.RLock()
	client, found := cf.clientsPool[certHash]
	cf.mu.RUnlock()
	if found {
		return client, nil
	}

	// second check with RWLock
	cf.mu.Lock()
	defer cf.mu.Unlock()
	if client, found = cf.clientsPool[certHash]; found {
		return client, nil
	}

	client = cf.newHttpClient(certHash)
	cf.clientsPool[certHash] = client

	return client, nil
}

func (cf *ClientFactory) newHttpClient(certHash models.CertHash) *http.Client {
	verifyFn := verifyPeerFn(certHash)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // #nosec custom verification used
		VerifyConnection: func(cs tls.ConnectionState) error {
			rawCerts := make([][]byte, len(cs.PeerCertificates))
			for i, cert := range cs.PeerCertificates {
				rawCerts[i] = cert.Raw
			}
			return verifyFn(rawCerts)
		},
	}

	baseTransport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         cf.dialer.DialContext,
		MaxIdleConns:        cf.cfg.maxIdleConns,
		IdleConnTimeout:     cf.cfg.idleConnTimeout,
		TLSHandshakeTimeout: cf.cfg.tlsHandshakeTimeout,
		TLSClientConfig:     tlsConfig,
	}

	transport := &httpTransport{
		base: baseTransport,
		log:  cf.cfg.log,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   cf.cfg.clientTimeout,
	}

	return client
}

func verifyPeerFn(certHash models.CertHash) func(rawCerts [][]byte) error {
	return func(rawCerts [][]byte) error {
		if len(rawCerts) == 0 {
			return xerr.New("no certificate provided")
		}
		cert, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return xerr.WrapWithStack(err)
		}
		sum := sha256Sum(cert.Raw)
		if sum != certHash {
			return xerr.New("certificate pinning failed")
		}
		return nil
	}
}

func sha256Sum(data []byte) [32]byte {
	return sha256.Sum256(data)
}
