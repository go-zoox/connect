# 简介

**GoZoox Connect** 运行在浏览器与业务系统之间：访客先经过 Connect 完成认证（或按策略放行），再访问背后的前端或上游 HTTP 服务。适合「不想自建全套账号体系，又要把登录收敛到一个入口」的场景。

## 能做什么

| 能力 | 说明 |
|------|------|
| **OAuth2 登录** | 配置 `oauth2` 列表与 `auth.provider`，走 `/login/<name>/callback` 回调并完成会话。 |
| **上游代理** | 将已登录请求转发到单个 **upstream**，或拆成 **frontend**（静态）+ **backend**（API）。 |
| **无认证模式** | `auth.mode: none` 或 `connect none`，用于内网穿透、调试或与外层网关组合。 |
| **内置 JSON API** | 如 `/api/app`、`/api/user`、`/api/menus` 等（路径可通过配置改名）。 |

## 与 gzproxy / gzauth 的关系

Connect 侧重「认证连接器」与一体化路由；同类组件可参考仓库 README 中的 [gzproxy](https://github.com/go-zoox/gzproxy)、[gzauth](https://github.com/go-zoox/gzauth)。

## 下一步

- 尚未安装：先看 [安装](./installation.md)。
- 想最快跑起来：[快速开始](./quick-start.md)。
- 需要理解路由与上游：[部署形态](./architecture.md)。
- 字段级说明与完整 YAML 示例：[配置](./config.md)。
