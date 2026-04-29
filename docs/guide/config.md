# 配置

Connect 使用 **YAML**，通过 [`github.com/go-zoox/config`](https://github.com/go-zoox/config) 映射到 `app/config.Config`。命令行入口一般为：

```bash
connect serve -c /absolute/or/relative/path/config.yml
```

加载完成后仍会执行 **`Config.ApplyDefault()`**：缺失字段会被补上默认值；与环境变量合并时，常见做法是 **先有 YAML，再用 env 覆盖**（详见文末「环境变量」）。

---

## 仓库里的示例文件

| 文件 | 用途 |
|------|------|
| `conf/config.full.example` | OAuth2 + 前后端 + password + services + routes |
| `conf/config.oauth.yml.example` | OAuth2 + 前后端 |
| `conf/config.upstream.yml.example` | 单一 upstream + OAuth2 |
| `conf/config.local.yml.example` | 本地密码登录 + `services` 全 local |

下列字段名均来自源码中的 `config:"..."` 标签；YAML 中请使用 **snake_case**（与 Go 结构体标签一致）。

---

## 顶层字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `port` | int | HTTP 监听端口；默认 `8080`，可被环境变量 `PORT` 覆盖 |
| `mode` | string | `development` / `production`；影响静态资源、Cookie 策略等 |
| `secret_key` | string | Session 加密密钥；生产环境必填稳定值；可被 `SECRET_KEY` 覆盖 |
| `session_max_age` | int | 会话有效期，**单位为秒**（例如一天为 `86400`）；默认约一天量级 |
| `log_level` | string | 日志级别；可被 `LOG_LEVEL` 覆盖 |
| `loading_html` | string | 加载页 HTML（可选） |
| `index_html` | string | 首页 HTML（可选） |

::: tip `session_max_age` 常见误区
若误填成毫秒级大数（如 `86400000` 当作「一天」），会话会过长。请确认单位为 **秒**。
:::

---

## Upstream（单上游）

与 **frontend/backend 互斥使用**：要么只配 `upstream`，要么配 `frontend` + `backend`。

| 字段 | 说明 |
|------|------|
| `protocol` | `http` 或 `https`；省略时后续逻辑会补默认 |
| `host` | 主机名或 IP |
| `port` | 端口；`http` 默认 80、`https` 默认 443（见 `UpstreamService.String()`） |
| `change_origin` | 是否改写 Origin（按部署需求） |

等价环境变量： **`UPSTREAM`**（完整 URL，如 `https://httpbin.zcorky.com`）。

YAML 键名必须是 **`protocol`**，不要使用 `scheme`。

---

## Frontend / Backend（前后端分离）

| 节点 | 字段 | 说明 |
|------|------|------|
| **frontend** | `protocol`, `host`, `port`, `change_origin` | 静态或 SPA 所在地址 |
| **backend** | 同上 | API 服务地址 |
| **backend** | `prefix` | 后端 API 路径前缀，默认 `/api` |
| **backend** | `is_disable_prefix_rewrite` | 是否关闭前缀重写 |

等价环境变量：**`FRONTEND`**、**`BACKEND`**（均为完整 URL）。

---

## 认证：`auth`

| 字段 | 说明 |
|------|------|
| `mode` | `oauth2`、`password`、`none`、`openid` 等 |
| `provider` | 使用哪一条 **`oauth2[].name`**（例如 `doreamon`、`github`、`feishu`、`slack`） |
| `ignore_paths` | 不需要登录的路径前缀列表 |
| `is_ignore_paths_disabled` | 为 `true` 时忽略上述列表（全部需登录） |
| `is_ignore_when_header_authorization_found` | 已有 `Authorization` 头时可跳过部分鉴权逻辑 |
| `allow_usernames` | 仅允许列表中的用户名登录（如收紧 GitHub 登录用户） |

可被 **`AUTH_MODE`**、**`AUTH_PROVIDER`**、**`AUTH_IGNORE_PATHS`** 等覆盖。

---

## OAuth2：`oauth2`

数组；每条对应一种 IdP，登录回调路径一般为：

```txt
http(s)://<connect-host>:<port>/login/<name>/callback
```

| 字段 | 说明 |
|------|------|
| `name` | 与 `auth.provider` 对应；同时也是 URL 中的 `<name>` |
| `client_id` / `client_secret` | IdP 发放的凭证 |
| `redirect_uri` | 必须与 IdP 控制台注册完全一致 |
| `scope` | 可选；依 IdP 要求填写（如 Slack 逗号分隔 scope） |

可同时配置多条（例如同时声明 `doreamon` 与 `slack`，再通过 **`auth.provider`** 选用当前生效的一条）。

---

## 密码登录：`password`

| 字段 | 说明 |
|------|------|
| `mode` | `local`：内置列表；`service`：远程 HTTP 校验 |
| `local` | 数组，每项含 `username`、`password` |
| `service` | 远程校验 URL（`POST`） |

当 **`auth.mode`** 为 `password` 时使用。

---

## 内置业务数据：`services`

内置 HTTP API（如应用信息、当前用户、菜单等）的数据来源可以是远端 HTTP（`service`）或写死在配置里的 `local`。

| 键 | 含义 |
|----|------|
| `services.app` | 应用元数据 |
| `services.user` | 当前用户 |
| `services.menus` | 菜单树 |
| `services.permissions` | 权限列表 |
| `services.users` | 用户列表（结构同菜单项数组场景见源码） |
| `services.open_id` | OpenID 相关信息 |

每个服务常见子字段：

| 子字段 | 说明 |
|--------|------|
| `mode` | `service` 或 `local` |
| `service` | GET 拉取数据的 URL（依具体 API 约定） |
| `local` | 直接写在 YAML 里的对象或数组 |

未在 YAML 中填写时，`ApplyDefault()` 会为 Doreamon 生态补上默认的 `api.zcorky.com` 系列地址（可在自有部署中整体替换）。

---

## 自定义路由：`routes`

将某路径前缀转发到指定后端主机（常用于微服务网关场景）。

| 路径 | 字段 | 说明 |
|------|------|------|
| `routes[].path` | 匹配的前缀 |
| `routes[].backend.service_protocol` | 如 `https` |
| `routes[].backend.service_name` | 主机名 |
| `routes[].backend.service_port` | 端口 |
| `routes[].backend.disable_rewrite` | 是否禁用重写 |
| `routes[].backend.secret_key` | 可选；路由级密钥 |

示例：

```yaml
routes:
  - path: /api/ms/httpbin
    backend:
      service_protocol: https
      service_name: httpbin.zcorky.com
      service_port: 443
```

::: warning 键名大小写
请使用 **snake_case**（如上）。若仓库中旧示例曾出现 camelCase，请以源码 `config` 标签为准。
:::

---

## 内置 API 路径：`built_in_apis`

控制 `/api` 下各 JSON 端点的**路径后缀**（非主机名）。

| 字段 | 默认（摘自 `ApplyDefault`） |
|------|------------------------------|
| `app` | `/app` |
| `user` | `/user` |
| `menus` | `/menus` |
| `permissions` | `/permissions` |
| `users` | `/users` |
| `config` | `/config` |
| `qrcode` | `/qrcode` |
| `login` | `/login` |
| `built_in`（源码字段 `Public`） | `/_` |

也可用环境变量 **`BUILT_IN_APIS_APP`**、`BUILT_IN_APIS_USER` 等单独覆盖。

---

## 菜单项：`menus.local` 等结构

单条菜单常见字段（对应 `MenuItem`）：

| 字段 | 说明 |
|------|------|
| `id` | 可选 |
| `name` | 展示名称 |
| `path` | 前端路径 |
| `icon` | 图标标识 |
| `sort` | 排序 |
| `hidden` | 是否隐藏（YAML 键名为 **`hidden`**） |
| `expanded` | 是否默认展开 |
| `layout` | 布局 |
| `iframe` | 内嵌 iframe 地址 |
| `redirect` | 重定向 |

---

## 完整示例一：单一 upstream + 多条 OAuth2

适合：反向代理一个站点，并在多条 IdP 之间切换（通过 **`auth.provider`**）。

```yaml
port: 8080

secret_key: REPLACE_WITH_LONG_RANDOM_STRING

session_max_age: 86400

upstream:
  protocol: https
  host: httpbin.zcorky.com
  port: 443

oauth2:
  - name: doreamon
    client_id: <DOREAMON_CLIENT_ID>
    client_secret: <DOREAMON_CLIENT_SECRET>
    redirect_uri: http://127.0.0.1:8080/login/doreamon/callback

  - name: slack
    client_id: <SLACK_CLIENT_ID>
    client_secret: <SLACK_CLIENT_SECRET>
    redirect_uri: https://connect.example.com/login/slack/callback
    scope: identity.basic,identity.email,identity.avatar

auth:
  mode: oauth2
  provider: slack   # 当前登录使用 oauth2 名称为 slack 的这一条
```

---

## 完整示例二：前后端 + OAuth2 + password + services + routes

与 `conf/config.full.example` 同级复杂度（占位 URL 请替换为真实环境）。

```yaml
port: 8080

secret_key: go-zoox

session_max_age: 86400

frontend:
  protocol: http
  host: 127.0.0.1
  port: 8000

backend:
  protocol: http
  host: 127.0.0.1
  port: 8001
  prefix: /api
  is_disable_prefix_rewrite: false

oauth2:
  - name: doreamon
    client_id: <CLIENT_ID>
    client_secret: <CLIENT_SECRET>
    redirect_uri: http://127.0.0.1:8080/login/doreamon/callback

password:
  mode: local
  local:
    - username: admin
      password: admin
  service: https://api.example.com # auth.mode=password 且 mode=service 时使用

auth:
  mode: oauth2
  provider: doreamon

services:
  app:
    mode: service
    local:
      name: Lighthouse DNS
      logo: https://avatars.githubusercontent.com/u/7463687?v=4
      description: 轻量级 DNS 服务
    service: https://api.example.com/oauth/app

  user:
    mode: service
    local:
      id: '0x01'
      username: admin
      nickname: 管理员
      avatar: https://avatars.githubusercontent.com/u/7463687?v=4
      permissions:
        - /s/dns
        - /s/docker
    service: https://api.example.com/user

  menus:
    mode: service
    local:
      - name: Home
        path: /i
        icon: appstore
      - name: DNS 管理
        path: /i/dns
        icon: appstore
        layout: header
    service: https://api.example.com/menus

  users:
    mode: service
    service: https://api.example.com/users

  open_id:
    mode: service
    service: https://api.example.com/oauth/app/user/open_id

routes:
  - path: /api/ms/httpbin
    backend:
      service_protocol: https
      service_name: httpbin.zcorky.com
      service_port: 443
```

---

## 完整示例三：密码登录 + 全部 local services

适合离线演示或与 `conf/config.local.yml.example` 对照。

```yaml
port: 8080

secret_key: go-zoox

frontend:
  protocol: http
  host: 127.0.0.1
  port: 8000

backend:
  protocol: http
  host: 127.0.0.1
  port: 8001

password:
  mode: local
  local:
    - username: admin
      password: admin

auth:
  mode: password

services:
  app:
    mode: local
    local:
      name: Lighthouse DNS
      logo: https://avatars.githubusercontent.com/u/7463687?v=4
      description: 轻量级 DNS 服务

  user:
    mode: local
    local:
      id: '0x01'
      username: admin
      nickname: 管理员
      avatar: https://avatars.githubusercontent.com/u/7463687?v=4
      permissions:
        - /s/dns
        - /s/docker

  menus:
    mode: local
    local:
      - name: Home
        path: /i
        icon: appstore
      - name: DNS 管理
        path: /i/dns
        icon: appstore
        layout: header
```

---

## 环境变量与 YAML

以下变量在 **`ApplyDefault()`** 中最常被读取；若同时存在于 YAML，一般以 **环境变量覆盖内存中的配置** 为准（具体顺序以实现为准，建议：敏感信息走 env，结构走 YAML）。

| 变量 | 作用 |
|------|------|
| `PORT` | 端口 |
| `MODE` | 运行模式 |
| `SECRET_KEY` | 会话密钥 |
| `SESSION_MAX_AGE` | 会话时长（秒） |
| `LOG_LEVEL` | 日志 |
| `UPSTREAM` | 单上游 URL |
| `FRONTEND` / `BACKEND` | 前后端 URL |
| `AUTH_*` | 认证模式、忽略路径等 |
| `BUILT_IN_APIS_*` | 内置 API 路径前缀 |

更全列表见 [命令行](./cli.md)。

---

## 延伸阅读

- [部署形态](./architecture.md)：upstream 与 frontend/backend 选型  
- [命令行](./cli.md)：子命令与环境变量速查  
- [Docker 与编排](./docker.md)：镜像与 Compose
