# GoZoox Connect - The Lighweight, Powerful Auth Connect

[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-zoox/connect)](https://pkg.go.dev/github.com/go-zoox/connect)
[![Build Status](https://github.com/go-zoox/connect/actions/workflows/release.yml/badge.svg?branch=master)](https://github.com/go-zoox/connect/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-zoox/connect)](https://goreportcard.com/report/github.com/go-zoox/connect)
[![Coverage Status](https://coveralls.io/repos/github/go-zoox/connect/badge.svg?branch=master)](https://coveralls.io/github/go-zoox/connect?branch=master)
[![GitHub issues](https://img.shields.io/github/issues/go-zoox/connect.svg)](https://github.com/go-zoox/connect/issues)
[![Release](https://img.shields.io/github/tag/go-zoox/connect.svg?label=Release)](https://github.com/go-zoox/connect/tags)


GoZoox Connect 是一个 Auth Connect(or)，帮助你无痛接入认证，支持多种认证方式，优先支持 OAuth2，特别支持 Doreamon (类似 Auth0，支持统一应用、用户、权限、配置等)。它可以使用Docker进行部署，并支持私有化部署。此外，connect完全开源。

## Installation
To install the package, run:

```bash
# with go
go install github.com/go-zoox/connect/cmd/connect@latest
```

if you dont have go installed, you can use the install script (zmicro package manager):

```bash
curl -o- https://raw.githubusercontent.com/zcorky/zmicro/master/install | bash

zmicro package install connect
```

## Features

- [x] 支持 Oauth2 认证
  - [x] 支持 Doreamon 登录
  - [x] 支持 GitHub 登录
  - [x] 支持飞书登录
- [ ] 支持 Basic Auth 认证
- [ ] 支持 BearToken 认证
- [x] 使用Docker容器化部署
- [x] 支持私有化部署

## Quick Start

### Using Command Line

#### Full Mode
```bash
connect serve -c /path/to/config.yaml
```

#### Doreamon Mode
```bash
export SECRET_KEY=666
export UPSTREAM=https://httpbin.zcorky.com
#
export CLIENT_ID=<YOUR_DOREAMON_CLIENT_ID>
export CLIENT_SECRET=<YOUR_DOREAMON_CLIENT_SECRET>
export REDIRECT_URI=http://127.0.0.1:8080/login/doreamon/callback

connect doreamon
```

#### OAuth2 Mode - GitHub
```bash
export SECRET_KEY=666
export UPSTREAM=https://httpbin.zcorky.com
#
export CLIENT_ID=<YOUR_GITHUB_CLIENT_ID>
export CLIENT_SECRET=<YOUR_GITHUB_CLIENT_SECRET>
export REDIRECT_URI=http://127.0.0.1:8080/login/github/callback

connect github
```

#### OAuth2 Mode - Feishu
```bash
export SECRET_KEY=666
export UPSTREAM=https://httpbin.zcorky.com
#
export CLIENT_ID=<YOUR_GITHUB_CLIENT_ID>
export CLIENT_SECRET=<YOUR_GITHUB_CLIENT_SECRET>
export REDIRECT_URI=http://127.0.0.1:8080/login/github/callback

connect feishu
```

### Using Docker / Docker Compose

1. create `docker-compose.yml`，这里使用 Doreamon Mode:

```yaml
# 使用 basic auth
services:
  connect:
    restart: unless-stopped
    image: whatwewant/connect:latest
    ports:
      - 8080:8080
    environment:
      SECRET_KEY: 666
      UPSTREAM: https://httpbin.zcorky.com
      CLIENT_ID: <YOUR_DOREAMON_CLIENT_ID>
      CLIENT_SECRET: <YOUR_DOREAMON_CLIENT_SECRET>
      REDIRECT_URI: http://127.0.0.1:8080/login/doreamon/callback
```

2. 启动容器：

```bash
$ docker-compose up -d
```

## Usage

```bash
NAME:
   connect - The Connector

USAGE:
   connect [global options] command [command options] [arguments...]

VERSION:
   1.16.3

DESCRIPTION:
   Connect between auth with apps/services

COMMANDS:
   serve     Start Connect Server
   none      Start Connect Server using None Auth
   doreamon  Start Connect Server using Doreamon
   github    Start Connect Server using GitHub
   feishu    Start Connect Server using Feishu
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Refrence
* [go-zoox/gzproxy](https://github.com/go-zoox/gzproxy) - Easy to proxy with your http server or any another upstream. Built in supports Basic Auth, Bearer Toke, OAuth2 (GitHub, Feishu, Doreamon, etc.)
* [go-zoox/gzauth](https://github.com/go-zoox/gzauth) - Simple Your Auth for Web Service

## 贡献

欢迎您参与贡献connect！请参阅 [CONTRIBUTING.md](./CONTRIBUTING.md) 文件了解更多信息。

## 许可证

connect采用MIT许可证。请参阅 [LICENSE.md](./LICENSE.md) 文件了解详细信息。
