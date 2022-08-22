package services

import (
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
)

type permissionsService struct {
	PermissionsRepository data.PermissionsRepository
	Logger                *jsonlog.Logger
}

func NewPermissionsService(permissionRepo data.PermissionsRepository, logger *jsonlog.Logger) PermissionsService {
	return &permissionsService{permissionRepo, logger}
}

func (s *permissionsService) GetAllForUser(userId int64) (data.Permissions, error) {
	return s.PermissionsRepository.GetAllForUser(userId)
}
