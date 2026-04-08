# AGENTS

## Workflow Rules

1. Before starting any new requirement, read this file first.
2. After finishing a requirement, append/update the "Experience Log" section.
3. Keep default behavior backward-compatible unless the requirement explicitly says otherwise.
4. For security-sensitive routing/auth changes, add tests before claiming completion.

## Iterative upgrades (patch-style)

- Treat new behavior as **additive patches**: small, reviewable diffs, one concern per change when possible.
- **Compatibility first**: existing configs and default-off code paths must keep working; gate new behavior behind config flags or explicit opt-in when breaking risk exists.
- Prefer extending existing modules (`app/config`, `app/router`, `app/admin/...`) over parallel implementations; avoid copy-paste handlers.
- After each patch: run `go test ./...` and `go vet ./...` before merge.

## Test strategy

### Unit tests

- **Config**: `app/config/*_test.go` — validation, defaults, admin entry safety.
- **Admin domain**: `app/admin/bootstrap/*_test.go`, `app/admin/service/*_test.go` — DB bootstrap, auth, RBAC resolution in isolation.

### Integration tests

- **Router + HTTP**: `app/router/*_test.go` — `httptest` against `router.New` + `defaults.Application()`; covers route registration, auth middleware interaction, admin on/off.
- **Full app wiring**: `app/connect_integration_test.go` — `Connect.Setup` then `TestHTTPHandler()` (see `export_test.go`) to exercise the same path as production **except listening** (OAuth2 register, embed loading, bootstrap, router).

### Naming

- Prefix full-stack checks with `TestIntegration_` in `connect_integration_test.go`.
- Keep regression tests for **default mode** (`admin.enabled` false) whenever touching router or admin.

## Experience Log

### 2026-04-09 - Connect.Setup + integration tests

- Added `Connect.Setup` so tests can run full wiring without `Run` (listen).
- `app/connect_integration_test.go` locks default-mode backward compatibility and admin-mode login + API after `Setup`.
- `TestSetup_RejectsAdminEnabledWithoutCredentials` uses `Setup` (fails before listen).

### 2026-04-09 - Admin package layout

- Built-in admin code lives under `app/admin/`. Import path: `github.com/go-zoox/connect/app/admin/...`.
- `app/router` and `app/app` import `app/admin/bootstrap`, `app/admin/api`, `app/admin/static`.

### 2026-04-08 - Built-in Admin + DB rollout

#### What worked

- Keep built-in admin strictly config-gated (`admin.enabled`) and default-off.
- Fail fast at startup for invalid admin config (root credentials, DB driver/DSN, unsafe admin entry).
- Use `gormx` + `AutoMigrate` + idempotent root-user seed for fast bootstrap.
- Add route-switch tests for `admin.enabled=true/false` to prevent regressions.

#### Pitfalls found and fixed

- Do not expose admin data APIs on public built-in prefix (`/api/_/*`).
- Do not allow `admin.entry` to overlap `/api` tree, auth paths, `/`, or backend prefix; otherwise auth ignore rules can bypass protection.
- RBAC resolution must ignore soft-deleted roles/groups.
- Admin login cannot be a stub success path; it must validate credentials.

#### Implementation guardrails for future requirements

- When introducing a new path that may be ignored by auth middleware, validate path safety in config.
- Any new admin endpoint must be tested in both modes:
  - `admin.enabled=false` (legacy behavior unchanged)
  - `admin.enabled=true` (admin behavior active)
- Public-prefix APIs should only expose explicitly safe endpoints.
- Document user-facing changes in `README.md` and add focused docs under `docs/`.
