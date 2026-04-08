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
