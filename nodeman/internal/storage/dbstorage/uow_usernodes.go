package dbstorage

import (
	"context"
	"fmt"
	"strings"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (uow *uowctx) FindPendingSyncs(ctx context.Context, id models.NodeID) ([]models.UserSyncStatus, error) {
	query := queryReplacer.Replace(`
		SELECT
			u.{user_id},
			u.{user_name},
			u.{display_name},
			u.{vless_uuid},
			u.{user_target_status},
			COALESCE(s.{user_current_status}, {user_status_disabled}) AS current_status
		FROM {users} u
		LEFT JOIN {syncs} s
			ON s.{user_id} = u.{user_id} AND s.{node_id} = $1
		WHERE u.{user_target_status} <> COALESCE(s.{user_current_status}, {user_status_disabled})
	`)

	rows, err := uow.tx.QueryContext(ctx, query, id)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	defer rows.Close()

	var result []models.UserSyncStatus
	for rows.Next() {
		var us models.UserSyncStatus
		err := rows.Scan(
			&us.User.Profile.ID,
			&us.User.Profile.Name,
			&us.User.Profile.DisplayName,
			&us.User.Profile.VlessUUID,
			&us.User.TargetStatus,
			&us.CurrentStatus,
		)
		if err != nil {
			return nil, xerr.WrapWithStack(err)
		}
		result = append(result, us)
	}

	if err := rows.Err(); err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	return result, nil
}

func (uow *uowctx) SetNodeUsers(ctx context.Context, id models.NodeID, patch []models.UserStatusPatch) error {
	// remove old users
	// TODO: 3-steps via temp table
	delQuery := queryReplacer.Replace(`
		DELETE FROM {syncs}
		WHERE {node_id} = $1
	`)
	if _, err := uow.tx.ExecContext(ctx, delQuery, id); err != nil {
		return xerr.WrapWithStack(err)
	}

	if len(patch) == 0 {
		return nil
	}

	// bulk insert new users
	insertQuery := queryReplacer.Replace(`
		INSERT INTO {syncs} ({user_id}, {node_id}, {user_current_status})
		VALUES %s
	`)

	values := make([]interface{}, 0, len(patch)*3)
	placeholders := make([]string, 0, len(patch))
	for i, p := range patch {
		start := i*3 + 1
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d)", start, start+1, start+2))
		values = append(values, p.UserID, id, p.Status)
	}

	q := fmt.Sprintf(insertQuery, strings.Join(placeholders, ","))
	if _, err := uow.tx.ExecContext(ctx, q, values...); err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func (uow *uowctx) UpdateNodeUsers(ctx context.Context, id models.NodeID, patch []models.UserStatusPatch) error {
	if len(patch) == 0 {
		return nil
	}

	insertQuery := queryReplacer.Replace(`
		INSERT INTO {syncs} ({user_id}, {node_id}, {user_current_status})
		VALUES %s
		ON CONFLICT ({user_id}, {node_id})
		DO UPDATE SET {user_current_status} = EXCLUDED.{user_current_status}
	`)

	values := make([]interface{}, 0, len(patch)*3)
	placeholders := make([]string, 0, len(patch))
	for i, p := range patch {
		start := i*3 + 1
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d)", start, start+1, start+2))
		values = append(values, p.UserID, id, p.Status)
	}

	q := fmt.Sprintf(insertQuery, strings.Join(placeholders, ","))
	if _, err := uow.tx.ExecContext(ctx, q, values...); err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func (uow *uowctx) GetUserNodes(ctx context.Context, id models.UserID) ([]models.Node, error) {
	query := queryReplacer.Replace(`
		SELECT
			n.{node_id},
			n.{client_config_template},
			n.{node_endpoint},
			n.{node_access_key},
			n.{node_current_status},
			n.{node_target_status}
		FROM {nodes} n
		INNER JOIN {syncs} s
			ON s.{node_id} = n.{node_id}
		WHERE s.{user_id} = $1
		  AND s.{user_current_status} = {user_status_enabled}
		  AND n.{node_target_status} = {node_status_running}
		  AND n.{deleted_at} IS NULL
	`)

	rows, err := uow.tx.QueryContext(ctx, query, id)
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
