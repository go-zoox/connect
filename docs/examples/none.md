# 无认证穿透

`auth.mode` 为 **`none`** 时不挂载 OAuth /登录中间件，只做反向代理，适用于：

- 上游已有网关鉴权；
- 内网调试；
- 临时暴露静态站点。

## CLI

```bash
export UPSTREAM=http://127.0.0.1:3000
# 或使用 FRONTEND + BACKEND

connect none
```

同样必须满足：**upstream** 或 **frontend + backend** 其一。

## 风险提示

无认证模式不会对访客做强校验，请勿直接暴露在公网敏感资产前。
