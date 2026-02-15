# Go LDAP Admin UI (Dashug Fork)

前端管理界面，配合后端实现目录与第三方平台的可视化配置。

## 致谢

本项目基于 [example/go-ldap-admin-ui](https://github.com/example/go-ldap-admin-ui) 二次开发，感谢原作者与全部贡献者。

## 核心功能

- 用户/分组管理界面
- 目录类型展示（OpenLDAP / AD）
- 目录快速配置入口（地址、DN、类型、同步开关）
- 平台对接向导（钉钉/企微/飞书）
- 测试连接 + 保存配置

## 当前前端对接接口

- `GET /api/base/config`
- `POST /api/base/directoryConfig`
- `POST /api/base/thirdPartyConfig`
- `POST /api/base/thirdPartyConfig/test`

## 页面配置流程

### 目录快速配置

页面入口：`人员管理 -> 用户管理 -> 目录快速配置`

操作步骤：

1. 选择目录类型：`OpenLDAP` / `Windows AD`。
2. 填写地址和 DN 参数。
3. 点击保存。

字段说明（常用）：

- `LDAP地址`：如 `ldap://10.0.0.10:389` 或 `ldaps://ad.example.com:636`
- `Base DN`：如 `dc=example,dc=com`
- `管理员 DN`：用于服务端连接目录
- `用户 OU DN`：新建用户默认写入位置

### 平台对接向导

页面入口：`人员管理 -> 用户管理 -> 平台对接向导`

操作步骤：

1. 选择平台标签页（钉钉 / 企业微信 / 飞书）。
2. 填写平台凭证。
3. 点击 `测试连接`。
4. 成功后点击 `保存`。

建议先完成目录配置，再做平台对接与同步。

## 图文功能说明

### 1. 功能总览图（首页 / 用户管理 / 分组管理）

![首页](https://example.com/img/image_20220724_165545.png)
![用户管理](https://example.com/img/image_20220724_165623.png)
![分组管理](https://example.com/img/image_20220724_165701.png)

### 2. 目录快速配置（步骤图 1-2-3）

```mermaid
flowchart LR
  A["用户管理页面"] --> B["点击 目录快速配置"]
  B --> C["选择 OpenLDAP/AD"]
  C --> D["填写地址与DN"]
  D --> E["保存配置"]
```

### 3. 平台对接向导（钉钉 / 企微 / 飞书）

```mermaid
flowchart LR
  A["点击 平台对接向导"] --> B["选择平台"]
  B --> C["填写凭证"]
  C --> D["测试连接"]
  D --> E["保存"]
```

### 4. 测试连接与保存结果示例

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "platform": "wecom",
    "ok": true
  }
}
```

### 5. 常见错误提示（可选）

- `测试连接失败`：先检查凭证字段是否填错。
- `保存后不同步`：确认对应平台 `启用同步` 是否打开。
- `AD 配置失败`：优先使用 `ldaps://` 地址。

## 维护仓库

- 后端：<https://github.com/dashug/go-ldap-admin>
- 前端：<https://github.com/dashug/go-ldap-admin-ui>

## 部署说明

### 一键全栈部署（推荐）

请使用后端仓库提供的一键部署包（包含统一入口 + 后端，支持 OpenLDAP / AD）：

- 仓库：<https://github.com/dashug/go-ldap-admin>
- 目录：`deploy/full-stack`

执行：

```bash
cd deploy/full-stack
./setup.sh
```

访问：`http://localhost:8080`

### 本地开发运行

```bash
npm install
npm run dev
```

默认开发地址通常为：

- `http://localhost:9528`

### 生产构建

```bash
npm install
npm run build:prod
```

构建产物在 `dist/` 目录，可交给 Nginx/静态文件服务托管。

### Docker 快速部署（前端）

项目内提供了前端容器 `Dockerfile` 与 Nginx 配置 `default.conf`。  
在仓库根目录执行：

```bash
docker build -t go-ldap-admin-ui:latest .
docker run -d --name go-ldap-admin-ui -p 80:80 go-ldap-admin-ui:latest
```

说明：

- 前端默认通过 `/api/` 反向代理到 `http://go-ldap-admin-server:8888`。
- 若你的后端地址不同，请调整 `default.conf` 中 `proxy_pass` 后再构建镜像。
