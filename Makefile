.PHONY: help build clean run test docker-build docker-run

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
BINARY_NAME=wechatpadpro
VERSION=$(shell cat version.txt 2>/dev/null || echo "1.0.0")
BUILD_TIME=$(shell date +%Y%m%d_%H%M%S)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# 帮助信息
help:
	@echo "WeChatPadPro 构建工具"
	@echo ""
	@echo "可用命令:"
	@echo "  make build           - 构建当前平台二进制"
	@echo "  make build-all       - 构建所有平台二进制"
	@echo "  make run             - 运行开发服务器"
	@echo "  make test            - 运行测试"
	@echo "  make clean           - 清理构建文件"
	@echo "  make docker-build    - 构建Docker镜像"
	@echo "  make docker-run      - 运行Docker容器"
	@echo "  make docker-stop     - 停止Docker容器"

# 构建当前平台
build:
	@echo "构建 $(BINARY_NAME)..."
	@mkdir -p bin
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) cmd/main.go
	@echo "构建完成: bin/$(BINARY_NAME)"

# 构建所有平台
build-all:
	@echo "构建所有平台..."
	@mkdir -p releases

	@echo "构建 Linux AMD64..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o releases/$(BINARY_NAME)-$(VERSION)-linux-amd64 cmd/main.go

	@echo "构建 Linux ARM64..."
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o releases/$(BINARY_NAME)-$(VERSION)-linux-arm64 cmd/main.go

	@echo "构建 macOS AMD64..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o releases/$(BINARY_NAME)-$(VERSION)-macos-amd64 cmd/main.go

	@echo "构建 macOS ARM64..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o releases/$(BINARY_NAME)-$(VERSION)-macos-arm64 cmd/main.go

	@echo "构建 Windows AMD64..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o releases/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe cmd/main.go

	@echo "构建 Windows ARM64..."
	@GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o releases/$(BINARY_NAME)-$(VERSION)-windows-arm64.exe cmd/main.go

	@echo "所有平台构建完成!"

# 运行开发服务器
run:
	@go run cmd/main.go

# 运行测试
test:
	@go test -v ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf bin releases
	@echo "清理完成"

# 安装依赖
deps:
	@echo "安装依赖..."
	@go mod download
	@go mod tidy

# Docker 相关
docker-build:
	@echo "构建Docker镜像..."
	@docker build -t $(BINARY_NAME):$(VERSION) .
	@docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest
	@echo "Docker镜像构建完成"

docker-run:
	@echo "运行Docker容器..."
	@docker-compose up -d

docker-stop:
	@echo "停止Docker容器..."
	@docker-compose down

# 格式化代码
fmt:
	@echo "格式化代码..."
	@go fmt ./...
	@goimports -w .

# 代码检查
lint:
	@echo "代码检查..."
	@golangci-lint run

# 生成文档
docs:
	@echo "生成API文档..."
	@swag init -g cmd/main.go -o docs