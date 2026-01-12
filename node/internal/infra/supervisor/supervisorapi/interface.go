package supervisorapi

import "context"

type Supervisor interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Status(ctx context.Context) (ServiceStatus, error)
	Close(ctx context.Context) error
}
