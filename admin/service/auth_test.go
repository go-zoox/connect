package service

import (
	"strings"
	"testing"

	"github.com/go-zoox/connect/admin/model"
	"github.com/go-zoox/connect/admin/repository"
	"github.com/go-zoox/gormx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func openAdminTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	name := strings.ReplaceAll(t.Name(), "/", "_")
	dsn := "file:" + name + "?mode=memory&cache=shared"
	if err := gormx.LoadDB("sqlite", dsn); err != nil {
		t.Fatalf("LoadDB: %v", err)
	}
	db := gormx.GetDB()
	if err := db.AutoMigrate(model.MigrateModels()...); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}
	return db
}

func TestAuthLogin(t *testing.T) {
	db := openAdminTestDB(t)
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-horse"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if err := db.Create(&model.User{Username: "alice", PasswordHash: string(hash)}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	auth := NewAuthService(repository.NewUserRepo(db))

	t.Run("success", func(t *testing.T) {
		res, err := auth.Login("alice", "correct-horse")
		if err != nil {
			t.Fatalf("Login: %v", err)
		}
		if res.Status != AuthLoginOK {
			t.Fatalf("status: got %q want %q", res.Status, AuthLoginOK)
		}
		if res.User == nil || res.User.Username != "alice" {
			t.Fatalf("user: %+v", res.User)
		}
		if res.User.PasswordHash != "" {
			t.Fatal("expected password hash cleared on success")
		}
	})

	t.Run("invalid credential", func(t *testing.T) {
		res, err := auth.Login("alice", "wrong-password")
		if err != nil {
			t.Fatalf("Login: %v", err)
		}
		if res.Status != AuthLoginInvalidCredentials {
			t.Fatalf("wrong password: status %q", res.Status)
		}
		if res.User != nil {
			t.Fatalf("expected no user, got %+v", res.User)
		}

		res2, err := auth.Login("nobody", "correct-horse")
		if err != nil {
			t.Fatalf("Login: %v", err)
		}
		if res2.Status != AuthLoginInvalidCredentials {
			t.Fatalf("unknown user: status %q", res2.Status)
		}
	})
}
