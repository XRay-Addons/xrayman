package secrets

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func ensureTLS(dir string, iss string, exp time.Duration) (cert, key []byte, err error) {
	if cert, key, err = readTLS(dir); err == nil {
		return cert, key, nil
	}

	cert, key, err = generateTLS(iss, exp)
	if err != nil {
		return nil, nil, fmt.Errorf("ensure tls: %w", err)
	}

	if err := writeTLS(dir, cert, key); err != nil {
		return nil, nil, fmt.Errorf("ensure tls: %w", err)
	}

	return cert, key, nil
}

func readTLS(dir string) (cert, key []byte, err error) {
	certPath := filepath.Join(dir, CertFile)
	keyPath := filepath.Join(dir, KeyFile)

	cert, err = os.ReadFile(certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read tls: %w", err)
	}

	key, err = os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read tls: %w", err)
	}

	if _, err := tls.X509KeyPair(cert, key); err != nil {
		return nil, nil, fmt.Errorf("read tls: %w", err)
	}

	return cert, key, nil
}

func writeTLS(dir string, cert, key []byte) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("write tls: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir, CertFile), cert, 0o600); err != nil {
		return fmt.Errorf("write tls: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, KeyFile), key, 0o600); err != nil {
		return fmt.Errorf("write tls: %w", err)
	}

	return nil
}

func generateTLS(issuer string, exp time.Duration) (certPEM, keyPEM []byte, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("generate tls: %w", err)
	}

	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return nil, nil, fmt.Errorf("generate tls: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: issuer,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(exp),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("generate tls: %w", err)
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return certPEM, keyPEM, nil
}
