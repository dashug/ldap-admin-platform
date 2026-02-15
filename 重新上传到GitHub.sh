#!/bin/bash
# 删库后重新上传，保证提交记录和 Contributors 里只有 dashug
#
# 第一步：到 GitHub 删除仓库
#   https://github.com/dashug/ldap-admin-platform/settings → 最下方「Delete this repository」
#
# 第二步：打开「系统自带终端」（不要用 Cursor 里的终端），执行：
#   cd /Users/taotongtong/Downloads/ldap-admin-platform
#   bash 重新上传到GitHub.sh

set -e
cd "$(dirname "$0")"

echo "1. 创建 GitHub 空仓库（若已存在会提示，可忽略）..."
gh repo create dashug/ldap-admin-platform --public --description "LDAP 管理平台（前后端一体，基于 go-ldap-admin 二开）" 2>/dev/null || true

echo "2. 初始化 Git..."
git init -b main
git remote add origin https://github.com/dashug/ldap-admin-platform.git 2>/dev/null || true
git remote set-url origin https://github.com/dashug/ldap-admin-platform.git
git config user.name "dashug"
git config user.email "dashug@users.noreply.github.com"

echo "3. 添加文件并提交（作者与提交者均为 dashug）..."
export GIT_AUTHOR_NAME="dashug"
export GIT_AUTHOR_EMAIL="dashug@users.noreply.github.com"
export GIT_COMMITTER_NAME="dashug"
export GIT_COMMITTER_EMAIL="dashug@users.noreply.github.com"
git add -A
git commit -m "chore: initial commit - LDAP 管理平台（前后端一体，基于 go-ldap-admin 二开）"

echo "4. 推送到 GitHub..."
git push -u origin main --force

echo ""
echo "完成。请到 https://github.com/dashug/ldap-admin-platform 查看，提交记录与 Contributors 应只显示 dashug。"
