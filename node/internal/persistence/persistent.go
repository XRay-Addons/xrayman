package persistence

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

const (
	CertFileName   = "cert.pem"
	KeyFileName    = "key.pem"
	ConfigFileName = "config.json"
	RSAKeySize     = 2048
	JWTSecretSize  = 32
)

type Config struct {
	JWTSecret string `json:"jwtSecret"`
	CertHash  string `json:"certHash"`
}

type Persistent struct {
	dir string
	cfg Config
}

func New(dir string, log *zap.Logger) (*Persistent, error) {
	if dir == "" {
		return nil, fmt.Errorf("persistence: empty directory path")
	}

	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0700); err != nil {
				return nil, fmt.Errorf("persistence: create dir: %w", err)
			}
		} else {
			return nil, fmt.Errorf("persistence: stat dir: %w", err)
		}
	} else if !info.IsDir() {
		return nil, fmt.Errorf("persistence: path is not a directory")
	}

	p := &Persistent{dir: dir}

	certPath := p.CertPath()
	keyPath := p.KeyPath()
	configPath := p.ConfigPath()

	certExists := fileExists(certPath)
	keyExists := fileExists(keyPath)
	configExists := fileExists(configPath)

	if !certExists || !keyExists || !configExists {
		if err := p.generateAndSave(); err != nil {
			return nil, fmt.Errorf("persistence: generate and save: %w", err)
		}
	} else {
		if err := p.loadConfig(); err != nil {
			return nil, fmt.Errorf("persistence: load config: %w", err)
		}
	}

	if log != nil {
		log.Info("node config",
			zap.String("JWTSecret", p.JWTSecret()),
			zap.String("CertHash", p.CertHash()))
	}

	return p, nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func (p *Persistent) generateAndSave() error {
	privKey, err := rsa.GenerateKey(rand.Reader, RSAKeySize)
	if err != nil {
		return fmt.Errorf("generate key: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Node TLS Cert",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privKey.PublicKey, privKey)
	if err != nil {
		return fmt.Errorf("create certificate: %w", err)
	}

	if err := writePEMFile(p.KeyPath(), "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privKey)); err != nil {
		return fmt.Errorf("write key: %w", err)
	}

	if err := writePEMFile(p.CertPath(), "CERTIFICATE", certDER); err != nil {
		return fmt.Errorf("write cert: %w", err)
	}

	secretBytes := make([]byte, JWTSecretSize)
	if _, err := rand.Read(secretBytes); err != nil {
		return fmt.Errorf("generate jwt secret: %w", err)
	}
	jwtSecret := base64.StdEncoding.EncodeToString(secretBytes)

	certHashBytes := sha256.Sum256(certDER)
	certHash := base64.StdEncoding.EncodeToString(certHashBytes[:])

	p.cfg = Config{
		JWTSecret: jwtSecret,
		CertHash:  certHash,
	}

	cfgData, err := json.MarshalIndent(p.cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config json: %w", err)
	}

	if err := os.WriteFile(p.ConfigPath(), cfgData, 0600); err != nil {
		return fmt.Errorf("write config json: %w", err)
	}

	return nil
}

func writePEMFile(path, blockType string, bytes []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return pem.Encode(f, &pem.Block{
		Type:  blockType,
		Bytes: bytes,
	})
}

func (p *Persistent) loadConfig() error {
	data, err := os.ReadFile(p.ConfigPath())
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &p.cfg)
}

func (p *Persistent) CertPath() string {
	return filepath.Join(p.dir, CertFileName)
}

func (p *Persistent) KeyPath() string {
	return filepath.Join(p.dir, KeyFileName)
}

func (p *Persistent) ConfigPath() string {
	return filepath.Join(p.dir, ConfigFileName)
}

func (p *Persistent) JWTSecret() string {
	return p.cfg.JWTSecret
}

func (p *Persistent) CertHash() string {
	return p.cfg.CertHash
}
