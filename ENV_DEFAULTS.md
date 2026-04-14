# 环境变量默认值参考

## 概述

本文档列出了 All-in-One Docker 镜像中所有环境变量的默认值。

## 使用方法

可以通过以下方式覆盖环境变量：

```bash
# Docker 运行
docker run -e ADMIN_KEY=your_key -e MYSQL_PASSWORD=your_pass ...

# Docker Compose
environment:
  - ADMIN_KEY=your_key
  - MYSQL_PASSWORD=your_pass
```

---

## 基础系统配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `HOST` | `0.0.0.0` | 服务监听地址 |
| `PORT` | `1238` | 服务端口号 |
| `DEBUG` | `true` | 调试模式 |
| `ADMIN_KEY` | `7Kf9mP2xR5nQ8sV3wY6zA1bC4dE7gH0j` | 管理员密钥 |
| `VERSION` | `1.0.0` | 版本号 |
| `API_VERSION` | `` | API版本前缀 |
| `MCP_PORT` | `8099` | MCP服务端口 |
| `GH_WXID` | `` | 推广公众号微信ID |
| `TZ` | `Asia/Shanghai` | 时区 |

---

## MySQL 配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `MYSQL_ROOT_PASSWORD` | `9mK2pR5nQ8sV3wY6zA1bC4dE7gH0jM3p` | MySQL root密码 |
| `MYSQL_DATABASE` | `weixin` | 数据库名称 |
| `MYSQL_USER` | `weixin` | MySQL用户名 |
| `MYSQL_PASSWORD` | `5nQ8sV3wY6zA1bC4dE7gH0jM3pR6tS9v` | MySQL密码 |
| `DB_HOST` | `127.0.0.1` | 数据库主机 |
| `DB_PORT` | `3306` | 数据库端口 |
| `DB_USERNAME` | `weixin` | 数据库用户名 |
| `DB_PASSWORD` | `5nQ8sV3wY6zA1bC4dE7gH0jM3pR6tS9v` | 数据库密码 |
| `DB_DATABASE` | `weixin` | 数据库名 |
| `MYSQL_CONNECT_STR` | `weixin:5nQ8sV3wY6zA1bC4dE7gH0jM3pR6tS9v@tcp(127.0.0.1:3306)/weixin?charset=utf8mb4&parseTime=true&loc=Local` | MySQL连接字符串 |
| `MYSQL_MAX_OPEN_CONN` | `100` | 最大打开连接数 |
| `MYSQL_MAX_IDLE_CONN` | `10` | 最大空闲连接数 |
| `MYSQL_MAX_LIFETIME` | `1h` | 连接最大生命周期 |

---

## Redis 配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `REDIS_HOST` | `127.0.0.1` | Redis主机 |
| `REDIS_PORT` | `6379` | Redis端口 |
| `REDIS_PASS` | `8sV3wY6zA1bC4dE7gH0jM3pR6tS9vW2yK` | Redis密码 |
| `REDIS_DB` | `1` | Redis数据库编号 |
| `REDIS_MAX_IDLE` | `30` | 最大空闲连接 |
| `REDIS_MAX_ACTIVE` | `100` | 最大活动连接 |
| `REDIS_IDLE_TIMEOUT` | `5000` | 空闲超时(毫秒) |
| `REDIS_MAX_CONN_LIFETIME` | `3600` | 连接最大生命周期(秒) |
| `REDIS_CONNECT_TIMEOUT` | `5000` | 连接超时(毫秒) |
| `REDIS_READ_TIMEOUT` | `10000` | 读取超时(毫秒) |
| `REDIS_WRITE_TIMEOUT` | `10000` | 写入超时(毫秒) |
| `REDIS_POOL_SIZE` | `10` | 连接池大小 |
| `REDIS_MIN_IDLE_CONN` | `5` | 最小空闲连接 |
| `REDIS_MAX_RETRIES` | `3` | 最大重试次数 |
| `REDIS_DIAL_TIMEOUT` | `5s` | 拨号超时 |
| `REDIS_READ_TIMEOUT_SEC` | `3s` | 读取超时 |
| `REDIS_WRITE_TIMEOUT_SEC` | `3s` | 写入超时 |

---

## Webhook 配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `WEBHOOK_ENABLED` | `true` | 是否启用Webhook |
| `WEBHOOK_TIMEOUT` | `10` | 超时时间(秒) |
| `WEBHOOK_MAX_RETRIES` | `3` | 最大重试次数 |
| `WEBHOOK_QUEUE_SIZE` | `1000` | 队列大小 |
| `WEBHOOK_BATCH_SIZE` | `20` | 批处理大小 |
| `WEBHOOK_DIRECT_STREAM` | `true` | 直接流模式 |

---

## WebSocket 配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `WS_HANDSHAKE_TIMEOUT` | `10` | 握手超时(秒) |
| `WS_READ_BUFFER_SIZE` | `4096` | 读缓冲区大小 |
| `WS_WRITE_BUFFER_SIZE` | `4096` | 写缓冲区大小 |
| `WS_READ_DEADLINE` | `120` | 读超时(秒) |
| `WS_WRITE_DEADLINE` | `60` | 写超时(秒) |
| `WS_PING_INTERVAL` | `25` | Ping间隔(秒) |
| `WS_CONNECTION_CHECK_INTERVAL` | `45` | 连接检查间隔(秒) |
| `WS_MAX_MESSAGE_SIZE` | `8192` | 最大消息大小 |

