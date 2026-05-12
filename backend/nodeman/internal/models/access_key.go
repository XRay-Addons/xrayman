package models

import (
	"database/sql/driver"
	"encoding/base64"

	"github.com/XRay-Addons/xrayman/common/xerr"
)

type CertHash = [32]byte
type AccessSecret = [32]byte

type AccessKey struct {
	CertHash     CertHash
	AccessSecret AccessSecret
}

func (k *AccessKey) MarshalText() ([]byte, error) {
	data := k.getKeyData()
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encoded, data)
	return encoded, nil
}

func (k *AccessKey) UnmarshalText(text []byte) error {
	raw := make([]byte, base64.StdEncoding.DecodedLen(len(text)))
	n, err := base64.StdEncoding.Decode(raw, text)
	if err != nil {
		return xerr.WrapWithStack(err)
	}
	if n != len(k.CertHash)+len(k.AccessSecret) {
		return xerr.New("invalid access key length")
	}
	copy(k.CertHash[:], raw[:len(k.CertHash)])
	copy(k.AccessSecret[:], raw[len(k.CertHash):])
	return nil
}

func (k AccessKey) String() string {
	return base64.StdEncoding.EncodeToString(k.getKeyData())
}

func (k AccessKey) getKeyData() []byte {
	data := make([]byte, len(k.CertHash)+len(k.AccessSecret))
	copy(data[:len(k.CertHash)], k.CertHash[:])
	copy(data[len(k.CertHash):], k.AccessSecret[:])
	return data
}

func (k AccessKey) Value() (driver.Value, error) {
	return k.getKeyData(), nil
}

func (k *AccessKey) Scan(src any) error {
	if src == nil {
		*k = AccessKey{}
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return xerr.New("invalid type for AccessKey")
	}
	if len(b) != len(CertHash{})+len(AccessSecret{}) {
		return xerr.New("invalid length for AccessKey")
	}
	copy(k.CertHash[:], b[:len(k.CertHash)])
	copy(k.AccessSecret[:], b[len(k.CertHash):])
	return nil
}
