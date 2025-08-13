package seccfg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

func generateSecurity(dir string, issuer string, expire time.Duration) error {
	// generate tls cert
	privKey, err := generateRSAKey()
	if err != nil {
		return err
	}
	keyPEM, err := encodePrivateKeyPEM(privKey)
	if err != nil {
		return err
	}
	certDER, err := createSelfSignedCert(privKey, issuer, expire)
	if err != nil {
		return err
	}
	certPEM := encodeCertPEM(certDER)
	certHash := getCertHash(certDER)

	accessSecret, err := createAccessSecret()
	if err != nil {
		return err
	}

	accessKey := models.AccessKey{
		CertHash:     certHash,
		AccessSecret: accessSecret,
	}

	if err := ensureDir(dir); err != nil {
		return err
	}

	if err := writeFiles(dir, certPEM, keyPEM, accessKey); err != nil {
		return err
	}

	return nil
}

func generateRSAKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, RSAKeySize)
}

func createSelfSignedCert(privKey *rsa.PrivateKey, issuer string, expire time.Duration) ([]byte, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("createSelfSignedCert: generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{issuer},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(expire),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, fmt.Errorf("createSelfSignedCert: create certificate: %w", err)
	}
	return certDER, err

	//certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	//return certPEM, nil
}

func getCertHash(certDER []byte) models.CertHash {
	return sha256.Sum256(certDER)
}

func createAccessSecret() (models.AccessSecret, error) {
	var secret [32]byte
	if _, err := rand.Read(secret[:]); err != nil {
		return secret, fmt.Errorf("createAccessSecret: generate secret: %w", err)
	}
	return secret, nil
}

func encodePrivateKeyPEM(privKey *rsa.PrivateKey) ([]byte, error) {
	keyDER := x509.MarshalPKCS1PrivateKey(privKey)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyDER})
	return keyPEM, nil
}

func encodeCertPEM(certDER []byte) []byte {
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	return certPEM
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0o700)
}

func writeFiles(dir string, certPEM, keyPEM []byte, accessKey models.AccessKey) error {
	if err := os.WriteFile(path.Join(dir, CertFileName), certPEM, 0o600); err != nil {
		return fmt.Errorf("writeFiles: write cert: %w", err)
	}
	if err := os.WriteFile(path.Join(dir, KeyFileName), keyPEM, 0o600); err != nil {
		return fmt.Errorf("writeFiles: write key: %w", err)
	}
	if err := writeAccessKey(path.Join(dir, AccessFileName), accessKey); err != nil {
		return fmt.Errorf("writeFiles: write access key: %w", err)
	}
	return nil
}
