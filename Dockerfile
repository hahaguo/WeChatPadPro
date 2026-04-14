FROM alpine:3.20

LABEL maintainer="WeChatPadPro"
LABEL description="基于 WeChat Pad 协议的高級管理工具"

# 设置时区
ENV TZ=Asia/Shanghai

# 安装必要的 CA 证书
RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone && \
    apk del tzdata

# 创建工作目录
WORKDIR /app

# 复制可执行文件 (根据实际架构调整)
COPY wechatpadpro-linux-amd64-vios18.61-861 ./wechatpadpro

# 复制配置文件
COPY .env ./
COPY webhook_config.json ./

# 复制 assets 目录下的配置文件
COPY assets/owner.json ./assets/
COPY assets/webhook_config.json ./assets/
COPY assets/white_group.json ./assets/
COPY assets/meta-inf.yaml ./assets/
COPY assets/sae.dat ./assets/
COPY assets/ca-cert ./assets/

# 复制 Redis 配置文件 (如果需要)
COPY redis/redis.conf ./redis/redis.conf

# 复制 MySQL 初始化脚本
# COPY win数据库/01_InitMySQL ./01_InitMySQL

# 复制数据库 SQL 文件
COPY wechat_mmtls.sql ./wechat_mmtls.sql

# 复制模板文件 (如果 Web 界面需要)
COPY static/ ./static/

# 复制 Python webhook 客户端 (可选)
COPY webhook-client.py ./

# 设置权限
RUN chmod +x ./wechatpadpro

# 创建日志目录
RUN mkdir -p /app/logs

# 暴露端口
# 1238 - 管理端口
# 8080 - API 端口
# 8099 - MCP 服务端口 (可选)
EXPOSE 1238 8080 8099

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:1238/login/GetLoginStatus || exit 1

# 运行应用
CMD ["./wechatpadpro"]
