#!/usr/bin/env bash
# 生产环境一键构建脚本
# 用法:
#   ./deploy/production-build.sh              # 本机架构，版本从 git tag 或默认 v1.0.1
#   VERSION=v1.0.2 ./deploy/production-build.sh
#   BUILD_OS=linux BUILD_ARCH=amd64 ./deploy/production-build.sh   # 交叉编译 Linux
#   PACK=1 ./deploy/production-build.sh      # 构建后打 tar.gz 包

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT"

# 版本：环境变量 > git tag > 默认
VERSION="${VERSION:-$(git describe --tags --always 2>/dev/null | sed 's/^v//' || echo '1.0.1')}"
# 若 git describe 带 v 前缀，统一去掉便于文件名
VERSION="${VERSION#v}"

# 可选：交叉编译
BUILD_OS="${BUILD_OS:-}"
BUILD_ARCH="${BUILD_ARCH:-}"
# 是否打 tar 包
PACK="${PACK:-0}"

OUT_DIR="$ROOT/release/ldap-admin-platform-${VERSION}"
BINARY_NAME="go-ldap-admin"

echo "== LDAP 管理平台 - 生产构建 =="
echo "  版本: $VERSION"
echo "  输出: $OUT_DIR"
[ -n "$BUILD_OS" ] && echo "  目标: ${BUILD_OS}/${BUILD_ARCH}"
echo ""

# 1. 前端依赖
if [ ! -d "web/node_modules" ]; then
  echo "[1/4] 安装前端依赖..."
  (cd web && npm install)
else
  echo "[1/4] 前端依赖已存在，跳过 install"
fi

# 2. 前端构建
echo "[2/4] 构建前端..."
(cd web && npm run build:prod)
rm -rf public/static/dist
cp -r web/dist public/static/

# 3. 后端构建
echo "[3/4] 构建后端..."
export LDFLAGS="-X 'github.com/dashug/ldap-admin-platform/public/version.Version=v${VERSION}' \
  -X 'github.com/dashug/ldap-admin-platform/public/version.GitCommit=$(git rev-parse --short=7 HEAD 2>/dev/null || echo unknown)' \
  -X 'github.com/dashug/ldap-admin-platform/public/version.BuildTime=$(date "+%Y-%m-%d %H:%M:%S")' \
  -X 'github.com/dashug/ldap-admin-platform/public/version.GoVersion=$(go version | awk "{print \$3}")'"

if [ -n "$BUILD_OS" ] && [ -n "$BUILD_ARCH" ]; then
  export GOOS="$BUILD_OS"
  export GOARCH="$BUILD_ARCH"
  BINARY_SUFFIX="_${BUILD_OS}_${BUILD_ARCH}"
  go build -ldflags "$LDFLAGS" -o "${BINARY_NAME}${BINARY_SUFFIX}" main.go
  BINARY_OUT="${BINARY_NAME}${BINARY_SUFFIX}"
else
  go build -ldflags "$LDFLAGS" -o "$BINARY_NAME" main.go
  BINARY_OUT="$BINARY_NAME"
fi

# 4. 打包发布目录
echo "[4/4] 整理发布目录..."
mkdir -p "$OUT_DIR"
cp "$ROOT/$BINARY_OUT" "$OUT_DIR/$BINARY_NAME"
cp "$ROOT/config.yml" "$OUT_DIR/config.yml.example"

# 生产部署说明
cat > "$OUT_DIR/生产部署说明.md" <<DEPLOY
# LDAP 管理平台 - 生产部署说明

## 版本

- **v${VERSION}**
- 二进制: \`$BINARY_NAME\`

## 步骤

1. **上传**  
   将本目录下 \`$BINARY_NAME\`、\`config.yml.example\` 拷贝到目标服务器（如 \`/opt/ldap-admin\`）。

2. **配置**  
   - 将 \`config.yml.example\` 复制为 \`config.yml\` 并按实际环境修改：
     - \`system.port\`：监听端口（默认 8888）
     - \`system.mode\`：改为 \`release\`
     - \`system.init-data\`：首次部署设为 \`true\`，初始化完成后改为 \`false\`
     - \`ldap\`：LDAP/AD 地址、base-dn、admin-dn、admin-pass、user-dn 等
     - 若用 MySQL，修改 \`database.driver\`、\`database.source\` 及 \`mysql\` 段
   - \`logs.level\` 生产建议 \`0\`（Info）。

3. **首次运行**  
   \`system.init-data: true\` 时执行一次：
   \`\`\`
   ./$BINARY_NAME
   \`\`\`
   访问 http://服务器IP:8888，使用 **admin** / \`config.yml\` 中 \`ldap.admin-pass\`（默认 123456）登录。  
   登录成功后，将 \`system.init-data\` 改为 \`false\` 并重启。

4. **进程守护（推荐）**  
   使用 systemd 时，可创建 \`/etc/systemd/system/go-ldap-admin.service\`，内容见同目录 \`go-ldap-admin.service.example\`，然后：
   \`\`\`
   sudo systemctl daemon-reload
   sudo systemctl enable go-ldap-admin
   sudo systemctl start go-ldap-admin
   \`\`\`

## 目录建议

- 二进制与 \`config.yml\` 放在同一目录；程序会在当前工作目录下创建 \`data/\`（SQLite、日志等），请保证该目录可写。
- 或使用 \`systemd\` 的 \`WorkingDirectory\` 指定工作目录。
DEPLOY

# systemd 示例
cat > "$OUT_DIR/go-ldap-admin.service.example" <<SYSTEMD
[Unit]
Description=LDAP Admin Platform
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/ldap-admin
ExecStart=/opt/ldap-admin/go-ldap-admin
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
SYSTEMD

echo ""
echo "构建完成: $OUT_DIR"
echo "  - $BINARY_NAME"
echo "  - config.yml.example"
echo "  - 生产部署说明.md"
echo "  - go-ldap-admin.service.example"
echo ""

if [ "$PACK" = "1" ] || [ "$PACK" = "yes" ]; then
  TARBALL="$ROOT/release/ldap-admin-platform-v${VERSION}.tar.gz"
  (cd "$ROOT/release" && tar czf "$TARBALL" "ldap-admin-platform-${VERSION}")
  echo "已打包: $TARBALL"
  echo "  解压: tar xzf ldap-admin-platform-v${VERSION}.tar.gz && cd ldap-admin-platform-${VERSION}"
  echo ""
fi
