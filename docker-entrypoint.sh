#!/bin/bash

# WeChatPadPro Docker 入口脚本

set -e

echo "=========================================="
echo "WeChatPadPro Docker 启动脚本"
echo "=========================================="

# 等待 MySQL 启动
echo "等待 MySQL 启动..."
MAX_TRIES=30
COUNTER=0
while [ $COUNTER -lt $MAX_TRIES ]; do
    if mysqladmin ping -h127.0.0.1 -uroot -p${MYSQL_ROOT_PASSWORD} &> /dev/null; then
        echo "MySQL 已启动"
        break
    fi
    COUNTER=$((COUNTER+1))
    echo "等待 MySQL... ($COUNTER/$MAX_TRIES)"
    sleep 2
done

if [ $COUNTER -eq $MAX_TRIES ]; then
    echo "错误: MySQL 启动超时"
    exit 1
fi

# 检查 Redis
echo "检查 Redis..."
if ! redis-cli ping &> /dev/null; then
    echo "错误: Redis 未启动"
    exit 1
fi
echo "Redis 已启动"

# 如果是启动 wechatpadpro 服务
if [ "$1" = "wechatpadpro" ]; then
    echo "启动 WeChatPadPro 服务..."

    # 加载环境变量
    if [ -f /.env ]; then
        set -a
        source /.env
        set +a
    fi

    # 设置默认环境变量
    export HOST=${HOST:-"0.0.0.0"}
    export PORT=${PORT:-"1238"}
    export DB_HOST=${DB_HOST:-"127.0.0.1"}
    export DB_PORT=${DB_PORT:-"3306"}
    export DB_DATABASE=${DB_DATABASE:-"weixin"}
    export DB_USERNAME=${DB_USERNAME:-"weixin"}
    export DB_PASSWORD=${DB_PASSWORD:-"${MYSQL_PASSWORD}"}
    export REDIS_HOST=${REDIS_HOST:-"127.0.0.1"}
    export REDIS_PORT=${REDIS_PORT:-"6379"}
    export REDIS_PASS=${REDIS_PASS:-""}
    export ADMIN_KEY=${ADMIN_KEY:-"12345"}
    export TZ=${TZ:-"Asia/Shanghai"}

    # 生成 MySQL 连接字符串
    export MYSQL_CONNECT_STR="${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}?charset=utf8mb4&parseTime=true&loc=Local"

    echo "环境变量:"
    echo "  HOST=$HOST"
    echo "  PORT=$PORT"
    echo "  DB_HOST=$DB_HOST"
    echo "  REDIS_HOST=$REDIS_HOST"
    echo "  TZ=$TZ"

    # 启动应用
    echo "启动 WeChatPadPro..."
    exec /app/wechatpadpro
fi

# 如果是初始化脚本，则退出
exec "$@"