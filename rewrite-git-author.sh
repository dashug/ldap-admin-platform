#!/bin/bash
# 在「系统自带终端」运行此脚本，把历史提交的作者和提交者都改成 dashug，
# 避免 GitHub 提交记录 / Contributors 里出现 Cursor Agent。
#
# 用法：
#   1. 打开 终端.app（或 iTerm），不要用 Cursor 里的终端
#   2. cd /Users/taotongtong/Downloads/ldap-admin-platform
#   3. bash rewrite-git-author.sh
#   4. 在同一终端里执行：git push origin main --force
#
# 说明：GitHub 的 Contributors 有时会缓存，更新后可能要过几小时才只显示 dashug。
#       之后尽量用系统终端做 git push，这样动态里也不会显示「由 Cursor 推送」。

set -e
git filter-branch -f --env-filter '
export GIT_AUTHOR_NAME="dashug"
export GIT_AUTHOR_EMAIL="dashug@users.noreply.github.com"
export GIT_COMMITTER_NAME="dashug"
export GIT_COMMITTER_EMAIL="dashug@users.noreply.github.com"
' --tag-name-filter cat -- main

echo ""
echo "已重写完成。请在本机终端执行："
echo "  git push origin main --force"
echo ""
echo "务必在「系统终端」里执行 push，不要用 Cursor 里的终端。"
