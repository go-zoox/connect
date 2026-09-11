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
- After each patch: run `go test ./...`, `go test -tags=integration ./...`, and `go vet ./...` before merge.

## Test strategy

### Unit tests

- **Config**: `app/config/*_test.go` — validation, defaults, admin entry safety.
- **Admin domain**: `app/admin/bootstrap/*_test.go`, `app/admin/service/*_test.go` — DB bootstrap, auth, RBAC resolution in isolation.

### Integration tests

- **Router + HTTP**: `app/router/*_test.go` — `httptest` against `router.New` + `defaults.Application()`; covers route registration, auth middleware interaction, admin on/off.
- **Full app wiring** (build tag `integration`): `app/connect_integration_test.go` — `Connect.Setup` then `TestHTTPHandler()` (see `export_test.go`) to exercise the same path as production **except listening** (OAuth2 register, embed loading, bootstrap, router).

Run integration-only file with the Go **`-tags`** flag (plural), not `-tag`:

```bash
go test -tags=integration ./...
```

Default `go test ./...` still runs unit + router tests; tagged tests are skipped unless you pass `-tags=integration`.

### Naming

- Prefix full-stack checks with `TestIntegration_` in `connect_integration_test.go`.
- Keep regression tests for **default mode** (`admin.enabled` false) whenever touching router or admin.

## Experience Log

### 2026-04-09 - Connect.Setup + integration tests

- Added `Connect.Setup` so tests can run full wiring without `Run` (listen).
- `app/connect_integration_test.go` locks default-mode backward compatibility and admin-mode login + API after `Setup`.
- `TestSetup_RejectsAdminEnabledWithoutCredentials` uses `Setup` (fails before listen).

### 2026-04-09 - Integration build tag

- `app/connect_integration_test.go` is behind `//go:build integration`; run with `go test -tags=integration ./...`.

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

### 2026-09-11 - Permissions inside X-Connect-Token

#### What worked

- `user.User.Encode/Decode` now carries the `permissions` claim, so upstream services behind the gateway resolve authorization from the token instead of calling back into the permissions service.
- Permissions are filled by `fillPermissions` right before each `X-Connect-Token` is signed (routes with `backend.secret_key`, upstream mode, frontend + backend mode), reusing the cached `service.GetPermission`.
- Tests: `user/user_test.go` locks the claim round trip; `app/router/router_upstream_token_test.go` asserts the header received by the upstream really contains the permissions.

#### Pitfalls found and fixed

- A permissions lookup failure must not break the request: it is logged and leaves permissions empty, which upstream services treat as deny-by-default.
- Proxy handlers cannot be exercised with `httptest.NewRecorder()` — the proxy needs a `http.CloseNotifier` capable `ResponseWriter`, use `httptest.NewServer`.

### 2026-09-11 - Where permissions come from

- `user.User.Permissions` was already filled by `service.GetUser`: the user service response (`result`) is unmarshalled into the user struct, and doreamon's resource service answers `result.permissions` (its `/oauth/user`, reached as `api.zcorky.com/user`, returns `{...user, permissions}`).
- The permissions service (`services.permissions.service`, default `https://api.zcorky.com/permissions`) is a second source and is only asked when the user service answered no permission. In the doreamon gateway (`config/gateway.yml`) no `permissions` service is mounted at all, so that default URL 404s; the app scoped endpoint is `https://api.zcorky.com/oauth/app/permissions` (gateway `oauth` service -> resource `/oauth`).
- The doreamon permission values are application scoped codes built from `role.permissions` + `role.menuPermissions` (for example `global.system.permissions`), not literals like `ADMIN`.
- `@znode/connect` (the Node side decoding `x-connect-token`, used by the doreamon resource service) only reads `id/nickname/avatar/email/username`, the extra `permissions` claim is additive and does not break it.
