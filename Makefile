.PHONY: help build run test clean fmt vet lint docker-build docker-run

help:
	@echo "可用命令:"
	@echo "  make build      - 构建项目"
	@echo "  make run        - 运行项目"
	@echo "  make test       - 运行测试"
	@echo "  make clean      - 清理构建文件"
	@echo "  make fmt        - 格式化代码"
	@echo "  make vet        - 代码检查"
	@echo "  make lint       - 代码规范检查"
	@echo "  make docker-build - 构建 Docker 镜像"
	@echo "  make docker-run   - 运行 Docker 容器"

build:
	@echo "构建项目..."
	@go build -o bin/exchangeapp main.go

run:
	@echo "运行项目..."
	@go run main.go

test:
	@echo "运行测试..."
	@go test -v ./...

clean:
	@echo "清理构建文件..."
	@rm -rf bin/
	@rm -rf logs/

fmt:
	@echo "格式化代码..."
	@go fmt ./...

vet:
	@echo "代码检查..."
	@go vet ./...

lint:
	@echo "代码规范检查..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，跳过"; \
	fi

docker-build:
	@echo "构建 Docker 镜像..."
	@docker build -t exchangeapp:latest .

docker-run:
	@echo "运行 Docker 容器..."
	@docker-compose up -d

docker-stop:
	@echo "停止 Docker 容器..."
	@docker-compose down

mod-tidy:
	@echo "整理依赖..."
	@go mod tidy

mod-update:
	@echo "更新依赖..."
	@go get -u ./...
	@go mod tidy
