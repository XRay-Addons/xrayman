package models

import (
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
	realLen, err := base64.StdEncoding.Decode(raw, text)
	if err != nil {
		return xerr.WrapWithStack(err)
	}
	return k.setKeyData(raw[:realLen])
}

func (k AccessKey) String() string {
	return base64.StdEncoding.EncodeToString(k.getKeyData())
}

func (k AccessKey) Value() ([]byte, error) {
	return k.getKeyData(), nil
}

func (k *AccessKey) Scan(src []byte) error {
	return k.setKeyData(src)
}

func (k AccessKey) getKeyData() []byte {
	data := make([]byte, len(k.CertHash)+len(k.AccessSecret))
	copy(data[:len(k.CertHash)], k.CertHash[:])
	copy(data[len(k.CertHash):], k.AccessSecret[:])
	return data
}

func (k *AccessKey) setKeyData(src []byte) error {
	if len(src) != len(CertHash{})+len(AccessSecret{}) {
		return xerr.Newf("invalid length %d for AccessKey", len(src))
	}
	copy(k.CertHash[:], src[:len(k.CertHash)])
	copy(k.AccessSecret[:], src[len(k.CertHash):])
	return nil
}
