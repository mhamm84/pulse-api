package services

import (
	"context"
	"github.com/mhamm84/pulse-api/internal/data"
)

type permissionsService struct {
	PermissionsRepository data.PermissionsRepository
}

func NewPermissionsService(permissionRepo data.PermissionsRepository) PermissionsService {
	return &permissionsService{permissionRepo}
}

func (s *permissionsService) GetAllForUser(ctx context.Context, userId int64) (data.Permissions, error) {
	return s.PermissionsRepository.GetAllForUser(userId)
}
