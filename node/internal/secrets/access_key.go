package secrets

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

func ensureAccessKey(dir string) (*models.AccessKey, error) {
	certHash, err := getCertHash(dir)
	if err != nil {
		return nil, err
	}

	// if access key exists and matches cert - ok
	accessKey, err := readAccessKey(dir)
	if err == nil && accessKey.CertHash == *certHash {
		return accessKey, nil
	}

	// generate new access key
	accessSecret, err := generateAccessSecret()
	if err != nil {
		return nil, err
	}
	accessKey = &models.AccessKey{
		CertHash:     *certHash,
		AccessSecret: *accessSecret,
	}

	// write it
	if err := writeAccessKey(dir, *accessKey); err != nil {
		return nil, err
	}

	return accessKey, nil
}

func getCertHash(dir string) (*models.CertHash, error) {
	certPEM, err := os.ReadFile(filepath.Join(dir, CertFile))
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, errdefs.New("get cert hash: invalid certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	certHash := sha256.Sum256(cert.Raw)

	return &certHash, nil
}

type accessKeyWrapper struct {
	AccessKey models.AccessKey `json:"access_key"`
}

func readAccessKey(dir string) (*models.AccessKey, error) {
	data, err := os.ReadFile(filepath.Join(dir, AccessFile))
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	var wrapper accessKeyWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, errdefs.WithStack(err)
	}
	return &wrapper.AccessKey, nil
}

func writeAccessKey(dir string, key models.AccessKey) error {
	wrapper := accessKeyWrapper{AccessKey: key}

	data, err := json.MarshalIndent(&wrapper, "", "  ")
	if err != nil {
		return errdefs.WithStack(err)
	}
	if err := os.WriteFile(filepath.Join(dir, AccessFile), data, 0600); err != nil {
		return errdefs.WithStack(err)
	}
	return nil
}

func generateAccessSecret() (*models.AccessSecret, error) {
	var secret models.AccessSecret
	if _, err := rand.Read(secret[:]); err != nil {
		return nil, errdefs.WithStack(err)
	}
	return &secret, nil
}
