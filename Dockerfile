# 多阶段自包含构建：一步得到「前端已 embed 进 Go 单二进制」的最小运行镜像。
#   docker build -t go-ldap-admin .
# 旧版（仅打包预编译二进制）见 deploy/Dockerfile.prebuilt。
# syntax=docker/dockerfile:1

# ---------- 1) 构建前端（Vue 3 + Vite，需 Node 18+） ----------
FROM node:18-alpine AS web
WORKDIR /web
# 先装依赖（利用层缓存）：有 lock 用 npm ci，否则回退 npm install
COPY web/package.json web/package-lock.json* ./
RUN npm ci || npm install
COPY web/ ./
RUN npm run build:prod

# ---------- 2) 编译后端（把前端 embed 进单二进制；纯 Go sqlite，关闭 CGO） ----------
FROM golang:1.25-alpine AS server
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 用上一阶段构建好的前端覆盖 embed 目录（public/static/static.go: //go:embed all:dist）
RUN rm -rf public/static/dist
COPY --from=web /web/dist ./public/static/dist
ARG VERSION=docker
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w -X 'github.com/dashug/ldap-admin-platform/public/version.Version=${VERSION}'" \
    -o /out/go-ldap-admin main.go

# ---------- 3) 运行时（极简） ----------
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
COPY --from=server /out/go-ldap-admin ./go-ldap-admin
# 镜像内置一份默认 config.yml；正式部署建议用卷挂载覆盖（见 docker-compose.yml）
COPY config.yml ./config.yml
EXPOSE 8888
CMD ["./go-ldap-admin"]
