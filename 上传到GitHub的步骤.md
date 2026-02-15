# 上传到 GitHub 的步骤

本地已完成：Git 初始化、首次提交、远程地址已设为 `https://github.com/dashug/ldap-admin-platform.git`。

你只需要在 GitHub 上**新建一个空仓库**，然后在本机执行一次推送即可。

---

## 第一步：在 GitHub 上创建新仓库

1. 打开：**https://github.com/new**
2. **Repository name** 填：`ldap-admin-platform`
3. **Description** 可选填：`LDAP 管理平台（前后端一体，基于 go-ldap-admin 二开）`
4. 选择 **Public**
5. **不要**勾选 “Add a README file”、“Add .gitignore”等，保持仓库为空
6. 点击 **Create repository**

---

## 第二步：在本地推送代码

在终端进入项目目录后执行：

```bash
cd /Users/taotongtong/Downloads/ldap-admin-platform
git push -u origin main
```

若你的默认分支是 `master` 而不是 `main`，先执行：

```bash
git branch -M main
git push -u origin main
```

按提示输入 GitHub 用户名和密码（或 Personal Access Token）即可完成上传。

---

## 推送完成后

仓库地址为：**https://github.com/dashug/ldap-admin-platform**

之后修改代码可这样提交并推送：

```bash
git add .
git commit -m "你的提交说明"
git push
```
