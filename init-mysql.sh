#!/bin/bash

# MySQL 初始化脚本

set -e

echo "=========================================="
echo "MySQL 数据库初始化"
echo "=========================================="

# 等待 MySQL 完全启动
echo "等待 MySQL 完全启动..."
sleep 10

# 检查是否已经初始化
if [ -f /var/lib/mysql/.initialized ]; then
    echo "数据库已初始化，跳过"
    exit 0
fi

echo "开始初始化数据库..."

# 设置 root 密码
mysql -uroot -e "ALTER USER 'root'@'localhost' IDENTIFIED BY '${MYSQL_ROOT_PASSWORD}';" 2>/dev/null || true
mysql -uroot -p${MYSQL_ROOT_PASSWORD} -e "ALTER USER 'root'@'localhost' IDENTIFIED BY '${MYSQL_ROOT_PASSWORD}';" || true

# 创建数据库和用户
mysql -uroot -p${MYSQL_ROOT_PASSWORD} <<EOF
CREATE DATABASE IF NOT EXISTS \`${MYSQL_DATABASE}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS '${MYSQL_USER}'@'%' IDENTIFIED BY '${MYSQL_PASSWORD}';
GRANT ALL PRIVILEGES ON \`${MYSQL_DATABASE}\`.* TO '${MYSQL_USER}'@'%';
FLUSH PRIVILEGES;
EOF

# 导入初始数据
if [ -f /app/mysql_init.sql ]; then
    echo "导入初始数据..."
    mysql -uroot -p${MYSQL_ROOT_PASSWORD} ${MYSQL_DATABASE} < /app/mysql_init.sql
fi

# 标记为已初始化
touch /var/lib/mysql/.initialized

echo "数据库初始化完成!"

# 退出，让 supervisord 继续启动其他服务
exit 0