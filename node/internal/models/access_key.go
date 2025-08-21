package models

import (
	"encoding/base64"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

type CertHash = [32]byte
type AccessSecret = [32]byte

type AccessKey struct {
	CertHash     CertHash
	AccessSecret AccessSecret
}

func (k *AccessKey) MarshalText() ([]byte, error) {
	data := append(k.CertHash[:], k.AccessSecret[:]...)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encoded, data)
	return encoded, nil
}

func (k *AccessKey) UnmarshalText(text []byte) error {
	raw := make([]byte, base64.StdEncoding.DecodedLen(len(text)))
	n, err := base64.StdEncoding.Decode(raw, text)
	if err != nil {
		return errdefs.WithStack(err)
	}
	if n != len(k.CertHash)+len(k.AccessSecret) {
		return errdefs.New("access key: invalid length")
	}
	copy(k.CertHash[:], raw[:len(k.CertHash)])
	copy(k.AccessSecret[:], raw[len(k.CertHash):])
	return nil
}

func (k AccessKey) String() string {
	data := append(k.CertHash[:], k.AccessSecret[:]...)
	return base64.StdEncoding.EncodeToString(data)
}
