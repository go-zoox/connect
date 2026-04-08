package model

// UserRole links users to roles.
type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

func (UserRole) TableName() string {
	return "admin_user_roles"
}

// RolePermission links roles to permissions.
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

func (RolePermission) TableName() string {
	return "admin_role_permissions"
}

// GroupUser links groups to users.
type GroupUser struct {
	GroupID uint `gorm:"primaryKey"`
	UserID  uint `gorm:"primaryKey"`
}

func (GroupUser) TableName() string {
	return "admin_group_users"
}

// GroupRole links groups to roles.
type GroupRole struct {
	GroupID uint `gorm:"primaryKey"`
	RoleID  uint `gorm:"primaryKey"`
}

func (GroupRole) TableName() string {
	return "admin_group_roles"
}

// MigrateModels returns core entities and relation tables for AutoMigrate.
func MigrateModels() []any {
	return []any{
		&User{},
		&Role{},
		&Permission{},
		&Group{},
		&UserRole{},
		&RolePermission{},
		&GroupUser{},
		&GroupRole{},
	}
}
