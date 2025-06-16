package xrayapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	handlerService "github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/vless"
)

func AddUsers(
	ctx context.Context,
	hs handlerService.HandlerServiceClient,
	ins []models.Inbound,
	users []models.User,
) error {
	if hs == nil {
		return fmt.Errorf("%w: api not exists", errdefs.ErrIPE)
	}

	// execute in tx
	var tx ApiTx
	for _, in := range ins {
		for _, u := range users {
			tx.AddFn(TxFn{
				Fn: func() error { return addUser(ctx, hs, in, u) },
				Rb: func() error { return delUser(ctx, hs, in, u) },
			})
		}
	}
	if err := tx.Execute(); err != nil {
		return fmt.Errorf("%w: add users: %v", errdefs.ErrXRay, err)
	}

	return nil
}

func DelUsers(
	ctx context.Context,
	hs handlerService.HandlerServiceClient,
	ins []models.Inbound,
	users []models.User,
) error {
	if hs == nil {
		return fmt.Errorf("%w: api not exists", errdefs.ErrIPE)
	}

	// execute in tx
	var tx ApiTx
	for _, in := range ins {
		for _, u := range users {
			tx.AddFn(TxFn{
				Fn: func() error { return delUser(ctx, hs, in, u) },
				Rb: func() error { return addUser(ctx, hs, in, u) },
			})
		}
	}
	if err := tx.Execute(); err != nil {
		return fmt.Errorf("%w: del users: %v", errdefs.ErrXRay, err)
	}

	return nil
}

func addUser(
	ctx context.Context,
	hs handlerService.HandlerServiceClient,
	in models.Inbound,
	u models.User,
) error {
	addUserOp := func(u *protocol.User) *serial.TypedMessage {
		return serial.ToTypedMessage(&handlerService.AddUserOperation{User: u})
	}
	err := alterUser(ctx, hs, in, u, addUserOp)

	// already exists is not an error for us
	alreadyExistsErrPattern := fmt.Sprintf("User %s already exists", u.Name)
	if err != nil && strings.Contains(err.Error(), alreadyExistsErrPattern) {
		return nil
	}

	return fmt.Errorf("add user: %w", err)
}

func delUser(
	ctx context.Context,
	hs handlerService.HandlerServiceClient,
	in models.Inbound,
	u models.User,
) error {
	delUserOp := func(u *protocol.User) *serial.TypedMessage {
		return serial.ToTypedMessage(&handlerService.RemoveUserOperation{Email: u.Email})
	}

	err := alterUser(ctx, hs, in, u, delUserOp)

	// not exists is not an error for us
	notFoundErrPattern := fmt.Sprintf("User %s not found", u.Name)
	if err != nil && strings.Contains(err.Error(), notFoundErrPattern) {
		return nil
	}

	return fmt.Errorf("del user: %w", err)
}

type ApiUserOp = func(*protocol.User) *serial.TypedMessage

func alterUser(
	ctx context.Context,
	hs handlerService.HandlerServiceClient,
	in models.Inbound,
	u models.User,
	op ApiUserOp,
) error {
	config, err := toProtocolUser(in, u)
	if err != nil {
		return err
	}
	_, err = hs.AlterInbound(ctx, &handlerService.AlterInboundRequest{
		Tag:       in.Tag,
		Operation: op(config),
	})

	return err
}

func toProtocolUser(in models.Inbound, u models.User) (*protocol.User, error) {
	var account *vless.Account
	switch in.Type {
	case models.VlessTcpReality:
		account = &vless.Account{
			Id:         u.UUID,
			Encryption: "none",
			Flow:       "xtls-rprx-vision",
		}
	case models.VlessXHTTP:
		account = &vless.Account{
			Id:         u.UUID,
			Encryption: "none",
		}
	default:
		return nil, fmt.Errorf("%w: unsupported inbound", errdefs.ErrIPE)
	}

	return &protocol.User{
		Email:   u.Name,
		Account: serial.ToTypedMessage(account),
	}, nil
}
