# 快速开始

下文假设你已按 [安装](./installation.md) 装好 `connect`。

## 1. 选定模式

| 场景 | 命令 | 说明 |
|------|------|------|
| 已有完整 YAML | `connect serve -c /path/to/config.yaml` | 适合自定义 upstream / 前后端拆分 |
| Doreamon | `connect doreamon` | OAuth2 Provider 为 doreamon，常用环境变量注入 |
| GitHub | `connect github` | GitHub OAuth |
| 飞书 | `connect feishu` | 飞书 OAuth |
| 无认证 | `connect none` | 只做代理，不做登录 |

便捷子命令会从环境变量组装配置；完整字段控制请用 **serve + YAML**，见 [配置](./config.md)。

## 2. Doreamon 一键示例（环境变量）

```bash
export SECRET_KEY=your-secret
export UPSTREAM=https://httpbin.zcorky.com   # 或同时设置 FRONTEND / BACKEND，见架构说明

export CLIENT_ID=<YOUR_CLIENT_ID>
export CLIENT_SECRET=<YOUR_CLIENT_SECRET>
export REDIRECT_URI=http://127.0.0.1:8080/login/doreamon/callback

connect doreamon
```

浏览器访问 `http://127.0.0.1:8080`（默认端口 `8080`，可通过 `PORT` 修改）。

## 3. GitHub OAuth 示例

```bash
export SECRET_KEY=your-secret
export UPSTREAM=https://httpbin.zcorky.com

export CLIENT_ID=<GITHUB_OAUTH_APP_CLIENT_ID>
export CLIENT_SECRET=<GITHUB_OAUTH_APP_SECRET>
export REDIRECT_URI=http://127.0.0.1:8080/login/github/callback

connect github
```

## 4. 下一步

- 弄清 **上游一种还是前后端两种**：读 [部署形态](./architecture.md)。
- 罗列所有子命令与变量：读 [命令行](./cli.md)。
- 编辑 YAML、与环境变量对照：读 [配置](./config.md)。
- 可复制粘贴的 Compose / 片段：读 [示例目录](/examples/)。
