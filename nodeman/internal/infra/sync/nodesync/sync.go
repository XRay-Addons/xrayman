package nodesync

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
)

func SyncState(ctx context.Context, client Client, storage Storage) error {
	if client == nil {
		return errdefs.NilArg("client")
	}
	if storage == nil {
		return errdefs.NilArg("storage")
	}
	s := syncer{
		storage: storage,
		client:  client,
	}
	if err := s.SyncNodeState(ctx); err != nil {
		return err
	}
	return nil
}
