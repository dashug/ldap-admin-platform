# 功能页面截图说明

README 中的「功能界面预览」需要本目录下的截图文件。请按以下步骤生成：

## 1. 启动本地服务

```bash
# 在项目根目录
make run
# 或
./go-ldap-admin
```

浏览器访问 **http://127.0.0.1:8888**，使用管理员账号登录（默认 admin，密码与 config.yml 中 `ldap.admin-pass` 一致，首次初始化多为 `123456`）。

## 2. 按页面截取并保存

| 功能       | 访问路径（登录后）     | 保存文件名        |
|------------|------------------------|-------------------|
| 登录页     | http://127.0.0.1:8888/login | `login.png`       |
| 首页/仪表盘 | http://127.0.0.1:8888/#/dashboard | `dashboard.png`   |
| 用户管理   | http://127.0.0.1:8888/#/personnel/user | `user.png`        |
| 部门管理   | http://127.0.0.1:8888/#/personnel/group | `group.png`       |
| 同步字段映射 | http://127.0.0.1:8888/#/personnel/fieldRelation | `field-relation.png` |
| 角色与权限 | http://127.0.0.1:8888/#/system/role | `role.png`        |
| 菜单       | http://127.0.0.1:8888/#/system/menu | `menu.png`        |
| 接口       | http://127.0.0.1:8888/#/system/api | `api.png`         |
| 系统信息   | http://127.0.0.1:8888/#/system/info | `system-info.png` |
| API 密钥   | http://127.0.0.1:8888/#/system/apiKey | `api-key.png`     |
| 操作日志   | http://127.0.0.1:8888/#/log/operation-log | `operation-log.png` |

将截图保存到 **本目录**（`docs/screenshots/`），文件名与上表一致即可，README 中会自动引用。
