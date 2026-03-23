package dbstorage

import (
	"database/sql"

	"github.com/XRay-Addons/xrayman/nodeman/internal/poolsyncer"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service"
	"github.com/XRay-Addons/xrayman/nodeman/internal/subscrman"
)

type uowctx struct {
	tx *sql.Tx
}

var _ service.UoWContext = (*uowctx)(nil)
var _ poolsyncer.UoWContext = (*uowctx)(nil)
var _ subscrman.UoWContext = (*uowctx)(nil)
