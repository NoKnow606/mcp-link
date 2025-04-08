# MCP Link Docker 部署指南

本文档提供了如何使用 Docker 和 Docker Compose 部署 MCP Link 应用的详细说明。

## 先决条件

- 安装 [Docker](https://docs.docker.com/get-docker/)
- 安装 [Docker Compose](https://docs.docker.com/compose/install/)

## 快速开始

1. 克隆仓库：

```bash
git clone https://github.com/anyisalin/mcp-openapi-to-mcp-adapter.git
cd mcp-openapi-to-mcp-adapter
```

2. （可选）修改 `.env` 文件配置：

```bash
# 根据实际情况修改环境变量
nano .env
```

3. 构建并启动服务：

```bash
docker-compose up -d
```

4. 验证服务运行状态：

```bash
# 检查服务状态
docker-compose ps

# 检查日志
docker-compose logs -f
```

5. 访问应用：

服务启动后，可通过 http://localhost:8080 访问 MCP Link 应用。

## 环境变量配置

可以通过修改 `.env` 文件或在 `docker-compose up` 命令中直接传递环境变量来配置应用：

| 变量名 | 描述 | 默认值 |
|-------|-----|-------|
| BASE_URL | 应用基础 URL，用于生成 SSE 链接 | http://localhost:8080 |
| MONGODB_URI | MongoDB 连接字符串 | mongodb://mongo:27017 |
| MONGODB_DATABASE | MongoDB 数据库名称 | mcp_link |
| MONGODB_USERNAME | MongoDB 用户名（可选） | - |
| MONGODB_PASSWORD | MongoDB 密码（可选） | - |
| MONGODB_AUTH_DATABASE | MongoDB 认证数据库（可选） | admin |
| MONGODB_HEARTBEAT_INTERVAL | MongoDB 心跳间隔（秒） | 10 |

## 生产环境部署

对于生产环境部署，建议采取以下措施：

1. 启用 MongoDB 认证：
   - 设置 MONGODB_USERNAME 和 MONGODB_PASSWORD
   - 修改 MONGODB_URI 包含认证信息

2. 使用 HTTPS：
   - 配置反向代理（如 Nginx）处理 SSL/TLS
   - 更新 BASE_URL 为 https:// 开头

3. 数据持久化：
   - 默认配置已包含 MongoDB 数据持久化
   - 重要数据建议定期备份

## 示例：使用 Nginx 反向代理

1. 创建 `nginx.conf` 文件：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://mcp-link:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # SSE 设置
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_buffering off;
        proxy_cache off;
        proxy_read_timeout 3600s;
    }
}
```

2. 将 Nginx 添加到 docker-compose.yml：

```yaml
nginx:
  image: nginx:alpine
  ports:
    - "80:80"
    - "443:443"
  volumes:
    - ./nginx.conf:/etc/nginx/conf.d/default.conf
    - ./certs:/path/to/certs
  depends_on:
    - mcp-link
  networks:
    - mcp-network
```

## 故障排除

1. 服务无法启动：
   - 检查日志：`docker-compose logs mcp-link`
   - 验证 MongoDB 连接：`docker-compose logs mongo`

2. 无法连接 MongoDB：
   - 确认 MONGODB_URI 正确
   - 检查 MongoDB 容器是否健康：`docker-compose ps`

3. 健康检查失败：
   - 访问 http://localhost:8080/health 查看详细错误信息
   - 检查应用日志：`docker-compose logs mcp-link` 