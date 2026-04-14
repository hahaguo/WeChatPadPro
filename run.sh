#!/bin/bash

echo "WeChatPadPro 启动脚本"
echo "版本: ios18.61-861"

# 检测操作系统类型
OS=$(uname -s)
ARCH=$(uname -m)

# 转换架构名称为脚本使用的格式
case $ARCH in
    x86_64)
        ARCH_NAME="amd64"
        ;;
    aarch64|arm64)
        ARCH_NAME="arm64"
        ;;
    mips64)
        ARCH_NAME="mips64"
        ;;
    ppc64)
        ARCH_NAME="ppc64"
        ;;
    riscv64)
        ARCH_NAME="riscv64"
        ;;
    *)
        echo "警告：无法识别的架构 $ARCH，默认使用amd64"
        ARCH_NAME="amd64"
        ;;
esac

# 根据操作系统类型选择可执行文件
case $OS in
    Darwin)
        EXECUTABLE="wechatpadpro-macos-${ARCH_NAME}-vios18.61-861"
        ;;
    Linux)
        EXECUTABLE="wechatpadpro-linux-${ARCH_NAME}-vios18.61-861"
        ;;
    FreeBSD)
        EXECUTABLE="wechatpadpro-freebsd-${ARCH_NAME}-vios18.61-861"
        ;;
    OpenBSD)
        EXECUTABLE="wechatpadpro-openbsd-${ARCH_NAME}-vios18.61-861"
        ;;
    *)
        echo "错误：不支持的操作系统 $OS"
        exit 1
        ;;
esac

# 检查是否存在.env文件
if [ ! -f .env ]; then
    echo "警告：未找到.env文件，将使用.env.example创建"
    cp .env.example .env
fi

# 检查可执行文件权限
if [ ! -x "$EXECUTABLE" ]; then
    echo "设置可执行权限..."
    chmod +x "$EXECUTABLE"
fi

# 运行程序
echo "正在启动WeChatPadPro..."
./$EXECUTABLE

echo "如需查看日志，请检查程序目录下的日志文件" 