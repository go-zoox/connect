package bootstrap

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-zoox/connect/admin/model"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/gormx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Bootstrap initializes the admin database via gormx when admin is enabled:
// connects using admin.database driver/dsn, auto-migrates models, and seeds
// the root admin user from config if no user with that username exists.
func Bootstrap(cfg *config.Config) error {
	if cfg == nil || !cfg.Admin.Enabled {
		return nil
	}

	if err := cfg.ValidateAdmin(); err != nil {
		return err
	}

	engine, err := adminEngine(cfg.Admin.Database.Driver)
	if err != nil {
		return err
	}

	dsn := strings.TrimSpace(cfg.Admin.Database.DSN)

	if err := gormx.LoadDB(engine, dsn); err != nil {
		return err
	}

	db := gormx.GetDB()
	if err := db.AutoMigrate(model.MigrateModels()...); err != nil {
		return fmt.Errorf("admin: automigrate: %w", err)
	}

	return seedRootAdmin(db, cfg)
}

func adminEngine(driver config.AdminDatabaseDriver) (string, error) {
	d := strings.TrimSpace(strings.ToLower(string(driver)))
	switch d {
	case "mysql":
		return "mysql", nil
	case "postgres", "postgresql":
		return "postgres", nil
	case "sqlite", "sqlite3":
		return "sqlite", nil
	default:
		return "", fmt.Errorf("admin: unsupported database.driver %q (use mysql, postgres, or sqlite)", driver)
	}
}

func seedRootAdmin(db *gorm.DB, cfg *config.Config) error {
	username := strings.TrimSpace(cfg.Admin.Auth.Admin.Username)
	password := cfg.Admin.Auth.Admin.Password

	var existing model.User
	err := db.Where("username = ?", username).First(&existing).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("admin: check root user: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("admin: hash password: %w", err)
	}

	u := &model.User{
		Username:     username,
		PasswordHash: string(hash),
	}
	if err := db.Create(u).Error; err != nil {
		return fmt.Errorf("admin: seed root user: %w", err)
	}

	return nil
}
