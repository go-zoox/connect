# 飞书 OAuth

飞书开放平台创建应用后，拿到 App ID / App Secret，回调一般为 `/login/feishu/callback`（与配置的 `oauth2[].name` 一致）。

## CLI

```bash
export SECRET_KEY=$(openssl rand -hex 16)
export UPSTREAM=https://your-internal-gateway.example.com

export CLIENT_ID=<FEISHU_APP_ID>
export CLIENT_SECRET=<FEISHU_APP_SECRET>
export REDIRECT_URI=https://connect.example.com/login/feishu/callback

connect feishu
```

## 注意事项

- 飞书侧「重定向 URL」需与白名单一致。
- 与其它 OAuth 模式相同：必须提供 **`UPSTREAM`** 或 **`FRONTEND`+`BACKEND`**。

更细的 Scope、租户开通流程请参阅 [飞书开放平台文档](https://open.feishu.cn/)（外链可能变更，以官网为准）。
