# 上游代理（单站点）

仅代理到一个 HTTP(S) 上游时，使用 **`upstream`** 块或环境变量 **`UPSTREAM`**。

## YAML 最小示例

对应仓库 `conf/config.upstream.yml.example`：

```yaml
port: 8080
secret_key: go-zoox
session_max_age: 86400

upstream:
  host: 127.0.0.1
  port: 8000

oauth2:
  - name: doreamon
    client_id: <CLIENT_ID>
    client_secret: <CLIENT_SECRET>
    redirect_uri: http://127.0.0.1:8080/login/doreamon/callback
```

启动：

```bash
connect serve -c ./conf/config.upstream.yml
```

## 环境变量等价写法

```bash
export UPSTREAM=http://127.0.0.1:8000
export CLIENT_ID=...
export CLIENT_SECRET=...
export REDIRECT_URI=http://127.0.0.1:8080/login/doreamon/callback

connect doreamon
```

区别：YAML 可同时声明多条 `oauth2`、自定义 `routes`；纯 env 适合容器注入。
