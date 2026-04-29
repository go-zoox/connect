# 部署形态

Connect 向外暴露一个 HTTP 服务（默认 `:8080`），向内则有两种典型接法。

## 形态 A：单一 Upstream

只配置 **upstream**：所有（通过鉴权后的）请求按路径转发到同一目标（协议 + 主机 + 端口）。

适用于：

- 单体站点、静态 + API 同源的反向代理；
- 已有网关后面的单一 origin。

配置示例见仓库 `conf/config.upstream.yml.example` 与 [示例：上游代理](/examples/upstream)。

YAML 中与 upstream 相关的键为 `upstream.host` / `upstream.port` / `upstream.protocol`（或通过环境变量 `UPSTREAM` 一次性注入 URL）。

## 形态 B：Frontend + Backend

同时配置 **frontend** 与 **backend**：

- **frontend**：一般为静态资源或 SPA 所在地址；
- **backend**：API 服务；默认会以 `/api` 等为前缀做路由与重写（可用配置关闭前缀重写）。

适用于前后端分离、静态与 API 域名或端口不一致的情况。

参考 `conf/config.oauth.yml.example`、`conf/config.full.example`。

## 与环境变量的关系

子命令 `doreamon` / `github` / `feishu` / `none` 会在内存中生成一份 `Config`，等价于手写 YAML 的一部分；你也可以始终使用 `connect serve -c file.yml` 完全显式控制。

无论哪种形态，OAuth 回调路径通常为：

```txt
http(s)://<connect-host>:<port>/login/<oauth2-name>/callback
```

`<oauth2-name>` 与配置里 `oauth2[].name` 一致（如 `doreamon`、`github`、`feishu`）。

## 会话与安全

生产环境务必设置稳定的 `secret_key`，并通过 HTTPS + 正确的 `redirect_uri` 注册到 IdP。关于 Cookie `Secure` 与 `SESSION_SECURE` 等行为，以运行时配置为准（参见源码中 `ApplyDefault` 与 OAuth 中间件）。
