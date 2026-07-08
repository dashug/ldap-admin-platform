#!/usr/bin/env bash
#
# purge-git-secrets.sh — 从 git 历史中彻底抹除曾提交的敏感文件（RSA 私钥/公钥）。
#
# 背景：
#   config/go-ldap-admin-priv.pem 曾被提交进历史（gitleaks 定位于 commit 10141e92），
#   等同于一把「公开的」私钥。当前代码已改为【运行期生成/注入】RSA 密钥（config/rsa_key.go），
#   仓库不再需要任何 .pem，因此应把历史中的这把私钥抹除。
#
# ⚠️  该操作会【重写整个 git 历史】：所有 commit hash 变化、必须 force-push、
#     所有协作者都要重新 clone。执行前务必：
#       1) 备份（脚本会自动做一个镜像备份）
#       2) 通知所有协作者
#       3) 事后仍需在目录服务侧【轮换】该密钥 —— 历史一旦被 push/clone 过即视为已泄露。
#
# 用法：
#   bash scripts/purge-git-secrets.sh
#
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"

# 要从【全部历史】中移除的路径（可按需增删）
PATHS=(
  "config/go-ldap-admin-priv.pem"
  "config/go-ldap-admin-pub.pem"
)

echo "==> 目标仓库: $REPO_ROOT"
echo "==> 当前分支: $(git rev-parse --abbrev-ref HEAD)"
echo "==> 将从【全部历史】中移除："
printf '      - %s\n' "${PATHS[@]}"
echo

if [ -n "$(git status --porcelain)" ]; then
  echo "✋ 工作区有未提交改动，请先提交或 stash 后再运行。"
  exit 1
fi

read -r -p "已理解此操作会重写历史、已通知协作者？输入 yes 继续: " ok
[ "${ok:-}" = "yes" ] || { echo "已取消。"; exit 1; }

# 1) 镜像备份（可回滚）
BACKUP="../$(basename "$REPO_ROOT")-backup-$(date +%Y%m%d-%H%M%S).git"
echo "==> 创建镜像备份: $BACKUP"
git clone --mirror . "$BACKUP" >/dev/null
echo "    备份完成。如需回滚：从 $BACKUP 重新 clone。"
echo

# 2) 重写历史：优先 git-filter-repo，其次同名 python 模块，最后给出 BFG 提示
build_args() { local a=(); for p in "${PATHS[@]}"; do a+=(--path "$p"); done; printf '%s\n' "${a[@]}"; }

if command -v git-filter-repo >/dev/null 2>&1; then
  echo "==> 使用 git-filter-repo 移除文件…"
  mapfile -t ARGS < <(build_args)
  git filter-repo --force --invert-paths "${ARGS[@]}"
elif python3 -c "import git_filter_repo" >/dev/null 2>&1; then
  echo "==> 使用 python -m git_filter_repo 移除文件…"
  mapfile -t ARGS < <(build_args)
  python3 -m git_filter_repo --force --invert-paths "${ARGS[@]}"
else
  cat <<'MSG'
✋ 未检测到 git-filter-repo。请先安装其一后重跑本脚本：
      brew install git-filter-repo          # macOS
      pipx install git-filter-repo          # 跨平台
      pip3 install git-filter-repo

  或改用 BFG（手动执行）：
      brew install bfg
      bfg --delete-files go-ldap-admin-priv.pem
      bfg --delete-files go-ldap-admin-pub.pem
      git reflog expire --expire=now --all && git gc --prune=now --aggressive
MSG
  exit 1
fi

# 3) 清理 reflog 与松散对象
echo "==> 清理引用与对象…"
git reflog expire --expire=now --all || true
git gc --prune=now --aggressive || true

# 4) 校验：历史中应已无该文件
echo
echo "==> 校验历史是否还残留目标文件："
leftover=0
for p in "${PATHS[@]}"; do
  if git log --all --oneline -- "$p" | grep -q .; then
    echo "    ✗ 仍存在: $p"; leftover=1
  else
    echo "    ✓ 已清除: $p"
  fi
done
[ "$leftover" -eq 0 ] && echo "    历史清理完成。"

cat <<'NEXT'

==> 后续手动步骤（确认无误后执行）：
      # git-filter-repo 会移除 origin，需要重新添加你的远程地址：
      git remote add origin <你的远程仓库地址>
      git push --force --all  origin
      git push --force --tags origin

⚠️  force-push 后：所有协作者必须【重新 clone】（旧本地副本仍含泄露密钥）。
⚠️  该私钥应视为已泄露：本项目已改为运行期生成 RSA 密钥，生产请用
    RSA_PRIVATE_KEY(_FILE) 环境变量注入各自独立的密钥，无需再提交任何 .pem。
NEXT
