# 使用官方Go镜像作为构建阶段
FROM golang:1.25-alpine AS builder
WORKDIR /app
# 复制依赖文件（利用Docker缓存层，加速后续构建）
COPY go.mod go.sum ./
RUN go mod download
# 复制源代码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./main.go

# 使用极简的Alpine镜像作为运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# 从构建阶段拷贝编译好的二进制文件
COPY --from=builder /app/server .
# 拷贝配置文件（重要！）
COPY ./config ./config
# 暴露你的应用端口（与config.yml中一致）
EXPOSE 3000
# 启动应用
CMD ["./server"]