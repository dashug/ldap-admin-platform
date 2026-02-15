#!/bin/bash
# 请在「系统自带终端」执行此脚本，避免提交记录和 Contributors 出现 Cursor Agent
# 用法：打开 终端.app，cd 到项目目录，执行：bash 提交并推送到GitHub-请在系统终端执行.sh

set -e
cd "$(dirname "$0")"

export GIT_AUTHOR_NAME="dashug"
export GIT_AUTHOR_EMAIL="dashug@users.noreply.github.com"
export GIT_COMMITTER_NAME="dashug"
export GIT_COMMITTER_EMAIL="dashug@users.noreply.github.com"

git add -A
git status
git commit -m "fix: 统一为新仓库地址；README 面向使用者；移除发布小节；Contributing 链接修正"
git push origin main
echo "已推送，提交记录与 Contributors 应为 dashug。"
