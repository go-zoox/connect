# 命令行

二进制入口为 `connect`（包路径：`github.com/go-zoox/connect/cmd/connect`）。以下为 README 与仓库 Dockerfile 中约定的子命令；**完整 Flags 以本地 `connect --help` / `connect <command> --help` 为准**。

## 全局选项

常见全局项包括：

| 选项 | 含义 |
|------|------|
| `-h`, `--help` | 帮助 |
| `-v`, `--version` | 版本号 |

## 子命令一览

| 命令 | 作用 |
|------|------|
| `serve` | 读取 YAML 配置启动，例如 `connect serve -c ./conf/config.yml` |
| `none` | 无认证模式启动（需 upstream 或 frontend+backend） |
| `doreamon` | 使用 Doreamon OAuth2 预设与环境变量 |
| `github` | 使用 GitHub OAuth 预设与环境变量 |
| `feishu` | 使用飞书 OAuth 预设与环境变量 |

## 环境变量（通用）

多数变量可被 YAML 中的同类字段覆盖；下列名称与 `app/config/defaults.go` 中读取逻辑对应。

### 服务进程

| 变量 | 说明 |
|------|------|
| `PORT` | 监听端口，默认 `8080` |
| `MODE` | `development` / `production` 等 |
| `SECRET_KEY` | 会话加密密钥；未设则随机生成（不推荐生产） |
| `SESSION_MAX_AGE` | 会话时长（秒） |
| `LOG_LEVEL` | 日志级别 |

### 上游与路由

| 变量 | 说明 |
|------|------|
| `UPSTREAM` | 单一上游完整 URL，如 `https://example.com` |
| `FRONTEND` | 前端根 URL（与 `BACKEND` 成对使用） |
| `BACKEND` | 后端 API 根 URL |
| `DISABLE_PREFIX_REWRITE` | 是否禁用后端前缀重写（非空则禁用） |

### 认证总开关

| 变量 | 说明 |
|------|------|
| `AUTH_MODE` | 如 `oauth2`、`password`、`none` |
| `AUTH_PROVIDER` | OAuth2 时的 provider，如 `doreamon`、`github`、`feishu` |
| `AUTH_IGNORE_PATHS` | 逗号分隔的路径前缀，可跳过鉴权 |
| `AUTH_IS_IGNORE_PATHS_DISABLED` | `true` 时不使用忽略列表 |
| `AUTH_IS_IGNORE_WHEN_HEADER_AUTHORIZATION_FOUND` | 请求头带 Authorization 时可跳过部分逻辑 |

GitHub / 飞书场景可在 YAML 中配置 **`auth.allow_usernames`**，仅允许列表中的用户登录同名账户（环境变量以当前版本 `connect` 的 `defaults` 为准，优先使用 `connect serve -c` + YAML）。

### OAuth Client（便捷模式）

便捷子命令通常读取：

| 变量 | 说明 |
|------|------|
| `CLIENT_ID` | OAuth2 Client ID |
| `CLIENT_SECRET` | OAuth2 Client Secret |
| `REDIRECT_URI` | 必须与 IdP 注册一致，且含正确回调路径 |

### 内置 API 路径前缀

可通过 `BUILT_IN_APIS_*` 系列变量改写，例如 `BUILT_IN_APIS_APP`、`BUILT_IN_APIS_USER` 等；默认值多为 `/app`、`/user`、`/menus` 等（参见 [配置](./config.md)）。

## 与 YAML 配置的分工

- **一键脚本 / 容器**：多用环境变量。
- **复杂路由、多 oauth2、自定义 services**：使用 YAML + `connect serve -c`。