---

## Worker 配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `WORKER_POOL_SIZE` | `500` | 工作池大小 |
| `MAX_WORKER_TASK_LEN` | `1000` | 最大任务队列长度 |
| `TASK_EXEC_WAIT_TIMES` | `500` | 任务执行等待时间(毫秒) |

---

## Auth 配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `TOKEN_EXPIRE` | `24h` | Token过期时间 |
| `REFRESH_EXPIRE` | `168h` | Refresh Token过期时间 |
| `AUTO_AUTH_INTERVAL` | `30m` | 自动认证间隔 |

---

## Web 配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `WEB_DOMAIN` | `localhost:1238` | Web域名 |
| `WEB_TASK_NAME` | `` | Web任务名称 |
| `WEB_TASK_APP_NUMBER` | `` | Web任务应用编号 |
| `NEWS_SYN_WXID` | `true` | 按微信ID同步消息 |
| `DT` | `true` | 启用DT |

---

## 消息队列配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `TOPIC` | `wx_sync_msg_topic` | 消息主题 |
| `ROCKET_MQ_ENABLED` | `false` | 启用RocketMQ |
| `ROCKET_MQ_HOST` | `127.0.0.1:9876` | RocketMQ主机 |
| `ROCKET_ACCESS_KEY` | `2pR6tS9vW2yK5nQ8sV3wY6zA1bC4dE7gH` | RocketMQ访问密钥 |
| `ROCKET_SECRET_KEY` | `4dE7gH0jM3pR6tS9vW2yK5nQ8sV3wY6zA1bC!@#` | RocketMQ密钥 |
| `RABBIT_MQ_ENABLED` | `false` | 启用RabbitMQ |
| `RABBIT_MQ_URL` | `amqp://user:6tS9vW2yK5nQ8sV3wY6zA1bC4dE7gH0jM3p@127.0.0.1:5672/` | RabbitMQ URL |
| `KAFKA_ENABLED` | `false` | 启用Kafka |
| `KAFKA_URL` | `` | Kafka URL |
| `KAFKA_USERNAME` | `` | Kafka用户名 |
| `KAFKA_PASSWORD` | `` | Kafka密码 |

---

## 任务配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `TASK_RETRY_COUNT` | `3` | 任务重试次数 |
| `TASK_RETRY_INTERVAL` | `5` | 任务重试间隔(秒) |
| `HEARTBEAT_INTERVAL` | `25` | 心跳间隔(秒) |
| `AUTO_SYNC_INTERVAL_MINUTES` | `30` | 自动同步间隔(分钟) |
| `QUEUE_EXPIRE_TIME` | `86400` | 队列过期时间(秒) |

---

## 集群配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `CLUSTER_NAME` | `` | 集群名称 |
| `ZK_ADDR` | `` | ZooKeeper地址 |
| `ETCD_ADDR` | `` | ETCD地址 |

---

## 其他配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `DISABLED_CMD_LIST` | `` | 禁用命令列表 |

---

## 安全说明

### 默认密码

| 变量 | 默认值 |
|------|--------|
| `ADMIN_KEY` | `7Kf9mP2xR5nQ8sV3wY6zA1bC4dE7gH0j` |
| `MYSQL_ROOT_PASSWORD` | `9mK2pR5nQ8sV3wY6zA1bC4dE7gH0jM3p` |
| `MYSQL_PASSWORD` | `5nQ8sV3wY6zA1bC4dE7gH0jM3pR6tS9v` |
| `REDIS_PASS` | `8sV3wY6zA1bC4dE7gH0jM3pR6tS9vW2yK` |

> ⚠️ 虽然默认密码较为复杂，但生产环境仍建议使用自定义密码。

### 生成安全密码

使用以下命令生成随机密码：

```bash
# 生成 32 位随机密码
openssl rand -base64 32

# 生成 40 位随机密码
openssl rand -base64 40

# 生成字母数字混合密码
cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1
```

### 修改密码示例

```bash
# Docker
docker run -d \
  -e ADMIN_KEY="$(openssl rand -base64 32)" \
  -e MYSQL_ROOT_PASSWORD="$(openssl rand -base64 32)" \
  -e MYSQL_PASSWORD="$(openssl rand -base64 32)" \
  -e REDIS_PASS="$(openssl rand -base64 32)" \
  wechatpadpro:all-in-one

# Docker Compose (在 docker-compose-all-in-one.yml 中或通过 .env 文件)
environment:
  - ADMIN_KEY=${ADMIN_KEY}
  - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
  - MYSQL_PASSWORD=${MYSQL_PASSWORD}
  - REDIS_PASS=${REDIS_PASS}
```

创建 `.env` 文件：

```ini
ADMIN_KEY=your_custom_admin_key_here
MYSQL_ROOT_PASSWORD=your_custom_root_password_here
MYSQL_PASSWORD=your_custom_mysql_password_here
REDIS_PASS=your_custom_redis_password_here
```