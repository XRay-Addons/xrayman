package node

import (
	"net/http"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type HTTPClientFactory interface {
	GetNodeClient(certHash models.CertHash) (*http.Client, error)
}
