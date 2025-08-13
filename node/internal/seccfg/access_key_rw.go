package seccfg

import (
	"encoding/json"
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type accessKeyJSON struct {
	AccessKey string `json:"accessKey"`
}

func readAccessKey(file string) (models.AccessKey, error) {
	var key models.AccessKey

	data, err := os.ReadFile(file)
	if err != nil {
		return key, err
	}

	var container accessKeyJSON
	if err := json.Unmarshal(data, &container); err != nil {
		return key, err
	}

	err = key.UnmarshalText([]byte(container.AccessKey))
	if err != nil {
		return key, err
	}

	return key, nil
}

func writeAccessKey(file string, key models.AccessKey) error {
	text, err := key.MarshalText()
	if err != nil {
		return err
	}

	container := accessKeyJSON{
		AccessKey: string(text),
	}

	data, err := json.MarshalIndent(container, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0o600)
}
