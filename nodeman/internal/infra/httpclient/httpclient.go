package httpclient

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

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

func NewClientFactory(opts ...Option) *ClientFactory {
	cfg := config{
		dialerTimeout:       30 * time.Second,
		keepAlive:           30 * time.Second,
		maxIdleConns:        24,
		idleConnTimeout:     90 * time.Second,
		tlsHandshakeTimeout: 10 * time.Second,
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

func (cf *ClientFactory) GetNodeClient(certHash CertHash) (*http.Client, error) {
	if cf == nil {
		return nil, errdefs.NewNilCall()
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
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,

		VerifyPeerCertificate: verifyPeerFn(certHash),
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

func verifyPeerFn(certHash CertHash) func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		if len(rawCerts) == 0 {
			return errors.New("no certificate provided")
		}
		cert, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return err
		}
		sum := sha256Sum(cert.Raw)
		fmt.Println(base64.StdEncoding.EncodeToString(cert.Raw))
		fmt.Println(base64.StdEncoding.EncodeToString(sum[:]))
		fmt.Println(base64.StdEncoding.EncodeToString(certHash[:]))

		if sum != certHash {
			return errors.New("certificate pinning failed")
		}
		return nil
	}
}

func sha256Sum(data []byte) [32]byte {
	return sha256.Sum256(data)
}
