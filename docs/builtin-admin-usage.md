# Built-in Admin Usage Guide

## Overview

`connect` now supports an optional built-in admin mode with local database storage.

- Default behavior stays unchanged (`admin.enabled: false`)
- Built-in mode is activated only when `admin.enabled: true`
- Built-in mode uses local account/password auth + session/cookie
- Built-in mode uses `gormx` for DB initialization and `AutoMigrate`

## Quick Start (SQLite)

Use this minimal config:

```yaml
admin:
  enabled: true
  entry: /admin
  auth:
    admin:
      username: admin
      password: "change-me-now"
  database:
    driver: sqlite
    dsn: "file:./data/connect-admin.db?cache=shared"
```

Run:

```bash
connect serve -c ./conf/config.local.yml
```

Open admin UI:

- `http://127.0.0.1:8080/admin`

## Configuration

### `admin.enabled`

- `false` (default): keep existing behavior
- `true`: enable built-in admin mode

### `admin.entry`

- Default: `/admin`
- Must not overlap API/auth routes
- Invalid examples:
  - `/`
  - `/api`
  - `/api/app`
  - `/login`
  - value equal to `backend.prefix`

### `admin.auth.admin`

Required when enabled:

- `username`
- `password`

If missing, startup fails fast.

### `admin.database`

Required when enabled:

- `driver`
- `dsn`

Configured through `gormx`.

## Built-in Mode API Behavior

When `admin.enabled: true`:

- Built-in handlers are mounted under `/api`
- Includes:
  - existing: `app`, `user`, `users`, `menus`, `permissions`, `login`
  - new: `roles`, `groups`
- Public built-in prefix (`/api/_/*`) does not expose admin data endpoints

When `admin.enabled: false`:

- Existing route behavior remains unchanged

## Data Initialization

On startup (built-in mode only):

1. Validate admin config
2. Initialize database with `gormx`
3. Run `AutoMigrate`
4. Seed root admin user if missing

Password is stored as bcrypt hash.

## Development and Build

- Development: frontend/backend separated
- Build/runtime: admin static assets embedded in `connect` binary
- Admin UI is mounted only in built-in mode

## Verification Checklist

- Automated: `go test ./...` (unit + router), then `go test -tags=integration ./...` (includes full-app `Connect.Setup` tests; Go uses **`-tags`**, not `-tag`).
- Start with `admin.enabled: false` and verify existing flow is unchanged
- Start with `admin.enabled: true` and verify:
  - admin login works
  - `/api/roles` and `/api/groups` are reachable with auth
  - `/api/_/roles` is not exposing admin data
  - admin UI available at `admin.entry`
