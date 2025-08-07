package models

type NewNodeParams struct {
	Endpoint string
}

type NewNodeResult struct {
	ID           NodeID
	Endpoint     string
	AccessSecret []byte
}

type StartNodeParams struct {
	ID NodeID
}

type StartNodeResult struct {
}

type StopNodeParams struct {
	ID NodeID
}

type StopNodeResult struct {
}

type ListNodeParams struct {
}

type ListNodeResult struct {
	Nodes []Node
}

/*type NewUserRequest struct {
	Name string
}

type NewUserResponse struct {
	User            User
	SubscriptionURL string
}

type EnableUserRequest struct {
	ID UserID
}

type EnableUserResponse struct {
}

type DisableUserRequest struct {
	ID UserID
}

type DisableUserResponse struct {
}*/
