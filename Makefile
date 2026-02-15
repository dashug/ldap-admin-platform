# LDAP 管理平台 - 前后端一体构建
# 后端：Go，根目录；前端：Vue，web/ 目录

BINARY_NAME=go-ldap-admin
VERSION ?= dev
LDFLAGS := -X 'github.com/dashug/ldap-admin-platform/public/version.Version=$(VERSION)' \
	-X 'github.com/dashug/ldap-admin-platform/public/version.GitCommit=unknown' \
	-X 'github.com/dashug/ldap-admin-platform/public/version.BuildTime=$(shell date "+%Y-%m-%d %H:%M:%S")' \
	-X 'github.com/dashug/ldap-admin-platform/public/version.GoVersion=$(shell go version | awk "{print \$$3}")'

.PHONY: help
help:
	@echo "LDAP 管理平台 - 常用命令"
	@echo "  make deps-web    安装前端依赖 (web/node_modules)"
	@echo "  make build-web   构建前端，产出到 public/static/dist"
	@echo "  make build       仅构建 Go 后端（需已存在 public/static/dist）"
	@echo "  make all         先构建前端，再构建后端，得到单可执行文件"
	@echo "  make run         运行后端（开发时前端可单独 npm run dev）"
	@echo "  make clean       清理前端构建产物与后端二进制"

.PHONY: deps-web
deps-web:
	cd web && npm install

.PHONY: build-web
build-web:
	cd web && npm run build:prod
	@rm -rf public/static/dist
	@cp -r web/dist public/static/

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) main.go

.PHONY: all
all: build-web build
	@echo "构建完成: ./$(BINARY_NAME)"

.PHONY: run
run:
	go run main.go

.PHONY: clean
clean:
	rm -rf web/dist public/static/dist $(BINARY_NAME)
