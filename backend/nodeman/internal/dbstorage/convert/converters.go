package convert

import (
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func NewNodeReq(r *models.Node) (*queries.NewNodeParams, error) {
	return cnv(r,
		func(from *models.Node, to *queries.NewNodeParams) {
			to.NodeEndpoint = from.Config.ConnectionInfo.Endpoint
			to.NodeCurrentStatus = int16(from.CurrentStatus)
			to.NodeTargetStatus = int16(from.TargetStatus)
			to.Version = from.Config.Settings.Version
		},
		func(from *models.Node, to *queries.NewNodeParams) (err error) {
			to.ClientCfgTemplate, err = from.Config.Settings.ClientConfigTemplate.Value()
			return
		},
		func(from *models.Node, to *queries.NewNodeParams) (err error) {
			to.NodeAccessKey, err = from.Config.ConnectionInfo.AccessKey.Value()
			return
		},
	)
}

func GetNodeResp(r *queries.GetNodeRow) (*models.Node, error) {
	return cnv(r,
		func(from *queries.GetNodeRow, to *models.Node) {
			to.ID = models.NodeID(from.NodeID)
			to.CurrentStatus = models.NodeStatus(from.NodeCurrentStatus)
			to.TargetStatus = models.NodeStatus(from.NodeTargetStatus)
			to.Config.ConnectionInfo.Endpoint = from.NodeEndpoint
			to.Config.Settings.Version = from.Version
		},
		func(from *queries.GetNodeRow, to *models.Node) error {
			return to.Config.ConnectionInfo.AccessKey.Scan(from.NodeAccessKey)
		},
		func(from *queries.GetNodeRow, to *models.Node) error {
			return to.Config.Settings.ClientConfigTemplate.Scan(from.ClientCfgTemplate)
		})
}

func ListNodesResp(r []queries.ListNodesRow) ([]models.Node, error) {
	return cnvArr(r,
		func(from *queries.ListNodesRow, to *models.Node) {
			to.ID = models.NodeID(from.NodeID)
			to.CurrentStatus = models.NodeStatus(from.NodeCurrentStatus)
			to.TargetStatus = models.NodeStatus(from.NodeTargetStatus)
			to.Config.ConnectionInfo.Endpoint = from.NodeEndpoint
			to.Config.Settings.Version = from.Version

		},
		func(from *queries.ListNodesRow, to *models.Node) error {
			return to.Config.ConnectionInfo.AccessKey.Scan(from.NodeAccessKey)
		},
		func(from *queries.ListNodesRow, to *models.Node) error {
			return to.Config.Settings.ClientConfigTemplate.Scan(from.ClientCfgTemplate)
		})
}

func SetNodeSettingsReq(id models.NodeID,
	cfg *models.NodeSettings,
) (*queries.SetNodeSettingsParams, error) {
	tmpl, err := cfg.ClientConfigTemplate.Value()
	if err != nil {
		return nil, err
	}
	return &queries.SetNodeSettingsParams{
		NodeID:            int64(id),
		ClientCfgTemplate: tmpl,
		Version:           cfg.Version,
	}, nil
}

func NewUserReq(r *models.User) *queries.NewUserParams {
	return cnvNoErr(r,
		func(from *models.User, to *queries.NewUserParams) {
			to.DisplayName = from.Profile.DisplayName
			to.UserName = from.Profile.Name
			to.UserTargetStatus = int16(from.TargetStatus)
			to.VlessUuid = from.Profile.VlessUUID
		})
}

func GetUserViewResp(r *queries.GetUserViewRow) *models.UserView {
	return cnvNoErr(r,
		func(from *queries.GetUserViewRow, to *models.UserView) {
			to.User.Profile.ID = models.UserID(from.UserID)
			to.User.Profile.Name = from.UserName
			to.User.Profile.DisplayName = from.DisplayName
			to.User.Profile.VlessUUID = from.VlessUuid
			to.User.TargetStatus = models.UserStatus(from.UserTargetStatus)
			to.Traffic.Total.Download = from.DownloadTotal
			to.Traffic.Total.Upload = from.UploadTotal
			to.Traffic.LastMonth.Download = from.DownloadLastDays
			to.Traffic.LastMonth.Upload = from.UploadLastDays
		})
}

func ListUsersResp(r []queries.ListUsersRow) []models.User {
	return cnvArrNoErr(r,
		func(from *queries.ListUsersRow, to *models.User) {
			to.Profile.ID = models.UserID(from.UserID)
			to.Profile.Name = from.UserName
			to.Profile.DisplayName = from.DisplayName
			to.Profile.VlessUUID = from.VlessUuid
			to.TargetStatus = models.UserStatus(from.UserTargetStatus)
		},
	)
}

func ListUserViewsResp(r []queries.ListUserViewsRow) []models.UserView {
	return cnvArrNoErr(r,
		func(from *queries.ListUserViewsRow, to *models.UserView) {
			to.User.Profile.ID = models.UserID(from.UserID)
			to.User.Profile.Name = from.UserName
			to.User.Profile.DisplayName = from.DisplayName
			to.User.Profile.VlessUUID = from.VlessUuid
			to.User.TargetStatus = models.UserStatus(from.UserTargetStatus)
			to.Traffic.Total.Upload = from.UploadTotal
			to.Traffic.Total.Download = from.DownloadTotal
			to.Traffic.LastMonth.Download = from.DownloadLastDays
			to.Traffic.LastMonth.Upload = from.UploadLastDays
		},
	)
}

func FindPendingSyncsResp(r []queries.FindPendingSyncsRow) []models.UserSyncStatus {
	return cnvArrNoErr(r,
		func(from *queries.FindPendingSyncsRow, to *models.UserSyncStatus) {
			to.CurrentStatus = models.UserStatus(from.UserCurrentStatus)
			to.User.TargetStatus = models.UserStatus(from.UserTargetStatus)
			to.User.Profile.DisplayName = from.DisplayName
			to.User.Profile.ID = models.UserID(from.UserID)
			to.User.Profile.Name = from.UserName
			to.User.Profile.VlessUUID = from.VlessUuid
		},
	)
}

func UpdateNodeUsersReq(id models.NodeID,
	patch []models.UserStatusPatch,
) queries.InsertNodeUsersParams {
	n := len(patch)
	arg := queries.InsertNodeUsersParams{
		NodeID:            int64(id),
		UserID:            make([]int64, n, n),
		UserCurrentStatus: make([]int16, n, n),
	}
	for i, p := range patch {
		arg.UserID[i] = int64(p.UserID)
		arg.UserCurrentStatus[i] = int16(p.Status)
	}
	return arg
}

func GetUserNodesResp(r []queries.GetUserNodesRow) ([]models.Node, error) {
	return cnvArr(r,
		func(from *queries.GetUserNodesRow, to *models.Node) {
			to.ID = models.NodeID(from.NodeID)
			to.CurrentStatus = models.NodeStatus(from.NodeCurrentStatus)
			to.TargetStatus = models.NodeStatus(from.NodeTargetStatus)
			to.Config.ConnectionInfo.Endpoint = from.NodeEndpoint
			to.Config.Settings.Version = from.Version
		},
		func(from *queries.GetUserNodesRow, to *models.Node) (err error) {
			err = to.Config.Settings.ClientConfigTemplate.Scan(from.ClientCfgTemplate)
			return
		},
		func(from *queries.GetUserNodesRow, to *models.Node) (err error) {
			err = to.Config.ConnectionInfo.AccessKey.Scan(from.NodeAccessKey)
			return
		},
	)
}

func UpdateNodeStatsReq(nodeID models.NodeID,
	stats models.NodeStats,
) queries.UpdateTotalStatsParams {
	n := len(stats.Users)
	req := queries.UpdateTotalStatsParams{
		NodeID:   int64(nodeID),
		UserID:   make([]int64, n, n),
		Upload:   make([]int64, n, n),
		Download: make([]int64, n, n),
	}
	for i, u := range stats.Users {
		req.UserID[i] = int64(u.ID)
		req.Upload[i] = u.Uplink
		req.Download[i] = u.Downlink
	}
	return req
}
