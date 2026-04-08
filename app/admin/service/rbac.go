package service

import "github.com/go-zoox/connect/app/admin/repository"

// RBACService resolves effective permission codes for admin users.
type RBACService struct {
	rbac *repository.RBACRepo
}

// NewRBACService builds an RBACService.
func NewRBACService(rbac *repository.RBACRepo) *RBACService {
	return &RBACService{rbac: rbac}
}

// ResolvePermissionCodes returns merged, deduplicated permission keys from direct
// roles and from roles attached to groups the user belongs to.
func (s *RBACService) ResolvePermissionCodes(userID uint) ([]string, error) {
	return s.rbac.EffectivePermissionCodes(userID)
}
