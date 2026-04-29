# 安装

## 使用 Go 安装（推荐）

```bash
go install github.com/go-zoox/connect/cmd/connect@latest
```

确保 `$GOPATH/bin` 或 `$HOME/go/bin` 已加入 `PATH`。安装完成后执行：

```bash
connect --help
connect --version
```

## 使用 zmicro 安装

若未安装 Go，可使用 [zmicro](https://github.com/zcorky/zmicro)：

```bash
curl -o- https://raw.githubusercontent.com/zcorky/zmicro/master/install | bash
zmicro package install connect
```

## Docker 镜像

README 与 Compose 示例中常用镜像形如 `whatwewant/connect:latest`、`whatwewant/connect-doreamon:v1`，具体标签以 [Docker Hub](https://hub.docker.com/) / 团队镜像仓库为准。仓库内提供：

| Dockerfile | 用途简述 |
|------------|----------|
| `Dockerfile` | 通用构建（入口依赖镜像内命令） |
| `Dockerfile.doreamon` | 打包 `connect doreamon` |
| `Dockerfile.github` | GitHub OAuth 场景 |
| `Dockerfile.feishu` | 飞书 OAuth 场景 |
| `Dockerfile.none` | 无认证 |

详见 [Docker 与编排](./docker.md)。

## 验证

```bash
connect version   # 或 connect --version，视构建版本而定
```

完整子命令说明见 [命令行](./cli.md)。
