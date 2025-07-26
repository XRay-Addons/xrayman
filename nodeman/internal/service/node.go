package service

// user description
/*type UserID int

type User struct {
	ID        UserID
	Name      string
	VlessUUID string
}

type UserStatus int

const (
	UserStatusUnknown UserStatus = iota
	UserDisabled
	UserEnabled
)

func (s UserStatus) String() string {
	switch s {
	case UserStatusUnknown:
		return "Unknown"
	case UserDisabled:
		return "Off"
	case UserEnabled:
		return "On"
	default:
		return "Unknown"
	}
}

// node description

type NodeStatus int

const (
	NodeStatusUnknown NodeStatus = iota
	NodeStopped
	NodeRunning
)

func (s NodeStatus) String() string {
	switch s {
	case NodeStatusUnknown:
		return "Unknow"
	case NodeStopped:
		return "Stopped"
	case NodeRunning:
		return "Running"
	default:
		return "Unknown"
	}
}

type UserState struct {
	User   User
	Status UserStatus
}

type NodeConfig struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mocks/node_api.go -package=mocks
type NodeAPI interface {
	Start(ctx context.Context, users []User) (*NodeConfig, error)
	Stop(ctx context.Context) error
	Status(ctx context.Context) (NodeStatus, error)
	EditUsers(ctx context.Context, users []UserState) error
	Close(ctx context.Context) error
}

type NodeUpdater interface {
	SetConfig(cfg *NodeConfig)
	SetStatus(state NodeStatus)
	SetUsers(users []UserState)
	Apply(ctx context.Context) error
}

type NodeStorage interface {
	GetNodeState(ctx context.Context) (actual, required NodeStatus, err error)
	GetOutOfSyncUsers(ctx context.Context) ([]UserState, error)
	GetAllUsers(ctx context.Context) ([]UserState, error)

	GetUpdater(ctx context.Context) (NodeUpdater, error)
}

type NodeController struct {
	nodeAPI      NodeAPI
	stateStorage NodeStorage
	log          *zap.SugaredLogger
}

func New(api NodeAPI, storage NodeStorage, log *zap.Logger) (*NodeController, error) {
	if api == nil {
		return nil, fmt.Errorf("node init: api: %w", errdefs.ErrNilArgPassed)
	}
	if storage == nil {
		return nil, fmt.Errorf("node init: storage: %w", errdefs.ErrNilArgPassed)
	}
	if log == nil {
		return nil, fmt.Errorf("node init: logger: %w", errdefs.ErrNilArgPassed)
	}

	return &NodeController{
		nodeAPI:      api,
		stateStorage: storage,
		log:          log.Sugar(),
	}, nil

}

func (c *NodeController) Close(ctx context.Context) error {
	if c == nil || c.nodeAPI == nil {
		return nil
	}
	if err := c.nodeAPI.Close(ctx); err != nil {
		return fmt.Errorf("node: close: %w", err)
	}
	return nil
}

func (c *NodeController) SyncNodeStatus(ctx context.Context) (err error) {
	if c == nil || c.nodeAPI == nil || c.stateStorage == nil {
		return fmt.Errorf("node: sync: %w", errdefs.ErrNilObjectCall)
	}

	c.log.Info("sync node")

	// get current node state from storage
	prevState, targetState, err := c.stateStorage.GetNodeState(ctx)
	if err != nil {
		return fmt.Errorf("node: sync: %w", err)
	}
	c.log.Infof("sync node: prev: %v, target: %v", prevState, targetState)

	// get node status if required
	currState := prevState
	if targetState == NodeRunning && prevState != NodeStopped {
		if currState, err = c.nodeAPI.Status(ctx); err != nil {
			return fmt.Errorf("node: sync: check status: %w", err)
		}
		c.log.Infof("sync node: update curr: %v", currState)
	}

	// update node status in storage if changed
	if currState != prevState {
		if err := c.changeNodeState(ctx, currState); err != nil {
			return fmt.Errorf("sync node: %w", err)
		}
	}

	// start, stop or edit node users
	switch {
	case targetState == NodeRunning && currState == NodeStopped:
		err = c.startNode(ctx)
	case targetState == NodeRunning && currState == NodeRunning:
		err = c.syncNodeUsers(ctx)
	case targetState == NodeStopped && currState == NodeRunning:
		err = c.stopNode(ctx)
	}

	if err != nil {
		return fmt.Errorf("node: sync: %w", err)
	}
	c.log.Info("sync node: state changed OK")

	return nil
}

func (c *NodeController) changeNodeState(ctx context.Context, state NodeStatus) error {
	upd, err := c.stateStorage.GetUpdater(ctx)
	if err != nil {
		return fmt.Errorf("set node unavailable: %w", err)
	}
	upd.SetStatus(state)
	if err = upd.Apply(ctx); err != nil {
		return fmt.Errorf("set node unavailable: %w", err)
	}
	return nil
}

func (c *NodeController) startNode(ctx context.Context) error {
	// get list of all users
	allUsers, err := c.stateStorage.GetAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("start node: get all users: %w", err)
	}

	// select only enabled users,
	enabledUsers := make([]User, 0, len(allUsers))
	for _, u := range allUsers {
		if u.Status != UserEnabled {
			continue
		}
		enabledUsers = append(enabledUsers, u.User)
	}

	if err := c.lockNodeState(ctx); err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	// start node
	cfg, err := c.nodeAPI.Start(ctx, enabledUsers)
	if err != nil {
		c.unlockNodeState(ctx, NodeStopped)
		return fmt.Errorf("start node: start via api: %w", err)
	}

	// update node state
	upd, err := c.stateStorage.GetUpdater(ctx)
	if err != nil {
		c.unlockNodeState(ctx, NodeStopped)
		return fmt.Errorf("start node: update state: %w", err)
	}
	upd.SetStatus(NodeRunning)
	upd.SetUsers(allUsers)
	upd.SetConfig(cfg)
	if err = upd.Apply(ctx); err != nil {
		return fmt.Errorf("start node: update state: %w", err)
	}

	return nil
}

func (c *NodeController) stopNode(ctx context.Context) error {
	// set node status to stopped
	if err := c.lockNodeState(ctx); err != nil {
		return fmt.Errorf("stop node: %w", err)
	}

	// stop node
	if err := c.nodeAPI.Stop(ctx); err != nil {
		c.unlockNodeState(ctx, NodeRunning)
		return fmt.Errorf("stop node: stop via api: %w", err)
	}
	c.unlockNodeState(ctx, NodeStopped)

	return nil
}

func (c *NodeController) syncNodeUsers(ctx context.Context) error {
	oosUsers, err := c.stateStorage.GetOutOfSyncUsers(ctx)
	if err != nil {
		return fmt.Errorf("edit node users: %w", err)
	}

	if err := c.lockNodeUsersState(ctx, oosUsers); err != nil {
		return fmt.Errorf("sync node users: %w", err)
	}
	if err := c.nodeAPI.EditUsers(ctx, oosUsers); err != nil {
		return fmt.Errorf("edit node users: %w", err)
	}
	upd, err := c.stateStorage.GetUpdater(ctx)
	if err != nil {
		return fmt.Errorf("sync node users: update synced users: %w", err)
	}
	upd.SetUsers(oosUsers)
	if err := upd.Apply(ctx); err != nil {
		return fmt.Errorf("sync node users: update synced users: %w", err)
	}
	return nil
}

func (c *NodeController) lockNodeState(ctx context.Context) error {
	upd, err := c.stateStorage.GetUpdater(ctx)
	if err != nil {
		return fmt.Errorf("lock node state: %w", err)
	}
	upd.SetStatus(NodeStatusUnknown)
	if err := upd.Apply(ctx); err != nil {
		return fmt.Errorf("lock node state: %w", err)
	}
	return nil
}

func (c *NodeController) unlockNodeState(ctx context.Context, state NodeStatus) error {
	upd, err := c.stateStorage.GetUpdater(ctx)
	if err != nil {
		return fmt.Errorf("lock node state: %w", err)
	}
	upd.SetStatus(state)
	if err := upd.Apply(ctx); err != nil {
		return fmt.Errorf("unlock node state: %w", err)
	}
	return nil
}

func (c *NodeController) lockNodeUsersState(ctx context.Context, users []UserState) error {
	upd, err := c.stateStorage.GetUpdater(ctx)
	if err != nil {
		return fmt.Errorf("lock node state: %w", err)
	}

	updUsers := make([]UserState, 0, len(users))
	for _, u := range users {
		updUsers = append(updUsers, UserState{
			User:   u.User,
			Status: UserStatusUnknown,
		})
	}

	upd.SetUsers(updUsers)
	if err := upd.Apply(ctx); err != nil {
		return fmt.Errorf("lock node users state: %w", err)
	}
	return nil
}

func (c *NodeController) unlockNodeUsersState(ctx context.Context, users []UserState) error {
	upd, err := c.stateStorage.GetUpdater(ctx)
	if err != nil {
		return fmt.Errorf("lock node state: %w", err)
	}
	upd.SetUsers(users)
	if err := upd.Apply(ctx); err != nil {
		return fmt.Errorf("unlock node users state: %w", err)
	}
	return nil
}*/
