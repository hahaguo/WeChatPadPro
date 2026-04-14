# All-in-One Docker 部署指南

## 概述

`Dockerfile.all-in-one` 将 WeChatPadPro、MySQL 和 Redis 打包到一个镜像中，简化部署流程。

## 镜像特性

- ✅ 包含 WeChatPadPro 主应用
- ✅ 包含 MySQL 8.0 数据库
- ✅ 包含 Redis 缓存
- ✅ 使用 supervisord 管理进程
- ✅ 自动初始化数据库
- ✅ 健康检查支持
- ✅ 安全默认密码配置

## 默认安全密码

| 变量 | 默认值 |
|------|--------|
| `ADMIN_KEY` | `7Kf9mP2xR5nQ8sV3wY6zA1bC4dE7gH0j` |
| `MYSQL_ROOT_PASSWORD` | `9mK2pR5nQ8sV3wY6zA1bC4dE7gH0jM3p` |
| `MYSQL_PASSWORD` | `5nQ8sV3wY6zA1bC4dE7gH0jM3pR6tS9v` |
| `REDIS_PASS` | `8sV3wY6zA1bC4dE7gH0jM3pR6tS9vW2yK` |

> ⚠️ 生产环境仍建议使用自定义密码！

## 快速开始

### 方式一：使用默认密码（快速测试）

```bash
# 1. 构建镜像
docker build -f Dockerfile.all-in-one -t wechatpadpro:all-in-one .

# 2. 运行容器
docker run -d \
  --name wechatpadpro \
  -p 3306:3306 \
  -p 6379:6379 \
  -p 8080:8080 \
  -p 1238:1238 \
  -p 8099:8099 \
  -v wechatpadpro-mysql:/var/lib/mysql \
  -v wechatpadpro-redis:/var/lib/redis \
  -v wechatpadpro-logs:/app/logs \
  wechatpadpro:all-in-one

# 3. 获取授权码（使用默认 ADMIN_KEY）
curl "http://localhost:1238/api/login/GenAuthKey2?key=7Kf9mP2xR5nQ8sV3wY6zA1bC4dE7gH0j&count=1&days=365"
```

### 方式二：使用自定义密码（推荐）

```bash
# 1. 生成安全密码
./generate-passwords.sh

# 或手动生成
ADMIN_KEY=$(openssl rand -base64 32)
MYSQL_PASSWORD=$(openssl rand -base64 32)

# 2. 运行容器
docker run -d \
  --name wechatpadpro \
  -p 3306:3306 \
  -p 6379:6379 \
  -p 8080:8080 \
  -p 1238:1238 \
  -p 8099:8099 \
  -e ADMIN_KEY="$ADMIN_KEY" \
  -e MYSQL_ROOT_PASSWORD="$MYSQL_PASSWORD" \
  -e MYSQL_PASSWORD="$MYSQL_PASSWORD" \
  -e REDIS_PASS="$REDIS_PASS" \
  -v wechatpadpro-mysql:/var/lib/mysql \
  -v wechatpadpro-redis:/var/lib/redis \
  -v wechatpadpro-logs:/app/logs \
  wechatpadpro:all-in-one

# 3. 使用自定义 ADMIN_KEY 获取授权码
curl "http://localhost:1238/api/login/GenAuthKey2?key=$ADMIN_KEY&count=1&days=365"
```

### 方式三：使用 Docker Compose

```bash
# 1. 复制模板并修改密码
cp .env.template .env

# 编辑 .env 文件，修改密码配置
vim .env

# 或使用密码生成脚本
./generate-passwords.sh

# 2. 启动服务
docker-compose -f docker-compose-all-in-one.yml up -d

# 3. 查看日志
docker-compose -f docker-compose-all-in-one.yml logs -f

# 4. 停止服务
docker-compose -f docker-compose-all-in-one.yml down
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `ADMIN_KEY` | `7Kf9mP2xR5nQ8sV3wY6zA1bC4dE7gH0j` | 管理员密钥 |
| `MYSQL_ROOT_PASSWORD` | `9mK2pR5nQ8sV3wY6zA1bC4dE7gH0jM3p` | MySQL root 密码 |
| `MYSQL_DATABASE` | `weixin` | 数据库名称 |
| `MYSQL_USER` | `weixin` | MySQL 用户名 |
| `MYSQL_PASSWORD` | `5nQ8sV3wY6zA1bC4dE7gH0jM3pR6tS9v` | MySQL 密码 |
| `REDIS_HOST` | `127.0.0.1` | Redis 主机 |
| `REDIS_PORT` | `6379` | Redis 端口 |
| `REDIS_PASS` | `8sV3wY6zA1bC4dE7gH0jM3pR6tS9vW2yK` | Redis 密码 |
| `TZ` | `Asia/Shanghai` | 时区 |

> 完整环境变量列表请参考 [ENV_DEFAULTS.md](ENV_DEFAULTS.md)

## 端口映射

| 端口 | 服务 | 说明 |
|------|------|------|
| 3306 | MySQL | 数据库 |
| 6379 | Redis | 缓存 |
| 8080 | API | WeChatPadPro API |
| 1238 | Web | 管理界面 |
| 8099 | MCP | MCP 服务（可选） |

## 数据持久化

```bash
# MySQL 数据
-v wechatpadpro-mysql:/var/lib/mysql

