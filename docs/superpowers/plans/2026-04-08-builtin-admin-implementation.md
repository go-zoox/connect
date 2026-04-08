# Built-in Admin and Database Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an optional built-in admin mode in `connect` that uses `gormx` + local username/password auth to manage app/user/users/menus/rbac while keeping default behavior unchanged.

**Architecture:** Keep existing routes and runtime defaults, then add a top-level `admin/` module that is activated only when `admin.enabled=true`. In built-in mode, core handlers delegate to DB-backed services, while default mode keeps current remote-service behavior. Frontend remains split in development and embedded for production.

**Tech Stack:** Go, `github.com/go-zoox/gormx`, Zoox router/middleware, session+cookie auth, go test

---

## File Structure Map

- Create: `app/config/admin.go` (admin config structs + validation helpers)
- Modify: `app/config/config.go` (add `Admin` field)
- Modify: `app/config/defaults.go` (default values + env mapping)
- Create: `admin/bootstrap/bootstrap.go` (init DB, auto migrate, seed admin)
- Create: `admin/model/*.go` (users/roles/permissions/groups + relation models)
- Create: `admin/repository/*.go` (CRUD + relation persistence)
- Create: `admin/service/*.go` (auth, rbac resolver, entity services)
- Create: `admin/api/*.go` (HTTP handlers for users/roles/groups/permissions/menus/app/user/login)
- Modify: `app/router/router.go` (mode switch and route wiring)
- Modify: `app/middleware/auth.go` (allow built-in mode auth flow where needed)
- Create: `admin/static/embed.go` (embedded assets mounting helper)
- Create: `admin/*_test.go` and `app/router/*_test.go` (unit/integration tests)
- Modify: `go.mod`, `go.sum` (add gormx and any direct deps)
- Modify: `README.md`, `conf/config.full.example` (docs and sample config)

### Task 1: Add admin config schema and validation gate

**Files:**
- Create: `app/config/admin.go`
- Modify: `app/config/config.go`
- Modify: `app/config/defaults.go`
- Test: `app/config/admin_test.go`

- [ ] **Step 1: Write failing tests for admin config defaults and required fields**

```go
func TestAdminDefaultDisabled(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefault()
	if cfg.Admin.Enabled {
		t.Fatalf("expected admin disabled by default")
	}
	if cfg.Admin.Entry != "/admin" {
		t.Fatalf("expected default admin entry /admin, got %s", cfg.Admin.Entry)
	}
}

func TestAdminEnabledRequiresRootAccount(t *testing.T) {
	cfg := &Config{}
	cfg.Admin.Enabled = true
	cfg.Admin.Auth.Admin.Username = ""
	cfg.Admin.Auth.Admin.Password = ""
	if err := cfg.Admin.Validate(); err == nil {
		t.Fatalf("expected validation error when admin credentials are missing")
	}
}
```

- [ ] **Step 2: Run config tests to confirm failure**

Run: `go test ./app/config -run TestAdmin -v`  
Expected: FAIL with missing `Admin` struct/validate logic.

- [ ] **Step 3: Implement config structs + defaults + validation**

```go
type AdminConfig struct {
	Enabled  bool            `config:"enabled"`
	Entry    string          `config:"entry"`
	Auth     AdminAuth       `config:"auth"`
	Database map[string]any  `config:"database"`
}

func (a *AdminConfig) ApplyDefault() {
	if a.Entry == "" {
		a.Entry = "/admin"
	}
}

func (a *AdminConfig) Validate() error {
	if !a.Enabled {
		return nil
	}
	if a.Auth.Admin.Username == "" || a.Auth.Admin.Password == "" {
		return fmt.Errorf("admin.auth.admin.username/password is required when admin.enabled=true")
	}
	return nil
}
```

- [ ] **Step 4: Run tests to verify pass**

Run: `go test ./app/config -run TestAdmin -v`  
Expected: PASS.

- [ ] **Step 5: Commit**

Run:
```bash
git add app/config/admin.go app/config/config.go app/config/defaults.go app/config/admin_test.go
git commit -m "feat(config): add built-in admin config and validation gate"
```

### Task 2: Add gormx bootstrap, automigrate, and admin seeding

**Files:**
- Create: `admin/bootstrap/bootstrap.go`
- Create: `admin/model/user.go`
- Create: `admin/model/role.go`
- Create: `admin/model/permission.go`
- Create: `admin/model/group.go`
- Create: `admin/model/relation.go`
- Test: `admin/bootstrap/bootstrap_test.go`

- [ ] **Step 1: Write failing bootstrap test for enabled mode**

```go
func TestBootstrapAdminMode_AutoMigrateAndSeed(t *testing.T) {
	cfg := testAdminConfigWithSQLiteMemory()
	db, err := bootstrap.Init(cfg)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}
	var count int64
	if err := db.Model(&model.User{}).Where("username = ?", cfg.Admin.Auth.Admin.Username).Count(&count).Error; err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected seeded admin user, got %d", count)
	}
}
```

- [ ] **Step 2: Run bootstrap tests to confirm failure**

