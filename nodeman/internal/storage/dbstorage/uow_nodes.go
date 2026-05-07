package dbstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// service.UoWContext::Storage impl
func (uow *uowctx) NewNode(ctx context.Context, node *models.Node) error {
	query := queryReplacer.Replace(`
		INSERT INTO {nodes} (
			{client_config_template},
			{node_endpoint},
			{node_access_key},
			{node_current_status},
			{node_target_status}
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING {node_id}
	`)

	clientConfigTemplate := ClientConfigTemplate(node.Config.ClientConfigTemplate)

	err := uow.tx.QueryRowContext(ctx, query,
		clientConfigTemplate,
		node.Config.ConnectionInfo.Endpoint,
		node.Config.ConnectionInfo.AccessKey,
		node.CurrentStatus,
		node.TargetStatus,
	).Scan(&node.ID)

	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func (uow *uowctx) GetNode(ctx context.Context, id models.NodeID) (*models.Node, bool, error) {
	query := queryReplacer.Replace(`
		SELECT
			{node_id},
			{client_config_template},
			{node_endpoint},
			{node_access_key},
			{node_current_status},
			{node_target_status}
		FROM {nodes}
		WHERE {node_id} = $1
		  AND {deleted_at} IS NULL
	`)

	var n models.Node
	var clientConfigTemplate ClientConfigTemplate

	err := uow.tx.QueryRowContext(ctx, query, id).Scan(
		&n.ID,
		&clientConfigTemplate,
		&n.Config.ConnectionInfo.Endpoint,
		&n.Config.ConnectionInfo.AccessKey,
		&n.CurrentStatus,
		&n.TargetStatus,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// нет записи
			return nil, false, nil
		}
		return nil, false, xerr.WrapWithStack(err)
	}

	n.Config.ClientConfigTemplate = models.ClientConfigTemplate(clientConfigTemplate)
	return &n, true, nil
}

func (uow *uowctx) ListNodes(ctx context.Context) ([]models.Node, error) {
	query := queryReplacer.Replace(`
		SELECT
			{node_id},
			{client_config_template},
			{node_endpoint},
			{node_access_key},
			{node_current_status},
			{node_target_status}
		FROM {nodes}
		WHERE {deleted_at} IS NULL
		ORDER BY {node_id} ASC
	`)

	rows, err := uow.tx.QueryContext(ctx, query)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	defer rows.Close()

	var nodes []models.Node
	for rows.Next() {
		var n models.Node
		var clientConfigTemplate ClientConfigTemplate
		err := rows.Scan(
			&n.ID,
			&clientConfigTemplate,
			&n.Config.ConnectionInfo.Endpoint,
			&n.Config.ConnectionInfo.AccessKey,
			&n.CurrentStatus,
			&n.TargetStatus,
		)
		if err != nil {
			return nil, xerr.WrapWithStack(err)
		}
		n.Config.ClientConfigTemplate = models.ClientConfigTemplate(clientConfigTemplate)
		nodes = append(nodes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	return nodes, nil
}

func (uow *uowctx) SetTargetNodeStatus(ctx context.Context, id models.NodeID, status models.NodeStatus) error {
	query := queryReplacer.Replace(`
		UPDATE {nodes}
		SET
			{node_target_status} = $1,
			{updated_at} = now()
		WHERE {node_id} = $2
		  AND {deleted_at} IS NULL
	`)

	_, err := uow.tx.ExecContext(ctx, query, status, id)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func (uow *uowctx) SetCurrentNodeStatus(ctx context.Context, id models.NodeID, status models.NodeStatus) error {
	query := queryReplacer.Replace(`
		UPDATE {nodes}
		SET
			{node_current_status} = $1,
			{updated_at} = now()
		WHERE {node_id} = $2
		  AND {deleted_at} IS NULL
	`)

	_, err := uow.tx.ExecContext(ctx, query, status, id)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}
func (uow *uowctx) SetClientConfig(ctx context.Context, id models.NodeID, cfg models.ClientConfigTemplate) error {
	query := queryReplacer.Replace(`
		UPDATE {nodes}
		SET
			{client_config_template} = $1,
			{updated_at} = now()
		WHERE {node_id} = $2
		  AND {deleted_at} IS NULL
	`)

	dbCfg := ClientConfigTemplate(cfg)

	_, err := uow.tx.ExecContext(ctx, query, dbCfg, id)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func (uow *uowctx) DeleteNode(ctx context.Context, id models.NodeID) error {
	query := queryReplacer.Replace(`
		UPDATE {nodes}
		SET {deleted_at} = now()
		WHERE {node_id} = $1
		  AND {deleted_at} IS NULL
	`)

	_, err := uow.tx.ExecContext(ctx, query, id)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}
