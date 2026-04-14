# Docker 部署指南

## 快速开始

### 1. 准备环境

确保已安装以下软件：
- Docker (20.10+)
- Docker Compose (2.0+)

### 2. 配置环境变量

复制示例配置文件并修改：

```bash
cp .env.example .env
```

编辑 `.env` 文件，至少修改以下重要配置：

```ini
# 管理员密钥（必须修改）
ADMIN_KEY=your_secure_admin_key_here

# MySQL 密码（建议修改）
MYSQL_ROOT_PASSWORD=your_root_password
MYSQL_PASSWORD=your_mysql_password

# MySQL 连接字符串（注意修改密码）
MYSQL_CONNECT_STR=weixin:your_mysql_password@tcp(mysql:3306)/weixin?charset=utf8mb4&parseTime=true&loc=Local
```

### 3. 构建并启动

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 4. 访问服务

- **管理界面**: http://localhost:1238
- **API 接口**: http://localhost:8080
- **MCP 服务**: http://localhost:8099 (如启用)

### 5. 获取授权码

```bash
curl "http://localhost:1238/login/GenAuthKey2?key=YOUR_ADMIN_KEY&count=1&days=365"
```

## 常用命令

```bash
# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f wechatpadpro
docker-compose logs -f mysql
docker-compose logs -f redis

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 停止并删除数据卷（慎用！）
docker-compose down -v

# 进入容器
docker-compose exec wechatpadpro sh

# 查看实时日志
docker-compose logs --tail=100 -f wechatpadpro
```

## 目录结构

```
WeChatPadPro/
├── Dockerfile              # 主应用镜像构建文件
├── docker-compose.yml      # 多服务编排文件
├── .dockerignore          # Docker 构建忽略文件
├── .env                   # 环境变量配置
├── .env.example           # 环境变量配置示例
├── wechatpadpro-linux-amd64-vios18.61-861  # Linux 可执行文件
├── assets/                # 配置文件目录
├── logs/                  # 日志目录（挂载到宿主机）
├── redis/                 # Redis 配置
├── mysql/                 # MySQL 配置
└── wechat_mmtls.sql       # 数据库初始化脚本
```

## 端口说明

| 端口 | 用途 | 内部端口 |
|------|------|----------|
| 1238 | 管理界面 | 1238 |
| 8080 | API 接口 | 8080 |
| 3306 | MySQL 数据库 | 3306 |
| 6379 | Redis 缓存 | 6379 |
| 8099 | MCP 服务（可选）| 8099 |

## 数据持久化

以下目录已配置为数据卷：

- `mysql_data`: MySQL 数据目录
- `redis_data`: Redis 数据目录
- `./logs`: 应用日志目录

## 故障排查

### 服务无法启动

1. 检查端口是否被占用：
```bash
netstat -tuln | grep -E '1238|8080|3306|6379'
```

2. 查看日志定位问题：
```bash
docker-compose logs wechatpadpro
```

### 数据库连接失败

1. 检查 MySQL 容器状态：
```bash
docker-compose ps mysql
```

2. 进入 MySQL 容器测试：
```bash
docker-compose exec mysql mysql -u weixin -p
```

### Redis 连接失败

1. 检查 Redis 容器状态：
```bash
docker-compose ps redis
```

2. 测试 Redis 连接：
```bash
docker-compose exec redis redis-cli ping
```

## 安全建议

1. **修改默认密码** - 部署前务必修改所有默认密码
2. **限制网络访问** - 使用防火墙限制端口访问
3. **定期更新** - 及时更新 Docker 镜像
4. **数据备份** - 定期备份 MySQL 和 Redis 数据
5. **日志监控** - 定期检查日志文件

## 备份与恢复

### 备份数据

```bash
# 备份 MySQL
docker-compose exec mysql mysqldump -u weixin -p weixin > backup_$(date +%Y%m%d).sql

# 备份 Redis
docker-compose exec redis redis-cli --rdb /data/dump_$(date +%Y%m%d).rdb
```

### 恢复数据

```bash
# 恢复 MySQL
cat backup_20240101.sql | docker-compose exec -T mysql mysql -u weixin -p weixin

# 恢复 Redis
docker-compose exec redis redis-cli --rdb /data/dump_20240101.rdb
```
