# 多阶段自包含构建：一步得到「前端已 embed 进 Go 单二进制」的最小运行镜像。
#   docker build -t go-ldap-admin .
# 旧版（仅打包预编译二进制）见 deploy/Dockerfile.prebuilt。
# syntax=docker/dockerfile:1

# 多架构说明：前端与后端均在【构建平台】上运行（$BUILDPLATFORM，不走 QEMU 模拟），
# 前端产物与架构无关，Go 用 CGO=0 直接交叉编译到目标架构（$TARGETOS/$TARGETARCH），
# 因此 linux/amd64 + linux/arm64 多架构构建既快又稳。

# ---------- 1) 构建前端（Vue 3 + Vite，需 Node 18+；架构无关） ----------
FROM --platform=$BUILDPLATFORM node:18-alpine AS web
WORKDIR /web
# 先装依赖（利用层缓存）：有 lock 用 npm ci，否则回退 npm install
COPY web/package.json web/package-lock.json* ./
RUN npm ci || npm install
COPY web/ ./
RUN npm run build:prod

# ---------- 2) 编译后端（前端 embed 进单二进制；纯 Go sqlite，CGO=0 交叉编译） ----------
# Go 构建镜像跟随最新 1.26.x 补丁版：修复 govulncheck 报告的 stdlib 漏洞
# （crypto/x509、net/textproto、html/template、net/mail，均 go1.26.3/1.26.4 已修）
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS server
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 用上一阶段构建好的前端覆盖 embed 目录（public/static/static.go: //go:embed all:dist）
RUN rm -rf public/static/dist
COPY --from=web /web/dist ./public/static/dist
ARG VERSION=docker
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build \
    -ldflags "-s -w -X 'github.com/dashug/ldap-admin-platform/public/version.Version=${VERSION}'" \
    -o /out/go-ldap-admin main.go

# ---------- 3) 运行时（极简） ----------
# 说明：进程经 entrypoint 的 su-exec 降权为非 root(app,uid 10001) 运行；不用顶层 USER 指令是为了让
# entrypoint 先以 root 修正 bind mount 的 /app/data 属主再降权。trivy DS-0002 静态检查在此为误报，
# 已在 .trivyignore 记录豁免。
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata su-exec && mkdir -p /app/data \
    && addgroup -S -g 10001 app && adduser -S -u 10001 -G app app
ENV TZ=Asia/Shanghai
COPY --from=server /out/go-ldap-admin ./go-ldap-admin
# 镜像内置一份【纯占位】的默认配置模板；正式部署务必用卷挂载覆盖或用环境变量注入真实配置/密钥
# （见 docker-compose.yml 与 README）。绝不把含真实密钥的 config.yml 烤进镜像。
COPY config.example.yml ./config.yml
COPY deploy/docker/entrypoint.sh /entrypoint.sh
RUN chown -R app:app /app && chmod +x /entrypoint.sh
# 通过 entrypoint 以 root 入口先把挂载卷 /app/data 补正为 app 属主，再用 su-exec 降权为非 root 的
# app(uid 10001) 运行主进程：既保留容器逃逸加固（进程实际非 root），又兼容 bind mount——避免非 root
# 无法写入首次由 Docker 以 root 创建的挂载目录而导致 sqlite/日志/RSA 私钥写入失败、启动崩溃。
EXPOSE 8888
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./go-ldap-admin"]
