package service

import (
	"reflect"
	"testing"

	"github.com/go-zoox/connect/app/admin/model"
	"github.com/go-zoox/connect/app/admin/repository"
)

func TestRBACResolvePermissions(t *testing.T) {
	db := openAdminTestDB(t)

	pRead := &model.Permission{Key: "users.read"}
	pWrite := &model.Permission{Key: "users.write"}
	for _, p := range []*model.Permission{pRead, pWrite} {
		if err := db.Create(p).Error; err != nil {
			t.Fatalf("create permission: %v", err)
		}
	}

	roleDirect := &model.Role{Name: "direct-role"}
	roleViaGroup := &model.Role{Name: "group-role"}
	for _, r := range []*model.Role{roleDirect, roleViaGroup} {
		if err := db.Create(r).Error; err != nil {
			t.Fatalf("create role: %v", err)
		}
	}

	// Direct role: only users.read
	if err := db.Create(&model.RolePermission{RoleID: roleDirect.ID, PermissionID: pRead.ID}).Error; err != nil {
		t.Fatalf("role perm direct: %v", err)
	}
	// Group role: users.read (duplicate path) + users.write
	for _, pid := range []uint{pRead.ID, pWrite.ID} {
		if err := db.Create(&model.RolePermission{RoleID: roleViaGroup.ID, PermissionID: pid}).Error; err != nil {
			t.Fatalf("role perm group: %v", err)
		}
	}

	u := &model.User{Username: "rbac-user", PasswordHash: "x"}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("user: %v", err)
	}

	if err := db.Create(&model.UserRole{UserID: u.ID, RoleID: roleDirect.ID}).Error; err != nil {
		t.Fatalf("user role: %v", err)
	}

	g := &model.Group{Name: "g1"}
	if err := db.Create(g).Error; err != nil {
		t.Fatalf("group: %v", err)
	}
	if err := db.Create(&model.GroupUser{GroupID: g.ID, UserID: u.ID}).Error; err != nil {
		t.Fatalf("group user: %v", err)
	}
	if err := db.Create(&model.GroupRole{GroupID: g.ID, RoleID: roleViaGroup.ID}).Error; err != nil {
		t.Fatalf("group role: %v", err)
	}

	svc := NewRBACService(repository.NewRBACRepo(db))
	codes, err := svc.ResolvePermissionCodes(u.ID)
	if err != nil {
		t.Fatalf("ResolvePermissionCodes: %v", err)
	}

	want := []string{"users.read", "users.write"}
	if !reflect.DeepEqual(codes, want) {
		t.Fatalf("codes = %v, want %v", codes, want)
	}
}

func TestRBACExcludesSoftDeletedDirectRole(t *testing.T) {
	db := openAdminTestDB(t)

	p := &model.Permission{Key: "perm.direct"}
	if err := db.Create(p).Error; err != nil {
		t.Fatalf("permission: %v", err)
	}
	role := &model.Role{Name: "r-del"}
	if err := db.Create(role).Error; err != nil {
		t.Fatalf("role: %v", err)
	}
	if err := db.Create(&model.RolePermission{RoleID: role.ID, PermissionID: p.ID}).Error; err != nil {
		t.Fatalf("role perm: %v", err)
	}
	u := &model.User{Username: "u-soft-role", PasswordHash: "x"}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("user: %v", err)
	}
	if err := db.Create(&model.UserRole{UserID: u.ID, RoleID: role.ID}).Error; err != nil {
		t.Fatalf("user role: %v", err)
	}
	if err := db.Delete(&model.Role{}, role.ID).Error; err != nil {
		t.Fatalf("soft delete role: %v", err)
	}

	svc := NewRBACService(repository.NewRBACRepo(db))
	codes, err := svc.ResolvePermissionCodes(u.ID)
	if err != nil {
		t.Fatalf("ResolvePermissionCodes: %v", err)
	}
	if len(codes) != 0 {
		t.Fatalf("expected no permissions from soft-deleted role, got %v", codes)
	}
}

