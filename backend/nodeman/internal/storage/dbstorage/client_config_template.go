package dbstorage

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type ClientConfigTemplate models.ClientConfigTemplate

func (c ClientConfigTemplate) Value() (driver.Value, error) {
	b, err := json.Marshal(models.ClientConfigTemplate(c))
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	return string(b), nil
}

func (c *ClientConfigTemplate) Scan(src any) error {
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return xerr.Newf("unsupported type %T for ClientConfigTemplate", src)
	}
	if err := json.Unmarshal(data, (*models.ClientConfigTemplate)(c)); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}
