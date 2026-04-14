#!/bin/bash

# WeChatPadPro 密码生成脚本

set -e

echo "=========================================="
echo "WeChatPadPro 密码生成工具"
echo "=========================================="
echo ""

# 函数：生成随机密码
generate_password() {
    local length=${1:-32}
    openssl rand -base64 $length | tr -d '=+/' | cut -c1-$length
}

# 生成密码
ADMIN_KEY=$(generate_password 32)
MYSQL_ROOT_PASSWORD=$(generate_password 32)
MYSQL_PASSWORD=$(generate_password 32)
REDIS_PASS=$(generate_password 32)

# 输出结果
echo "生成的安全密码："
echo "=========================================="
echo ""
echo "ADMIN_KEY=$ADMIN_KEY"
echo "MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD"
echo "MYSQL_PASSWORD=$MYSQL_PASSWORD"
echo "REDIS_PASS=$REDIS_PASS"
echo ""
echo "=========================================="
echo ""

# 询问是否保存到 .env 文件
read -p "是否保存到 .env 文件? (y/n): " save_env

if [ "$save_env" = "y" ] || [ "$save_env" = "Y" ]; then
    cat > .env << EOF
# WeChatPadPro 环境变量配置
# 自动生成于 $(date +%Y-%m-%d\ %H:%M:%S)

# ========== 基础系统配置 ==========
HOST=0.0.0.0
PORT=1238
DEBUG=false
ADMIN_KEY=$ADMIN_KEY
VERSION=1.0.0
TZ=Asia/Shanghai

# ========== MySQL 配置 ==========
MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD
MYSQL_DATABASE=weixin
MYSQL_USER=weixin
MYSQL_PASSWORD=$MYSQL_PASSWORD
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USERNAME=weixin
DB_DATABASE=weixin
DB_PASSWORD=$MYSQL_PASSWORD

# ========== Redis 配置 ==========
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASS=$REDIS_PASS
REDIS_DB=1

# ========== 其他配置 ==========
WEBHOOK_ENABLED=true
WORKER_POOL_SIZE=500
AUTO_AUTH_INTERVAL=30m
WEB_DOMAIN=localhost:1238
NEWS_SYN_WXID=true
DT=true
EOF

    echo "已保存到 .env 文件"
    echo ""
    echo "请妥善保管 .env 文件，不要提交到版本控制系统！"
else
    echo "密码已生成，请手动保存。"
fi

echo ""
echo "提示：使用 docker-compose 启动时，将自动读取 .env 文件中的配置"
echo ""
echo "启动命令："
echo "docker-compose -f docker-compose-all-in-one.yml up -d"
echo ""