Run: `go test ./admin/bootstrap -run TestBootstrapAdminMode -v`  
Expected: FAIL because bootstrap/models do not exist.

- [ ] **Step 3: Implement models, relation tables, and bootstrap init**

```go
func Init(cfg *config.Config) (*gorm.DB, error) {
	db, err := gormx.New(&gormx.Config{
		Driver: cfg.Admin.Database.Driver,
		DSN:    cfg.Admin.Database.DSN,
	})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&model.User{}, &model.Role{}, &model.Permission{}, &model.Group{},
		&model.UserRole{}, &model.GroupRole{}, &model.GroupUser{}, &model.RolePermission{},
	); err != nil {
		return nil, err
	}
	if err := seedAdmin(db, cfg); err != nil {
		return nil, err
	}
	return db, nil
}
```

- [ ] **Step 4: Run tests to verify pass**

Run: `go test ./admin/bootstrap -run TestBootstrapAdminMode -v`  
Expected: PASS.

- [ ] **Step 5: Commit**

Run:
```bash
git add admin/bootstrap/bootstrap.go admin/model/*.go admin/bootstrap/bootstrap_test.go go.mod go.sum
git commit -m "feat(admin): add gormx bootstrap with automigrate and admin seed"
```

### Task 3: Implement built-in auth and RBAC resolver services

**Files:**
- Create: `admin/service/auth.go`
- Create: `admin/service/rbac.go`
- Create: `admin/repository/user_repo.go`
- Create: `admin/repository/rbac_repo.go`
- Test: `admin/service/auth_test.go`
- Test: `admin/service/rbac_test.go`

- [ ] **Step 1: Write failing tests for login and effective permissions**

```go
func TestAuthLogin_Success(t *testing.T) {
	svc := newAuthServiceWithSeededUser(t)
	u, err := svc.Login("admin", "pass123")
	if err != nil || u.Username != "admin" {
		t.Fatalf("expected login success, got err=%v", err)
	}
}

func TestRBACResolvePermissions_FromRolesAndGroups(t *testing.T) {
	svc := newRBACServiceWithFixtures(t)
	codes, err := svc.ResolveUserPermissionCodes("u1")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	assertContains(t, codes, "menu:read")
	assertContains(t, codes, "user:write")
}
```

- [ ] **Step 2: Run service tests to confirm failure**

Run: `go test ./admin/service -run 'TestAuthLogin|TestRBACResolvePermissions' -v`  
Expected: FAIL due to missing service code.

- [ ] **Step 3: Implement auth + rbac service minimal pass**

```go
func (s *AuthService) Login(username, password string) (*model.User, error) {
	u, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredential
	}
	return u, nil
}

func (s *RBACService) ResolveUserPermissionCodes(userID string) ([]string, error) {
	codes, err := s.repo.ListPermissionCodesByUserAndGroups(userID)
	if err != nil {
		return nil, err
	}
	return dedupe(codes), nil
}
```

- [ ] **Step 4: Run service tests to verify pass**

Run: `go test ./admin/service -run 'TestAuthLogin|TestRBACResolvePermissions' -v`  
Expected: PASS.

- [ ] **Step 5: Commit**

Run:
```bash
git add admin/service/auth.go admin/service/rbac.go admin/repository/*.go admin/service/*_test.go
git commit -m "feat(admin): add built-in auth and rbac resolver services"
```

### Task 4: Wire built-in API handlers and route switch

**Files:**
- Create: `admin/api/login.go`
- Create: `admin/api/users.go`
- Create: `admin/api/roles.go`
- Create: `admin/api/groups.go`
- Create: `admin/api/permissions.go`
- Create: `admin/api/menus.go`
- Create: `admin/api/app.go`
- Create: `admin/api/user.go`
- Modify: `app/router/router.go`
- Test: `app/router/router_admin_mode_test.go`

- [ ] **Step 1: Write failing router integration test for mode switch**

```go
func TestRouter_AdminEnabled_RegistersBuiltinAdminEndpoints(t *testing.T) {
	app, cfg := newTestAppAndConfig(t)
	cfg.Admin.Enabled = true
	router.New(app, cfg)
	resp := performRequest(app, "GET", "/api"+cfg.BuiltInAPIs.Roles, "")
	if resp.Code == http.StatusNotFound {
		t.Fatalf("expected roles endpoint registered when admin enabled")
	}
}
```

- [ ] **Step 2: Run router tests to confirm failure**

Run: `go test ./app/router -run TestRouter_AdminEnabled -v`  
Expected: FAIL due to missing roles route/admin mode switch.

- [ ] **Step 3: Implement handlers and route selection**

