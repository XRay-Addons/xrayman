package seccfg

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/fs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

const (
	CertFileName   = "cert.pem"
	KeyFileName    = "key.pem"
	AccessFileName = "access.json"
	RSAKeySize     = 2048
	JWTSecretSize  = 32
)

type SecurityConfig struct {
	AccessKey models.AccessKey
	Cert      []byte
	Key       []byte
}

func New(dir string) (*SecurityConfig, error) {
	certFile := path.Join(dir, CertFileName)
	keyFile := path.Join(dir, KeyFileName)
	accessFile := path.Join(dir, AccessFileName)

	exists := true
	for _, f := range []string{certFile, keyFile, accessFile} {
		fexists, err := fs.AccessFile(f)
		if err != nil {
			return nil, fmt.Errorf("security: init: %w", err)
		}
		exists = exists && fexists
	}

	issuer := "xray node"
	expire := time.Hour * 24 * 365 * 10

	if !exists {
		if err := generateSecurity(dir, issuer, expire); err != nil {
			return nil, fmt.Errorf("security: init: %w", err)
		}
	}

	cert, err := os.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("security: init: %w", err)
	}
	key, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("security: init: %w", err)
	}
	accessKey, err := readAccessKey(accessFile)
	if err != nil {
		return nil, fmt.Errorf("security: init: %w", err)
	}

	security := &SecurityConfig{
		Cert:      cert,
		Key:       key,
		AccessKey: accessKey,
	}

	if err := security.validate(); err != nil {
		return nil, fmt.Errorf("security: init: %w", err)
	}

	return security, nil
}

func (s *SecurityConfig) validate() error {
	if len(s.Cert) == 0 {
		return fmt.Errorf("certificate is empty: %w", errdefs.ErrConfig)
	}
	if len(s.Key) == 0 {
		return fmt.Errorf("private key is empty: %w", errdefs.ErrConfig)
	}

	block, _ := pem.Decode(s.Cert)
	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("invalid certificate PEM: %w", errdefs.ErrConfig)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("certificate parsing: %v: %w", err, errdefs.ErrConfig)
	}

	actualHash := sha256.Sum256(cert.Raw)
	if actualHash != s.AccessKey.CertHash {
		return fmt.Errorf("certificate hash mismatch: %w", errdefs.ErrConfig)
	}
	fmt.Println("node cert hash:", base64.StdEncoding.EncodeToString(actualHash[:]))

	return nil
}
