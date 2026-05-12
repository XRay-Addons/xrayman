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
)

type Option func(opts *config)

type CertHash = [32]byte

type ClientFactory struct {
	cfg         config
	dialer      net.Dialer
	clientsPool map[CertHash]*http.Client
	mu          sync.RWMutex
}

type config struct {
	dialerTimeout       time.Duration
	keepAlive           time.Duration
	maxIdleConns        int
	idleConnTimeout     time.Duration
	tlsHandshakeTimeout time.Duration
}

const (
	defaultDialerTimeout       = 30 * time.Second
	defaultKeepAlive           = 30 * time.Second
	defaultMaxIdleConns        = 24
	defaultIdleConnTimeout     = 90 * time.Second
	defaultTlsHandshakeTimeout = 10 * time.Second
)

func NewClientFactory(opts ...Option) *ClientFactory {
	cfg := config{
		dialerTimeout:       defaultDialerTimeout,
		keepAlive:           defaultKeepAlive,
		maxIdleConns:        defaultMaxIdleConns,
		idleConnTimeout:     defaultIdleConnTimeout,
		tlsHandshakeTimeout: defaultTlsHandshakeTimeout,
	}
	for _, o := range opts {
		o(&cfg)
	}
	return &ClientFactory{
		cfg: cfg,
		dialer: net.Dialer{
			Timeout:   cfg.dialerTimeout,
			KeepAlive: cfg.keepAlive,
		},
		clientsPool: make(map[CertHash]*http.Client),
	}
}

func (cf *ClientFactory) Close() {
	if cf == nil {
		return
	}
	for _, c := range cf.clientsPool {
		c.CloseIdleConnections()
	}
}

func (cf *ClientFactory) GetNodeClient(certHash CertHash) (*http.Client, error) {
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

func (cf *ClientFactory) newHttpClient(certHash CertHash) *http.Client {
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

	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         cf.dialer.DialContext,
		MaxIdleConns:        cf.cfg.maxIdleConns,
		IdleConnTimeout:     cf.cfg.idleConnTimeout,
		TLSHandshakeTimeout: cf.cfg.tlsHandshakeTimeout,
		TLSClientConfig:     tlsConfig,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   cf.cfg.dialerTimeout + cf.cfg.tlsHandshakeTimeout,
	}
}

func verifyPeerFn(certHash CertHash) func(rawCerts [][]byte) error {
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
