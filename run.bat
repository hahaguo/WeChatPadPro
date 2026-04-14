@echo off
echo WeChatPadPro 启动脚本
echo 版本: ios18.61-861

REM 检查是否存在.env文件
if not exist .env (
    echo 警告：未找到.env文件，将使用.env.example创建
    copy .env.example .env
)

REM 运行程序
echo 正在启动WeChatPadPro...
start wechatpadpro-windows-amd64-vios18.61-861.exe

echo 如需查看日志，请检查程序目录下的日志文件 