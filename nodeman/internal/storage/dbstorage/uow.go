package dbstorage

import (
	"database/sql"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
)

type uowctx struct {
	tx *sql.Tx
}

var _ users.UoWContext = (*uowctx)(nil)
var _ nodes.UoWContext = (*uowctx)(nil)
var _ subscr.UoWContext = (*uowctx)(nil)
var _ auth.UoWContext = (*uowctx)(nil)
var _ poolsync.UoWContext = (*uowctx)(nil)
