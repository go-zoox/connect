# GitHub OAuth

使用 GitHub OAuth App 作为登录提供方：`auth.provider` 为 `github`，回调路径 `/login/github/callback`。

## CLI

```bash
export SECRET_KEY=$(openssl rand -hex 16)
export UPSTREAM=https://httpbin.zcorky.com

export CLIENT_ID=<GITHUB_APP_CLIENT_ID>
export CLIENT_SECRET=<GITHUB_APP_CLIENT_SECRET>
export REDIRECT_URI=http://127.0.0.1:8080/login/github/callback

connect github
```

## GitHub 控制台配置要点

1. Authorization callback URL 填写与 `REDIRECT_URI` 相同。
2. 若仅允许部分 GitHub 用户，可在配置里设置 **`auth.allow_usernames`**（参见 [配置](../guide/config.md)）。

## YAML 片段

参考 `conf/config.oauth.yml.example`，将 `oauth2` 中 `name` 改为 `github` 并填入 GitHub 发放的 Client ID/Secret。

```yaml
oauth2:
  - name: github
    client_id: ${GITHUB_CLIENT_ID}
    client_secret: ${GITHUB_CLIENT_SECRET}
    redirect_uri: https://connect.example.com/login/github/callback
```
