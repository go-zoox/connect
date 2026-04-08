package bootstrap

import (
	"testing"

	"github.com/go-zoox/connect/admin/model"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/gormx"
	"golang.org/x/crypto/bcrypt"
)

// Shared-cache SQLite URI keeps a single in-memory database for repeated LoadDB
// calls in this process (ordinary :memory: would allocate a new DB each open).
const testSQLiteMemoryDSN = "file:bootstrap_admin_shared?mode=memory&cache=shared"

func TestBootstrapAdminMode(t *testing.T) {
	var cfg config.Config
	cfg.ApplyDefault()
	cfg.Admin.Enabled = true
	cfg.Admin.Auth.Admin.Username = "root"
	cfg.Admin.Auth.Admin.Password = "secret-pass"
	cfg.Admin.Database.Driver = "sqlite"
	cfg.Admin.Database.DSN = testSQLiteMemoryDSN

	if err := cfg.ValidateAdmin(); err != nil {
		t.Fatalf("ValidateAdmin: %v", err)
	}

	if err := Bootstrap(&cfg); err != nil {
		t.Fatalf("first Bootstrap: %v", err)
	}

	db := gormx.GetDB()
	for _, m := range model.MigrateModels() {
		if !db.Migrator().HasTable(m) {
			t.Fatalf("after bootstrap, expected migrated table for %T", m)
		}
	}

	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		t.Fatalf("count users: %v", err)
	}
	if count != 1 {
		t.Fatalf("after first bootstrap, want 1 admin user, got %d", count)
	}

	if err := Bootstrap(&cfg); err != nil {
		t.Fatalf("second Bootstrap: %v", err)
	}
	db = gormx.GetDB()
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		t.Fatalf("count users (2): %v", err)
	}
	if count != 1 {
		t.Fatalf("after second bootstrap, want still 1 user (idempotent seed), got %d", count)
	}

	var u model.User
	if err := db.Where("username = ?", "root").First(&u).Error; err != nil {
		t.Fatalf("load root user: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("secret-pass")); err != nil {
		t.Fatalf("stored password is not valid bcrypt for configured password: %v", err)
	}
}
