package converter

import (
	"encoding/base64"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
)

// goverter:converter
// goverter:output:format function
// goverter:output:file ./converter_generated.go
// goverter:extend ConvertNodeID RConvertNodeID ConvertNodeStatusResult ConvertAccessSecret ConvertCertHash
//
//go:generate goverter gen .
type Converter interface {
	ConvertNewNodeRequest(r *api.NewNodeRequest) (*models.NewNodeParams, error)
	ConvertNewNodeResult(r *models.NewNodeResult) *api.NewNodeResponse

	ConvertStartNodeRequest(r *api.StartNodeRequest) (*models.StartNodeParams, error)
	ConvertStopNodeRequest(r *api.StopNodeRequest) (*models.StopNodeParams, error)

	ConvertListNodesResult(r *models.ListNodeResult) *api.ListNodeResponse
}

// goverter:extend
func ConvertNodeID(i models.NodeID) api.NodeID {
	return api.NodeID(i)
}

// goverter:extend
func RConvertNodeID(i api.NodeID) models.NodeID {
	return models.NodeID(i)
}

// goverter:extend
func ConvertAccessSecret(s []byte) (models.AccessSecret, error) {
	var secret models.AccessSecret

	decoded, err := base64.StdEncoding.DecodeString(string(s))
	if err != nil {
		return secret, fmt.Errorf("base64 decode error: %w", err)
	}

	if len(decoded) != len(secret) {
		return secret, fmt.Errorf("invalid length: expected %d, got %d", len(secret), len(decoded))
	}

	copy(secret[:], decoded)
	return secret, nil
}

// goverter:extend
func ConvertCertHash(h []byte) (models.CertHash, error) {
	var hash models.CertHash

	decoded, err := base64.StdEncoding.DecodeString(string(h))
	if err != nil {
		return hash, fmt.Errorf("base64 decode error: %w", err)
	}

	if len(decoded) != len(hash) {
		return hash, fmt.Errorf("invalid length: expected %d, got %d", len(hash), len(decoded))
	}

	copy(hash[:], decoded)
	return hash, nil
}

// goverter:extend
func ConvertNodeStatusResult(source models.NodeStatus) api.NodeStatus {
	var response api.NodeStatus
	switch source {
	case models.NodeStatusStopped:
		response = api.NodeStatusStopped
	case models.NodeStatusRunning:
		response = api.NodeStatusRunning
	case models.NodeStatusUnknown:
		response = api.NodeStatusUnknown
	default:
		panic(fmt.Sprintf("unexpected enum element: %v", source))
	}
	return response
}
