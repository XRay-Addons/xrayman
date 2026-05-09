package memstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// UsersStorage impl
func (s *Storage) NewUser(ctx context.Context,
	user *models.User,
) error {
	user.Profile.ID = models.UserID(len(s.users))
	s.users = append(s.users, *user)
	for nodeID := range s.syncStatus {
		s.syncStatus[nodeID] = append(s.syncStatus[nodeID], models.UserStatusDisabled)
	}
	return nil
}

func (s *Storage) SetTargetUserStatus(ctx context.Context,
	id models.UserID, status models.UserStatus,
) error {
	s.users[id].TargetStatus = status
	return nil
}

func (s *Storage) ListUsers(ctx context.Context) (
	[]models.User, error,
) {
	var users []models.User
	users = append(users, s.users...)
	return users, nil
}

func (s *Storage) GetUser(ctx context.Context,
	id models.UserID,
) (*models.User, bool, error) {
	return &s.users[id], true, nil
}

func (s *Storage) DeleteUser(ctx context.Context,
	id models.UserID,
) error {
	s.users[id].TargetStatus = 0
	return nil
}
