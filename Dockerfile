# 第一阶段：构建阶段
FROM golang:1.23-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mcp-link .

# 第二阶段：运行阶段
FROM alpine:3.18

# 添加运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 创建非 root 用户
RUN adduser -D -g '' appuser

# 从构建阶段复制可执行文件
COPY --from=builder /app/mcp-link /usr/local/bin/

# 设置工作目录
WORKDIR /app

# 将所有权转移给非 root 用户
RUN chown -R appuser:appuser /app

# 使用非 root 用户运行应用
USER appuser

# 暴露端口
EXPOSE 8080

# 运行应用
ENTRYPOINT ["mcp-link"]
CMD ["serve", "--host", "0.0.0.0", "--port", "8080"] 