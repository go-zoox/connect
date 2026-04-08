# Built-in Admin Implementation Report

## Goal

Deliver an optional built-in admin + database mode in `connect`, default-off, for fast third-party bootstrap without third-party OAuth2.

## Delivered Scope

### 1. Config and startup guard

- Added `admin` config block:
  - `enabled`
  - `entry`
  - `auth.admin.username/password`
  - `database.driver/dsn`
- Added validation and fail-fast startup checks when enabled

### 2. Database bootstrap (`gormx`)

- Added admin bootstrap flow:
  - DB init from admin config
  - `AutoMigrate`
  - root admin seeding
- Added idempotent seed behavior
- Password hashing uses bcrypt

### 3. Built-in auth and RBAC service layer

- Implemented username/password verification service
- Implemented RBAC effective permission resolver:
  - direct user roles permissions
  - group roles permissions
  - deduplicated permission codes
- Added soft-delete aware RBAC behavior for role/group resolution

### 4. Router mode switching

- Added `admin.enabled` based route switching
- In built-in mode, `/api` switches to built-in handlers
- Added new built-in endpoints: `/roles`, `/groups`
- In default mode, legacy behavior is preserved

### 5. Admin UI embedding and mount

- Embedded admin static assets in binary
- Mounted under configurable `admin.entry` (default `/admin`)
- Mounted only when built-in mode is enabled

### 6. Security hardening done in this delivery

- Rejected unsafe `admin.entry` values that can overlap protected API/auth paths
- Removed admin data endpoint exposure under public built-in prefix (`/api/_/*`) in built-in mode

### 7. Docs and regression tests

- Updated README and config examples
- Added default-mode regression tests and admin-mode route tests

## Validation Evidence

Executed in feature worktree:

- `go test ./...` passed
- `go vet ./...` passed

## Commits (feature branch)

- `43c8691` fix(router): block admin data endpoints on public built-in prefix
- `cd5bb2f` test/docs: add default-mode regression coverage and admin config docs
- `82ee44d` fix(admin-ui): disallow admin entry under api prefix
- `a7c19dc` fix(admin-ui): validate admin entry to avoid auth bypass overlaps
- `2d4e6ba` feat(admin-ui): embed admin frontend and mount at configurable entry
- `68e97f1` fix(admin): use real credential validation for builtin login
- `29310dd` feat(router): switch to built-in admin handlers when enabled
- `e9ee432` fix(admin): respect soft deletes in rbac permission resolution
- `5539e58` feat(admin): add built-in auth and rbac resolver services
- `549fd67` fix(admin): run bootstrap during startup with db config validation
- `1d76c11` feat(admin): add gormx bootstrap with automigrate and admin seed
- `cccde4a` fix(config): enforce admin credential validation during startup
- `7909260` feat(config): add built-in admin config and validation gate

## Known Follow-ups (non-blocking)

- Add rate limiting / lock policy for admin password login
- Clarify production deployment guidance when `services.user` is not local
- Consider reducing global DB handle coupling from `gormx` singleton usage