```go
if cfg.Admin.Enabled {
	group.Get(cfg.BuiltInAPIs.App, adminAPIApp.New(cfg, adminContainer))
	group.Get(cfg.BuiltInAPIs.User, adminAPIUser.New(cfg, adminContainer))
	group.Get(cfg.BuiltInAPIs.Users, adminAPIUsers.List(cfg, adminContainer))
	group.Get(cfg.BuiltInAPIs.Permissions, adminAPIPermissions.List(cfg, adminContainer))
	group.Get(cfg.BuiltInAPIs.Roles, adminAPIRoles.List(cfg, adminContainer))
	group.Get(cfg.BuiltInAPIs.Groups, adminAPIGroups.List(cfg, adminContainer))
	group.Post(cfg.BuiltInAPIs.Login, adminAPILogin.New(cfg, adminContainer))
} else {
	group.Get(cfg.BuiltInAPIs.App, apiApp.New(cfg))
	group.Get(cfg.BuiltInAPIs.User, apiUser.New(cfg))
	group.Get(cfg.BuiltInAPIs.Menus, apiMenus.New(cfg))
	group.Get(cfg.BuiltInAPIs.Permissions, apiPermissions.New(cfg))
	group.Get(cfg.BuiltInAPIs.Users, apiUser.GetUsers(cfg))
	group.Post(cfg.BuiltInAPIs.Login, apiUser.Login(cfg))
}
```

- [ ] **Step 4: Run router tests to verify pass**

Run: `go test ./app/router -run TestRouter_AdminEnabled -v`  
Expected: PASS.

- [ ] **Step 5: Commit**

Run:
```bash
git add admin/api/*.go app/router/router.go app/router/router_admin_mode_test.go
git commit -m "feat(router): switch to built-in admin handlers when enabled"
```

### Task 5: Embed admin frontend and configurable entry mount

**Files:**
- Create: `admin/static/embed.go`
- Modify: `app/router/router.go`
- Modify: `conf/config.full.example`
- Test: `app/router/router_admin_static_test.go`

- [ ] **Step 1: Write failing test for admin entry static mount**

```go
func TestRouter_AdminEnabled_MountsStaticAtConfiguredEntry(t *testing.T) {
	app, cfg := newTestAppAndConfig(t)
	cfg.Admin.Enabled = true
	cfg.Admin.Entry = "/console"
	router.New(app, cfg)
	resp := performRequest(app, "GET", "/console", "")
	if resp.Code == http.StatusNotFound {
		t.Fatalf("expected admin entry static route")
	}
}
```

- [ ] **Step 2: Run static route test to confirm failure**

Run: `go test ./app/router -run TestRouter_AdminEnabled_MountsStaticAtConfiguredEntry -v`  
Expected: FAIL because static embedding/mount not wired.

- [ ] **Step 3: Add embed FS and conditional mount**

```go
//go:embed dist/*
var AdminFS embed.FS

func Mount(app *zoox.Application, entry string) {
	app.StaticFS(entry+"/assets", http.FS(distFS))
	app.Get(entry, indexHandler)
	app.Get(entry+"/*", spaFallbackHandler)
}
```

- [ ] **Step 4: Run static route test to verify pass**

Run: `go test ./app/router -run TestRouter_AdminEnabled_MountsStaticAtConfiguredEntry -v`  
Expected: PASS.

- [ ] **Step 5: Commit**

Run:
```bash
git add admin/static/embed.go app/router/router.go conf/config.full.example app/router/router_admin_static_test.go
git commit -m "feat(admin-ui): embed admin frontend and mount at configurable entry"
```

### Task 6: Backward compatibility regression tests and docs

**Files:**
- Create: `app/router/router_default_mode_regression_test.go`
- Modify: `README.md`
- Modify: `conf/config.local.yml.example`
- Modify: `conf/config.full.example`

- [ ] **Step 1: Write failing regression test for default mode**

```go
func TestRouter_DefaultMode_KeepsExistingBehavior(t *testing.T) {
	app, cfg := newTestAppAndConfig(t)
	cfg.Admin.Enabled = false
	router.New(app, cfg)
	resp := performRequest(app, "GET", "/api"+cfg.BuiltInAPIs.User, "")
	if resp.Code == http.StatusNotFound {
		t.Fatalf("expected existing built-in user endpoint in default mode")
	}
}
```

- [ ] **Step 2: Run regression test to confirm failure (if any)**

Run: `go test ./app/router -run TestRouter_DefaultMode_KeepsExistingBehavior -v`  
Expected: FAIL if route switch broke defaults; otherwise PASS and keep as guard.

- [ ] **Step 3: Update docs with complete config and mode explanation**

```yaml
admin:
  enabled: false
  entry: /admin
  auth:
    admin:
      username: "admin"
      password: "strong-password"
  database:
    driver: sqlite
    dsn: "file:./data/connect-admin.db?cache=shared"
```

- [ ] **Step 4: Run full targeted verification**

Run: `go test ./app/config ./admin/... ./app/router -v`  
Expected: PASS.

- [ ] **Step 5: Commit**

Run:
```bash
git add app/router/router_default_mode_regression_test.go README.md conf/config.local.yml.example conf/config.full.example
git commit -m "test/docs: add default-mode regression coverage and admin config docs"
```

## Final Verification Checklist

- [ ] Run: `go test ./...`
- [ ] Run: `go vet ./...`
- [ ] Start default mode manually and verify existing flow unchanged.
- [ ] Start with `admin.enabled=true` and verify login, users, roles, groups, permissions, menus, and `/app` APIs.
- [ ] Verify admin UI is only reachable when enabled and mounted at configured `admin.entry`.
