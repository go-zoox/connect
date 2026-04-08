# AGENTS

## Workflow Rules

1. Before starting any new requirement, read this file first.
2. After finishing a requirement, append/update the "Experience Log" section.
3. Keep default behavior backward-compatible unless the requirement explicitly says otherwise.
4. For security-sensitive routing/auth changes, add tests before claiming completion.

## Experience Log

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
