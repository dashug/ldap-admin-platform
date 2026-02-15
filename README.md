# LDAP 管理平台（前后端一体）

基于 Go + Vue 的 **OpenLDAP / Active Directory** 管理后台。本仓库为**前后端合并**结构，一次构建得到一个可执行文件，同时提供 REST API 与 Web 管理界面。

---

## 项目说明与致谢

**本项目是在 [github.com/eryajf/go-ldap-admin](https://github.com/eryajf/go-ldap-admin) 基础上进行的二次开发与整合。**

- 在保留原有功能的前提下，将后端与前端合并到同一仓库，并支持单二进制部署（前端资源嵌入可执行文件）。
- 感谢原项目作者 [eryajf](https://github.com/eryajf) 及所有贡献者，本项目的核心能力与设计均来源于原项目。

---

## 功能说明

### 核心能力

| 模块 | 说明 |
|------|------|
| **LDAP/AD 管理** | 支持 **OpenLDAP** 与 **Active Directory**，统一管理用户、组织架构（OU/组） |
| **用户管理** | 用户的增删改查、批量导入、密码重置、初始密码策略、邮箱后缀等 |
| **分组/组织管理** | 分组（OU）的创建、修改、删除，与 LDAP 目录结构同步 |
| **字段映射** | 自定义数据库字段与 LDAP 属性的映射关系，适配不同目录规范 |
| **LDAP 同步** | 可选将现有 LDAP 中的数据同步到平台数据库，便于在界面中管理 |

### 权限与安全

| 模块 | 说明 |
|------|------|
| **RBAC 权限** | 基于 Casbin 的角色与菜单/API 权限控制，支持多角色、细粒度接口权限 |
| **API Key** | 支持 API Key 认证，便于第三方系统或脚本调用接口 |
| **JWT 登录** | 登录态使用 JWT，可配置过期时间与刷新策略 |
| **操作日志** | 记录关键操作日志，便于审计与排查问题 |
| **限流** | 令牌桶限流，防止接口被滥用 |

### 第三方集成

| 模块 | 说明 |
|------|------|
| **钉钉** | 从钉钉同步组织架构与用户，支持定时同步、部门过滤、离职人员处理 |
| **企业微信** | 从企业微信同步组织与用户，支持定时同步与更新策略 |
| **飞书** | 从飞书同步组织与用户，支持部门列表与定时任务 |
| **Webhook** | 用户/部门创建或同步完成后，可向指定 URL 发送 HTTP 回调 |
| **邮件** | 可选在新建/同步用户后发送通知邮件（需配置 SMTP） |

### 系统与运维

- **多数据库**：支持 **SQLite**（默认，免安装）与 **MySQL**，通过 `config.yml` 切换。
- **日志**：可配置日志级别、路径、轮转大小、保留天数与压缩。
- **运行模式**：支持 `debug` / `release` / `test`，生产环境建议使用 `release`。

---

## 功能界面概览

产品主要包含以下功能模块，登录后可在左侧菜单进入：

- **首页**：仪表盘与概览
- **组织与用户**：用户管理、部门管理、同步字段映射
- **系统**：角色与权限、菜单、接口、系统信息、API 密钥
- **审计**：操作日志

---

## 目录结构

```
ldap-admin-platform/
├── web/                    # 前端 Vue 项目源码
│   ├── src/
│   ├── public/
│   └── package.json
├── public/
│   └── static/
│       └── dist/           # 前端构建产出（make 写入，供 Go embed）
├── config.yml               # 后端配置文件（必改）
├── main.go                  # 后端入口
├── Makefile                 # 一键构建与常用命令
├── Dockerfile               # 镜像构建（需先 make all 生成二进制）
├── deploy/                  # 部署相关脚本与示例
└── README.md
```

---

## 环境要求

- **Go** 1.16+（用于从源码构建）
- **Node.js** 14+、**npm**（用于构建前端）
- 若使用 MySQL：需已有 MySQL 服务；若使用 SQLite，无需额外数据库服务

---

## 快速开始

### 1. 安装前端依赖

```bash
make deps-web
```

### 2. 一键构建（前端 + 后端）

```bash
make all
```

将依次：构建前端 → 拷贝到 `public/static/dist` → 编译 Go（embed 前端）→ 生成可执行文件 **`./go-ldap-admin`**。

### 3. 配置与首次运行

1. 复制或编辑根目录 **`config.yml`**：
   - 修改 **`system.port`**（默认 `8888`）等运行参数；
   - 配置 **`ldap`** 段：LDAP/AD 地址、base-dn、admin-dn、admin-pass、user-dn 等；
   - 若使用 MySQL，配置 **`mysql`** 段，并将 **`database.driver`** 改为 `mysql`、**`database.source`** 改为 DSN；
   - 首次部署时建议 **`system.init-data`** 设为 **`true`**，用于初始化库表与默认数据。

2. 启动服务：

```bash
./go-ldap-admin
# 或
make run
```

3. 浏览器访问 **`http://127.0.0.1:8888`**（端口以 `config.yml` 为准）。  
4. 首次初始化完成后，将 **`system.init-data`** 改为 **`false`**，避免重复初始化。

---

## 开发模式（前后端分离）

- **后端**：在项目根目录执行 `make run` 或 `go run main.go`。
- **前端**：`cd web && npm run dev`，通过前端 dev 代理访问后端 API（需在 `web` 中配置 proxy 指向后端地址）。

生产环境仍建议使用 **`make all`** 得到单二进制部署。

---

## 部署说明

### 生产一键构建（推荐）

在项目根目录执行 **`deploy/production-build.sh`**，会完成：前端依赖安装 → 前端构建 → 后端构建 → 产出到 **`release/ldap-admin-platform-<版本>/`**（含二进制、配置示例、部署说明、systemd 示例）。

```bash
./deploy/production-build.sh                    # 本机架构，版本自动或默认
VERSION=v1.0.2 ./deploy/production-build.sh    # 指定版本
BUILD_OS=linux BUILD_ARCH=amd64 ./deploy/production-build.sh   # 交叉编译 Linux
PACK=1 ./deploy/production-build.sh             # 构建后打 tar.gz 包，便于上传服务器
```

将 `release/` 下对应目录或 tar 包拷贝到目标服务器，按目录内 **生产部署说明.md** 操作即可。

### 方式一：二进制直接部署

1. 在具备 Go + Node 环境的机器上执行：
   ```bash
   make deps-web
   make all
   ```
2. 将生成的 **`go-ldap-admin`** 与 **`config.yml`** 拷贝到目标服务器同一目录。
3. 按需修改 **`config.yml`**（LDAP、数据库、端口、日志等）。
4. 首次运行前确保 **`system.init-data: true`**，启动一次完成初始化后改为 **`false`**。
5. 生产环境建议：
   - **`system.mode`** 设为 **`release`**；
   - **`logs.level`** 设为 **0**（Info）或更高，减少日志量；
   - 使用 systemd/supervisor 等做进程守护与开机自启。

### 方式二：Docker 部署

项目根目录提供 **Dockerfile**，适用于已构建好的二进制：

1. 先在本机构建二进制与前端：
   ```bash
   make all
   ```
2. 将二进制按 Dockerfile 约定放到 **`bin/`** 下（如 `bin/go-ldap-admin_linux_amd64`），并在 Dockerfile 同目录准备 **`config.yml`**。
3. 构建镜像：
   ```bash
   docker build -t ldap-admin-platform .
   ```
4. 运行容器时挂载 **`config.yml`** 与数据目录（如 `data/`），并暴露端口（如 8888）。

注意：当前 Dockerfile 中的 `config.yml` 会被内联修改（如将 `localhost` 改为容器内服务名），若使用外部 LDAP/MySQL，建议通过挂载覆盖 `config.yml`。

### 配置要点摘要

| 配置项 | 说明 |
|--------|------|
| **system.port** | 服务监听端口 |
| **system.init-data** | 是否初始化数据，首次为 true，之后改为 false |
| **system.mode** | debug / release / test，生产用 release |
| **database.driver** | sqlite3 或 mysql |
| **database.source** | SQLite 文件路径或 MySQL DSN |
| **ldap.*** | LDAP/AD 地址、base-dn、管理员账号、user-dn、密码策略等 |
| **jwt.key** | 生产环境务必改为随机强密钥 |
| **logs.path** | 日志目录，需保证进程有写权限 |

---

## Make 命令一览

| 命令 | 说明 |
|------|------|
| **make deps-web** | 安装前端依赖（web/node_modules） |
| **make build-web** | 仅构建前端，产出到 public/static/dist |
| **make build** | 仅构建 Go 后端（需已存在 public/static/dist） |
| **make all** | 先构建前端再构建后端，得到单可执行文件 |
| **make run** | 运行后端（开发时用） |
| **make clean** | 清理前端构建产物与二进制 |

---

## 常见问题

- **构建报错找不到 dist**：先执行 `make build-web` 再 `make build`，或直接 `make all`。
- **前端修改后不生效**：需重新执行 `make build-web` 再重新编译 Go（或 `make all`），因前端被 embed 进二进制。
- **无法连接 LDAP**：检查 `config.yml` 中 `ldap.url`、`ldap.admin-dn`、`ldap.admin-pass` 及网络/防火墙。
- **初始化失败或重复初始化**：确认首次运行后已将 **`system.init-data`** 改为 **`false`**。

---

## 开源协议

请遵循原项目 [go-ldap-admin](https://github.com/eryajf/go-ldap-admin) 所采用的开源协议。  
本项目为在该项目基础上的二次开发，再次感谢原作者与社区贡献者。
