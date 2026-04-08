package repository

import (
	"sort"

	"github.com/go-zoox/connect/app/admin/model"
	"gorm.io/gorm"
)

// RBACRepo resolves permission keys from admin RBAC tables.
type RBACRepo struct {
	db *gorm.DB
}

// NewRBACRepo returns an RBACRepo backed by db.
func NewRBACRepo(db *gorm.DB) *RBACRepo {
	return &RBACRepo{db: db}
}

// EffectivePermissionCodes returns distinct permission keys for userID from:
// direct user roles and roles inherited via group membership.
func (r *RBACRepo) EffectivePermissionCodes(userID uint) ([]string, error) {
	var directRoleIDs []uint
	if err := r.db.Model(&model.UserRole{}).Where("user_id = ?", userID).Pluck("role_id", &directRoleIDs).Error; err != nil {
		return nil, err
	}

	var groupIDs []uint
	if err := r.db.Model(&model.GroupUser{}).Where("user_id = ?", userID).Pluck("group_id", &groupIDs).Error; err != nil {
		return nil, err
	}

	var groupRoleIDs []uint
	if len(groupIDs) > 0 {
		if err := r.db.Model(&model.GroupRole{}).Where("group_id IN ?", groupIDs).Pluck("role_id", &groupRoleIDs).Error; err != nil {
			return nil, err
		}
	}

	roleSet := make(map[uint]struct{})
	for _, id := range directRoleIDs {
		roleSet[id] = struct{}{}
	}
	for _, id := range groupRoleIDs {
		roleSet[id] = struct{}{}
	}
	if len(roleSet) == 0 {
		return nil, nil
	}

	roleIDs := make([]uint, 0, len(roleSet))
	for id := range roleSet {
		roleIDs = append(roleIDs, id)
	}

	var permissionIDs []uint
	if err := r.db.Model(&model.RolePermission{}).Where("role_id IN ?", roleIDs).Pluck("permission_id", &permissionIDs).Error; err != nil {
		return nil, err
	}
	if len(permissionIDs) == 0 {
		return nil, nil
	}

	permIDSet := make(map[uint]struct{})
	for _, id := range permissionIDs {
		permIDSet[id] = struct{}{}
	}
	ids := make([]uint, 0, len(permIDSet))
	for id := range permIDSet {
		ids = append(ids, id)
	}

	var keys []string
	if err := r.db.Model(&model.Permission{}).Where("id IN ?", ids).Pluck("key", &keys).Error; err != nil {
		return nil, err
	}

	keySet := make(map[string]struct{})
	for _, k := range keys {
		keySet[k] = struct{}{}
	}
	out := make([]string, 0, len(keySet))
	for k := range keySet {
		out = append(out, k)
	}
	sort.Strings(out)
	return out, nil
}
