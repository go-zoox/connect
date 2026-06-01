---
layout: home

hero:
  name: GoZoox Connect
  text: 认证与上游一站接入
  tagline: 在浏览器用户与后端服务之间插入登录层：OAuth2、Doreamon、GitHub、飞书，以及上游反向代理与内置 API。
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/quick-start
    - theme: alt
      text: 浏览示例
      link: /examples/
    - theme: alt
      text: 源码
      link: https://github.com/go-zoox/connect

features:
  - icon: 🔐
    title: OAuth2 优先
    details: 统一登录回调与会话；支持 Doreamon、GitHub、飞书等 Provider，可按配置扩展多条 oauth2 条目。
  - icon: 🧩
    title: 部署形态灵活
    details: 既可代理单一 upstream，也可拆成「前端静态 + 后端 API」双路由；无认证模式用于纯穿透调试。
  - icon: 🏗️
    title: Doreamon 生态
    details: 与 api.zcorky.com 侧应用 / 用户 / 菜单 / 权限等服务联动（可切换为自有服务端点）。
  - icon: 📦
    title: 镜像与 Compose
    details: 提供多种 Dockerfile（none、doreamon、github、feishu 等），便于不同场景打包与编排。
  - icon: ⚙️
    title: 配置双通道
    details: YAML 配置 + 环境变量覆盖（端口、密钥、上游、OAuth Client、内置 API 路径等）。
  - icon: 🔌
    title: 内置 HTTP API
    details: 聚合应用信息、当前用户、菜单、权限等 JSON API，供控制台或前端动态加载。
---