func TestRBACExcludesSoftDeletedGroup(t *testing.T) {
	db := openAdminTestDB(t)

	pDirect := &model.Permission{Key: "perm.keep"}
	pGroup := &model.Permission{Key: "perm.drop"}
	for _, p := range []*model.Permission{pDirect, pGroup} {
		if err := db.Create(p).Error; err != nil {
			t.Fatalf("permission: %v", err)
		}
	}
	roleDirect := &model.Role{Name: "r-keep"}
	roleGroup := &model.Role{Name: "r-group"}
	for _, r := range []*model.Role{roleDirect, roleGroup} {
		if err := db.Create(r).Error; err != nil {
			t.Fatalf("role: %v", err)
		}
	}
	if err := db.Create(&model.RolePermission{RoleID: roleDirect.ID, PermissionID: pDirect.ID}).Error; err != nil {
		t.Fatalf("role perm direct: %v", err)
	}
	if err := db.Create(&model.RolePermission{RoleID: roleGroup.ID, PermissionID: pGroup.ID}).Error; err != nil {
		t.Fatalf("role perm group: %v", err)
	}
	u := &model.User{Username: "u-soft-group", PasswordHash: "x"}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("user: %v", err)
	}
	if err := db.Create(&model.UserRole{UserID: u.ID, RoleID: roleDirect.ID}).Error; err != nil {
		t.Fatalf("user role: %v", err)
	}
	g := &model.Group{Name: "g-del"}
	if err := db.Create(g).Error; err != nil {
		t.Fatalf("group: %v", err)
	}
	if err := db.Create(&model.GroupUser{GroupID: g.ID, UserID: u.ID}).Error; err != nil {
		t.Fatalf("group user: %v", err)
	}
	if err := db.Create(&model.GroupRole{GroupID: g.ID, RoleID: roleGroup.ID}).Error; err != nil {
		t.Fatalf("group role: %v", err)
	}
	if err := db.Delete(&model.Group{}, g.ID).Error; err != nil {
		t.Fatalf("soft delete group: %v", err)
	}

	svc := NewRBACService(repository.NewRBACRepo(db))
	codes, err := svc.ResolvePermissionCodes(u.ID)
	if err != nil {
		t.Fatalf("ResolvePermissionCodes: %v", err)
	}
	want := []string{"perm.keep"}
	if !reflect.DeepEqual(codes, want) {
		t.Fatalf("codes = %v, want %v", codes, want)
	}
}

func TestRBACExcludesSoftDeletedGroupRole(t *testing.T) {
	db := openAdminTestDB(t)

	p := &model.Permission{Key: "perm.via-role"}
	if err := db.Create(p).Error; err != nil {
		t.Fatalf("permission: %v", err)
	}
	role := &model.Role{Name: "r-gr-del"}
	if err := db.Create(role).Error; err != nil {
		t.Fatalf("role: %v", err)
	}
	if err := db.Create(&model.RolePermission{RoleID: role.ID, PermissionID: p.ID}).Error; err != nil {
		t.Fatalf("role perm: %v", err)
	}
	u := &model.User{Username: "u-soft-gr", PasswordHash: "x"}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("user: %v", err)
	}
	g := &model.Group{Name: "g1"}
	if err := db.Create(g).Error; err != nil {
		t.Fatalf("group: %v", err)
	}
	if err := db.Create(&model.GroupUser{GroupID: g.ID, UserID: u.ID}).Error; err != nil {
		t.Fatalf("group user: %v", err)
	}
	if err := db.Create(&model.GroupRole{GroupID: g.ID, RoleID: role.ID}).Error; err != nil {
		t.Fatalf("group role: %v", err)
	}
	if err := db.Delete(&model.Role{}, role.ID).Error; err != nil {
		t.Fatalf("soft delete role: %v", err)
	}

	svc := NewRBACService(repository.NewRBACRepo(db))
	codes, err := svc.ResolvePermissionCodes(u.ID)
	if err != nil {
		t.Fatalf("ResolvePermissionCodes: %v", err)
	}
	if len(codes) != 0 {
		t.Fatalf("expected no permissions from soft-deleted group-linked role, got %v", codes)
	}
}
