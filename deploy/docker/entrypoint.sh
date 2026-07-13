#!/bin/sh
set -e
# 容器入口：以非 root 的 app 用户运行主进程（容器逃逸加固），同时兼容宿主机 bind mount。
# bind mount 的 /app/data 首次由 Docker 以 root 创建，非 root 的 app 无法写入 sqlite/日志/RSA 私钥，
# 会导致启动崩溃。故这里以 root 入口先把数据目录属主补正为 app，再用 su-exec 降权执行。
# 若容器已被 compose `user:` 指定为非 root 启动，则跳过 chown 直接执行。
if [ "$(id -u)" = "0" ]; then
  mkdir -p /app/data
  chown -R app:app /app/data 2>/dev/null || true
  exec su-exec app "$@"
fi
exec "$@"
