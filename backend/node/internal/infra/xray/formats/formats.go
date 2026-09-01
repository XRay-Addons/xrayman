package formats

import "github.com/XRay-Addons/xrayman/node/internal/models"

func InboundFormats() []models.InboundFormat {
	return []models.InboundFormat{
		&VlessTCPReality{},
		&VlessXHTTP{},
	}
}