# Redis 数据
-v wechatpadpro-redis:/var/lib/redis

# 应用日志
-v wechatpadpro-logs:/app/logs

# 自定义配置（可选）
-v $(pwd)/.env:/app/.env:ro
```

## 密码生成工具

### 使用内置脚本

```bash
./generate-passwords.sh
```

脚本会：
1. 生成所有需要的安全密码
2. 显示生成的密码
3. 询问是否保存到 `.env` 文件

### 手动生成密码

```bash
# 生成 32 位随机密码
openssl rand -base64 32

# 生成 40 位随机密码
openssl rand -base64 40

# 生成字母数字混合密码
cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1
```

## 健康检查

```bash
# 检查容器健康状态
docker ps wechatpadpro

# 查看健康检查详情
docker inspect wechatpadpro | grep -A 10 Health

# 手动检查健康接口
curl http://localhost:1238/health
```

## 查看日志

```bash
# 查看所有日志
docker logs -f wechatpadpro

# 查看特定服务日志
docker exec wechatpadpro tail -f /var/log/mysql/mysqld.log
docker exec wechatpadpro tail -f /var/log/redis/redis.log
docker exec wechatpadpro tail -f /app/logs/wechatpadpro.log

# 查看 supervisord 日志
docker exec wechatpadpro tail -f /var/log/supervisord.log
```

## 故障排查

### 容器无法启动

```bash
# 查看详细日志
docker logs wechatpadpro

# 进入容器检查
docker exec -it wechatpadpro sh

# 检查服务状态
supervisorctl status
```

### 数据库连接失败

```bash
# 进入容器
docker exec -it wechatpadpro sh

# 检查 MySQL
mysql -uroot -p${MYSQL_ROOT_PASSWORD} -e "SHOW DATABASES;"

# 检查连接
mysql -uweixin -p${MYSQL_PASSWORD} -e "SELECT 1;"
```

### 服务未就绪

镜像使用 supervisord 管理服务，可能需要等待 30-60 秒让所有服务完全启动。

```bash
# 等待服务就绪
sleep 60

# 测试连接
curl http://localhost:1238/health
```

### 重置数据

```bash
# 停止容器
docker stop wechatpadpro
docker rm wechatpadpro

# 删除数据卷
docker volume rm wechatpadpro-mysql wechatpadpro-redis

# 重新启动
docker run -d \
  --name wechatpadpro \
  -v wechatpadpro-mysql:/var/lib/mysql \
  ... (其他参数)
  wechatpadpro:all-in-one
```

## 备份与恢复

### 备份数据

```bash
# 创建备份目录
mkdir -p backup

# 备份所有数据
docker run --rm \
  --volumes-from wechatpadpro \
  -v $(pwd)/backup:/backup \
  alpine tar czf /backup/wechatpadpro-backup-$(date +%Y%m%d_%H%M%S).tar.gz \
  /var/lib/mysql \
  /var/lib/redis

# 备份配置
cp .env backup/.env.backup
```

### 恢复数据

```bash
# 停止容器
docker stop wechatpadpro

# 恢复数据
docker run --rm \
  --volumes-from wechatpadpro \
  -v $(pwd)/backup:/backup \
  alpine tar xzf /backup/wechatpadpro-backup-YYYYMMDD_HHMMSS.tar.gz -C /

# 恢复配置
cp backup/.env.backup .env

# 启动容器
docker start wechatpadpro
```

## 更新密码

### 方法一：重新创建容器

```bash
# 1. 停止并删除旧容器
docker stop wechatpadpro
docker rm wechatpadpro

# 2. 使用新密码重新运行
docker run -d \
  --name wechatpadpro \
  -e ADMIN_KEY="新密码" \
  -e MYSQL_PASSWORD="新密码" \
  ... (其他参数)
  wechatpadpro:all-in-one
```

### 方法二：更新 .env 文件

```bash
# 1. 生成新密码
./generate-passwords.sh

# 2. 重新启动容器（会读取新的 .env 文件）
docker-compose -f docker-compose-all-in-one.yml down
docker-compose -f docker-compose-all-in-one.yml up -d
```

## 注意事项

1. **密码安全**:
   - 默认密码已设置为复杂随机密码
   - 生产环境仍建议使用自定义密码
   - 使用 `generate-passwords.sh` 生成安全密码

2. **资源要求**:
   - 建议至少分配 2GB 内存给容器
   - 确保有足够的磁盘空间用于数据持久化

3. **数据持久化**:
   - 务必挂载数据卷以实现持久化
   - 定期备份数据

4. **网络安全**:
   - 生产环境不要暴露 MySQL (3306) 端口
   - 使用防火墙限制访问
   - 考虑使用 Docker 网络隔离

5. **日志管理**:
   - 定期清理日志文件
   - 或使用日志轮转配置

## 相关文档

- [ENV_DEFAULTS.md](ENV_DEFAULTS.md) - 环境变量完整列表
- [DEPLOY.md](DEPLOY.md) - 分离式部署指南
- [DEVELOPMENT.md](DEVELOPMENT.md) - 开发指南