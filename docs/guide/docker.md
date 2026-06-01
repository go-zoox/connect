# Docker 与编排

仓库提供多套 **Dockerfile**，构建产物一般为 `/bin/connect`，镜像内在 Shell / Alpine 上执行镜像自带的 **entrypoint**（调用对应子命令，例如 `connect doreamon`）。

## 通用 Dockerfile（按需）

| 文件 | 典型用途 |
|------|-----------|
| `Dockerfile` | 默认镜像，`COPY .config.yml` 挂载配置的 Compose |
| `Dockerfile.doreamon` | `ENTRYPOINT` → `connect doreamon` |
| `Dockerfile.github` | GitHub OAuth |
| `Dockerfile.feishu` | 飞书 OAuth |
| `Dockerfile.none` | 无认证 |

实际镜像名与 tag（如 `whatwewant/connect`、`whatwewant/connect-doreamon:v1`）以构建流水线为准。

## docker-compose.yml（挂载 YAML）

仓库根目录示例大致如下——通过挂载 `.config.yml` 注入完整 YAML：

```yaml
version: '3.7'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - '8080:8080'
    volumes:
      - .config.yml:/app/.config.yml
      - ./data:/app/data
```

请先准备好等价于 `conf/config.*.example` 的 YAML 再挂载启动。

## Doreamon 专用 Compose（环境变量）

`docker-compose.doreamon.yml` 示意（节选）：

```yaml
services:
  app:
    image: whatwewant/connect-doreamon:v1
    ports:
      - '8080:8080'
    environment:
      DEBUG: ${DEBUG}
      SESSION_KEY: ${SESSION_KEY}
      UPSTREAM: ${UPSTREAM}
      CLIENT_ID: ${CLIENT_ID}
      CLIENT_SECRET: ${CLIENT_SECRET}
      REDIRECT_URI: ${REDIRECT_URI}
      FRONTEND: ${FRONTEND}
      BACKEND: ${BACKEND}
```

说明：

- 源码里会话密钥优先读取 **`SECRET_KEY`**（见 `defaults.go`）；若镜像文档写 `SESSION_KEY`，请以所用镜像说明为准或在本仓库环境变量表中统一为 `SECRET_KEY`。
- `UPSTREAM` 与 `FRONTEND`+`BACKEND` 二选一（与进程内校验一致）。

## README 中的极简示例

适合快速试跑（镜像 `whatwewant/connect:latest`）：

```yaml
services:
  connect:
    restart: unless-stopped
    image: whatwewant/connect:latest
    ports:
      - '8080:8080'
    environment:
      SECRET_KEY: '666'
      UPSTREAM: https://httpbin.zcorky.com
      CLIENT_ID: <YOUR_DOREAMON_CLIENT_ID>
      CLIENT_SECRET: <YOUR_DOREAMON_CLIENT_SECRET>
      REDIRECT_URI: http://127.0.0.1:8080/login/doreamon/callback
```

## 网络与外部 compose

根目录 `docker-compose.yml` 使用了 `external` 网络 `compose-ingress`，用于与其它栈互联；单机试用时可改为默认 bridge 或自建 network。

更多场景化片段见 [示例目录](/examples/)。
