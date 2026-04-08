# Connect Built-in Admin + Database Design

## Background and Goal

`connect` currently defaults to remote service mode for core entities (`app`, `user`, `users`, `menus`, `permissions`).  
This design adds an optional built-in admin system backed by database storage, so third-party apps can bootstrap quickly without third-party OAuth2.

Primary goals:

- Keep current behavior unchanged by default.
- Enable built-in admin + DB only when explicitly configured.
- Support account/password login with session/cookie.
- Manage `app`, `user`, `users`, `menus`, RBAC (`roles`, `permissions`, `groups`) in one self-contained mode.
- Support frontend dev/prod workflow: split in development, embedded static assets in build output.

## Scope

In scope:

- Config-gated built-in admin mode.
- GORMX-based database integration with AutoMigrate.
- Built-in auth (username/password) for admin and normal users.
- API compatibility for existing built-in endpoints and extensions for roles/groups.
- Admin frontend mount path configurable (default `/admin`) and embedded on build.

Out of scope (this phase):

- Third-party OAuth2 integration in built-in mode.
- Versioned SQL migration framework (AutoMigrate only).
- Multi-tenant isolation model.

## Mode and Compatibility

### Default mode (unchanged)

- `admin.enabled=false` (default): keep existing service behavior.
- Existing config and runtime paths remain backward compatible.

### Built-in mode

- `admin.enabled=true`: built-in DB-backed handlers become active for:
  - existing: `/app`, `/user`, `/users`, `/menus`, `/permissions`
  - new: `/roles`, `/groups`
- Third-party OAuth2 is not used in this mode.
- Login uses local username/password + current session/cookie infrastructure.

## Configuration Design

Add a new top-level section:

```yaml
admin:
  enabled: false
  entry: /admin
  auth:
    admin:
      username: admin
      password: change_me
  database:
    # managed by gormx config style
    # driver: sqlite|mysql|postgres
    # dsn: ...
```

Rules:

- `admin.enabled` default is `false`.
- `admin.entry` default is `/admin`, can be customized.
- If `admin.enabled=true`, `admin.auth.admin.username/password` are required.
- If required admin config is missing when enabled, server startup fails fast with clear error.
- Database settings follow `gormx` conventions (no custom DB abstraction invented in `connect`).

## Architecture

Introduce internal admin module boundaries:

- `app/admin/model`: entity models and relation models.
- `app/admin/repository`: DB CRUD and relation operations.
- `app/admin/service`: domain logic for auth, users, groups, roles, permissions.
- `app/admin/api`: HTTP handlers for built-in admin endpoints.
- `app/admin/bootstrap`: DB init, AutoMigrate, seed admin account.

Router integration:

- Keep current router entrypoints.
- At boot, choose implementation by mode:
  - default: existing service-backed handlers.
  - built-in: admin module handlers.

This keeps endpoint compatibility while isolating built-in complexity.

## Data Model (RBAC)

Core tables:

- `users`
  - id, username (unique), password_hash, nickname, email, avatar
  - is_admin, status, created_at, updated_at
- `roles`
  - id, name (unique), code (unique), description
- `permissions`
  - id, name, code (unique), description
- `groups`
  - id, name (unique), code (unique), description

Relation tables:

- `user_roles` (user_id, role_id)
- `group_roles` (group_id, role_id)
- `group_users` (group_id, user_id)
- `role_permissions` (role_id, permission_id)

Computed permission strategy:

- User effective permissions =
  direct user roles permissions + group roles permissions.
- Service layer computes and deduplicates permission codes.

## Authentication and Session

- Built-in login endpoint validates username/password from DB.
- Password storage uses hash (bcrypt recommended for first implementation).
- On login success, issue existing session/cookie mechanism (no JWT in this phase).
- Logout and current-user behaviors remain aligned with existing middleware pattern.

## API Surface

### Compatible endpoints

- `/app`: app metadata/config used by frontend.
- `/user`: current user profile + effective permissions snapshot.
- `/users`: user CRUD and membership/role association operations.
- `/menus`: menu CRUD and ordering.
- `/permissions`: permission CRUD.

### New endpoints

- `/roles`: role CRUD + role-permission binding.
- `/groups`: group CRUD + group-role binding + group-member binding.

Behavioral requirement:

- These endpoints are active in built-in mode and unchanged in default mode.

## Frontend Integration

Development:

- Frontend/backend separated, using the external frontend workflow as reference.

Build/Runtime:

- Frontend static assets are embedded into `connect` binary.
- When `admin.enabled=true`, mount static files under `admin.entry` (default `/admin`).
- When `admin.enabled=false`, do not expose built-in admin UI routes.

## Initialization and Seed Strategy

When built-in mode is enabled:

1. Initialize GORMX DB connection.
2. Run `AutoMigrate` for all built-in admin models.
3. Ensure initial admin account exists using configured admin credentials.
   - If absent, create it.
   - If present, keep idempotent behavior (no destructive reset by default).

## Error Handling and Security

- Startup validation errors:
  - missing admin credentials when enabled
  - invalid DB config when enabled
- Runtime authorization:
  - route-level guards for admin-only operations
  - non-admin users only access allowed resources
- Security baseline:
  - hash passwords only
  - do not log plain-text credentials
  - preserve current cookie security options

## Testing Strategy

Unit tests:

- auth service (password verify, session subject generation)
- RBAC resolver (user/group/role/permission composition)
- bootstrap idempotency for admin seed

Integration tests:

- mode switch behavior (`enabled=false` vs `enabled=true`)
- CRUD flows for users/roles/groups/permissions/menus
- login/logout/current-user session lifecycle

Regression tests:

- default mode behavior unchanged for existing deployments

## Delivery Plan (Phased)

Phase 1: Config + bootstrap + DB models

- Add config schema and validation.
- Integrate GORMX and AutoMigrate.
- Seed initial admin account from config.

Phase 2: Built-in auth + core entity APIs

- Implement login/session for built-in mode.
- Implement `/user`, `/users`, `/permissions`, `/roles`, `/groups`.

Phase 3: Compatibility and frontend embedding

- Wire `/app`, `/menus` compatibility.
- Add embedded static admin frontend and configurable mount entry.

Phase 4: Hardening

- complete tests, edge-case handling, docs and examples.

## Open Decisions Resolved

- DB layer: use `gormx`.
- Frontend: dev split, build embed.
- Built-in mode activation: config flag only, default disabled.
- Auth: local username/password + session/cookie.
- Migration: AutoMigrate.
- Initial admin: must be provided in config.
- API strategy: keep existing endpoint compatibility and extend roles/groups.
