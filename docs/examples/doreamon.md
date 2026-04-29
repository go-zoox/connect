# Doreamon 模式

适合已开通 **Doreamon**（或兼容同一 OAuth 契约的服务）的应用：`oauth2[].name` 为 `doreamon`，回调形如 `/login/doreamon/callback`。

## 环境变量启动

```bash
export SECRET_KEY=$(openssl rand -hex 16)
export UPSTREAM=https://your-demo-site.example.com

export CLIENT_ID=<CLIENT_ID>
export CLIENT_SECRET=<CLIENT_SECRET>
export REDIRECT_URI=https://connect.example.com/login/doreamon/callback

connect doreamon
```

若不用 `UPSTREAM`，可同时指定：

```bash
export FRONTEND=http://127.0.0.1:8000
export BACKEND=http://127.0.0.1:8001
```

## Docker Compose 片段

```yaml
services:
  connect:
    image: whatwewant/connect-doreamon:v1
    ports:
      - '8080:8080'
    environment:
      SECRET_KEY: ${SECRET_KEY}
      UPSTREAM: ${UPSTREAM}
      CLIENT_ID: ${CLIENT_ID}
      CLIENT_SECRET: ${CLIENT_SECRET}
      REDIRECT_URI: ${REDIRECT_URI}
```

## IdP 侧检查清单

- Redirect URI 与 `REDIRECT_URI` 完全一致（协议、端口、路径）。
- Client ID / Secret 与 Connect 环境变量一致。
- 生产环境使用 HTTPS，并与会话 Cookie 策略一致。
