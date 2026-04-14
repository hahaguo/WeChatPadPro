 # Webhook故障排查完整指南

## 📋 目录

- [常见问题](#常见问题)
- [诊断工具](#诊断工具)
- [日志分析](#日志分析)
- [配置检查](#配置检查)
- [网络诊断](#网络诊断)
- [性能优化](#性能优化)
- [紧急修复](#紧急修复)

## 🚨 常见问题

### 1. Webhook无法接收消息

#### 症状
- Webhook端点没有收到任何消息
- 日志显示"配置不存在或未启用"
- Swagger测试接口返回错误

#### 检查步骤
```bash
# 1. 检查配置是否启用
curl "http://localhost:1238/webhook/Status?key=your_key"

# 2. 测试Webhook连接
curl "http://localhost:1238/webhook/Test?key=your_key"

# 3. 检查配置列表
curl "http://localhost:1238/webhook/List"
```

#### 解决方案
```json
{
    "enabled": true,
    "useDirectStream": true,
    "useRedisSync": false,
    "independentMode": true
}
```

### 2. 消息延迟或丢失

#### 症状
- 消息推送延迟超过5秒
- 部分消息没有收到
- 日志显示"跳过连接前的消息"

#### 检查步骤
```bash
# 1. 检查独立模式配置
curl "http://localhost:1238/webhook/Diagnostics?key=your_key"

# 2. 检查连接时间
curl "http://localhost:1238/webhook/Status?key=your_key"

# 3. 重置连接
curl -X POST "http://localhost:1238/webhook/ResetConnection?key=your_key"
```

#### 解决方案
```json
{
    "useDirectStream": true,
    "useRedisSync": false,
    "independentMode": true,
    "timeout": 5
}
```

### 3. 频繁重试失败

#### 症状
- 日志显示大量重试错误
- 消息发送成功率低
- 服务端返回错误状态码

#### 检查步骤
```bash
# 1. 检查服务端状态
curl -I "http://your-webhook-url.com/webhook"

# 2. 测试网络连接
ping your-webhook-url.com

# 3. 检查超时设置
curl "http://localhost:1238/webhook/Status?key=your_key"
```

#### 解决方案
```json
{
    "timeout": 10,
    "retryCount": 2,
    "url": "https://your-webhook-url.com/webhook"
}
```

### 4. 签名验证失败

#### 症状
- 服务端报告签名错误
- 日志显示"签名验证失败"
- Webhook消息被拒绝

#### 检查步骤
```bash
# 1. 验证密钥配置
curl "http://localhost:1238/webhook/Status?key=your_key"

# 2. 检查密钥格式
# 确保密钥长度≥16字符，包含字母、数字、特殊字符
```

#### 解决方案
```json
{
    "secret": "wh_sk_2024_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

### 5. 消息类型过滤问题

#### 症状
- 某些类型的消息没有收到
- 日志显示"消息类型不匹配"
- 配置了特定消息类型但收不到

#### 检查步骤
```bash
# 1. 检查消息类型配置
curl "http://localhost:1238/webhook/Status?key=your_key"

# 2. 查看消息类型列表
# 1=文本, 3=图片, 47=表情, 49=链接, *=所有
```

#### 解决方案
```json
{
    "messageTypes": ["*"]  // 接收所有消息类型
}
```

### 6. 自己消息过滤问题

#### 症状
- 收不到自己发送的消息
- 收到不应该收到的自己消息
- `includeSelfMessage`配置不生效

#### 检查步骤
```bash
# 1. 检查配置
curl "http://localhost:1238/webhook/Status?key=your_key"

# 2. 查看诊断信息
curl "http://localhost:1238/webhook/Diagnostics?key=your_key"
```

#### 解决方案
```json
{
    "includeSelfMessage": true,  // 包含自己消息
    "wxId": ""                   // 留空表示所有账号
}
```

## 🔧 诊断工具

### API诊断接口

#### 1. 获取状态信息
```bash
curl "http://localhost:1238/webhook/Status?key=your_key"
```

**返回示例**:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "enabled": true,
        "url": "http://example.com/webhook",
        "lastSendTime": 1640995200,
        "lastSendStatus": true,
        "totalSent": 100,
        "totalFailed": 2
    }
}
```

#### 2. 获取诊断信息
```bash
curl "http://localhost:1238/webhook/Diagnostics?key=your_key"
```

**返回示例**:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "config": {
            "enabled": true,
            "useDirectStream": true,
            "useRedisSync": false,
            "independentMode": true
        },
        "connection": {
            "lastConnectTime": 1640995200,
            "isConnected": true,
            "responseTime": 150
        },
        "statistics": {
            "totalSent": 100,
            "totalFailed": 2,
            "successRate": 98.0
        }
    }
}
```

#### 3. 测试连接
```bash
curl "http://localhost:1238/webhook/Test?key=your_key"
```

#### 4. 重置连接
```bash
curl -X POST "http://localhost:1238/webhook/ResetConnection?key=your_key"
```

#### 5. 重新初始化
```bash
curl -X POST "http://localhost:1238/webhook/Reinitialize?key=your_key"
```

### 命令行诊断工具

#### 1. 检查服务状态
```bash
# 检查WeChatPadPro进程
ps aux | grep wechatpadpro

# 检查端口监听
netstat -tlnp | grep 1238

# 检查日志文件
tail -f logs/webhook.log
```

#### 2. 网络连接测试
```bash
# 测试Webhook URL连接
curl -I "http://your-webhook-url.com/webhook"

# 测试DNS解析
nslookup your-webhook-url.com

# 测试端口连通性
telnet your-webhook-url.com 80
```

#### 3. 性能监控
```bash
# 监控CPU和内存使用
top -p $(pgrep wechatpadpro)

# 监控网络连接
netstat -an | grep :1238

# 监控磁盘使用
df -h
```

## 📊 日志分析

### 日志文件位置
```
logs/webhook.log          # Webhook专用日志
logs/wechatpadpro.log     # 主程序日志
logs/error.log           # 错误日志
```

### 关键日志模式

#### 1. 消息处理日志
```bash
# 查看消息处理日志
grep "直接推送" logs/webhook.log | tail -20

# 查看消息过滤日志
grep "跳过" logs/webhook.log | tail -20

# 查看重试日志
grep "重试" logs/webhook.log | tail -20
```

#### 2. 连接状态日志
```bash
# 查看连接建立日志
grep "连接" logs/webhook.log | tail -20

# 查看连接错误日志
grep "ERROR" logs/webhook.log | tail -20

# 查看配置加载日志
grep "配置" logs/webhook.log | tail -20
```

#### 3. 性能监控日志
```bash
# 查看响应时间日志
grep "响应时间" logs/webhook.log | tail -20

# 查看队列状态日志
grep "队列" logs/webhook.log | tail -20

# 查看统计信息日志
grep "统计" logs/webhook.log | tail -20
```

### 日志分析脚本

#### 1. 错误统计脚本
```bash
#!/bin/bash
echo "=== Webhook错误统计 ==="
grep "ERROR" logs/webhook.log | wc -l
echo "=== 最近10个错误 ==="
grep "ERROR" logs/webhook.log | tail -10
```

#### 2. 性能分析脚本
```bash
#!/bin/bash
echo "=== 消息处理统计 ==="
echo "总消息数: $(grep "直接推送" logs/webhook.log | wc -l)"
echo "成功消息数: $(grep "发送成功" logs/webhook.log | wc -l)"
echo "失败消息数: $(grep "发送失败" logs/webhook.log | wc -l)"
```

#### 3. 连接状态脚本
```bash
#!/bin/bash
echo "=== 连接状态检查 ==="
grep "连接" logs/webhook.log | tail -5
echo "=== 最后连接时间 ==="
grep "连接时间" logs/webhook.log | tail -1
```

## 🔍 配置检查

### 配置验证清单

#### 1. 基础配置检查
- [ ] `enabled` 设置为 `true`
- [ ] `url` 格式正确且可访问
- [ ] `secret` 长度≥16字符
- [ ] `timeout` 在合理范围内（5-30秒）

#### 2. 独立模式检查
- [ ] `useDirectStream` 设置为 `true`
- [ ] `useRedisSync` 设置为 `false`
- [ ] `independentMode` 设置为 `true`

#### 3. 消息过滤检查
- [ ] `messageTypes` 包含需要的消息类型
- [ ] `includeSelfMessage` 设置符合需求
- [ ] `wxId` 设置正确（如需要过滤特定账号）

### 配置验证脚本

```bash
#!/bin/bash
echo "=== Webhook配置验证 ==="

# 检查配置文件
if [ -f "webhook_config.json" ]; then
    echo "✅ 配置文件存在"
    cat webhook_config.json | jq '.enabled' | grep -q "true" && echo "✅ 启用状态正确" || echo "❌ 启用状态错误"
    cat webhook_config.json | jq '.useDirectStream' | grep -q "true" && echo "✅ 直接消息流启用" || echo "❌ 直接消息流未启用"
    cat webhook_config.json | jq '.independentMode' | grep -q "true" && echo "✅ 独立模式启用" || echo "❌ 独立模式未启用"
else
    echo "❌ 配置文件不存在"
fi
```

## 🌐 网络诊断

### 网络连接测试

#### 1. 基础连接测试
```bash
# 测试HTTP连接
curl -I "http://your-webhook-url.com/webhook"

# 测试HTTPS连接
curl -I "https://your-webhook-url.com/webhook"

# 测试DNS解析
nslookup your-webhook-url.com
```

#### 2. 端口连通性测试
```bash
# 测试80端口
telnet your-webhook-url.com 80

# 测试443端口
telnet your-webhook-url.com 443

# 测试自定义端口
telnet your-webhook-url.com 8080
```

#### 3. 路由追踪
```bash
# 追踪路由路径
traceroute your-webhook-url.com

# 查看网络延迟
ping -c 10 your-webhook-url.com
```

### 防火墙检查

#### 1. 本地防火墙
```bash
# 检查本地防火墙状态
sudo ufw status

# 检查iptables规则
sudo iptables -L
```

#### 2. 网络防火墙
```bash
# 检查出站连接
curl -I "http://httpbin.org/status/200"

# 检查入站连接
netstat -tlnp | grep :1238
```

## ⚡ 性能优化

### 性能监控指标

#### 1. 响应时间监控
```bash
# 监控平均响应时间
grep "响应时间" logs/webhook.log | awk '{sum+=$NF; count++} END {print "平均响应时间: " sum/count "ms"}'

# 监控最大响应时间
grep "响应时间" logs/webhook.log | awk '{if($NF>max) max=$NF} END {print "最大响应时间: " max "ms"}'
```

#### 2. 吞吐量监控
```bash
# 监控消息处理速率
grep "直接推送" logs/webhook.log | wc -l

# 监控成功率
success=$(grep "发送成功" logs/webhook.log | wc -l)
total=$(grep "直接推送" logs/webhook.log | wc -l)
echo "成功率: $((success * 100 / total))%"
```

#### 3. 资源使用监控
```bash
# 监控内存使用
ps aux | grep wechatpadpro | awk '{print "内存使用: " $6 "KB"}'

# 监控CPU使用
top -p $(pgrep wechatpadpro) -n 1 | grep wechatpadpro
```

### 性能优化建议

#### 1. 超时优化
```json
{
    "timeout": 5,        // 高性能环境
    "retryCount": 1      // 减少重试次数
}
```

#### 2. 消息过滤优化
```json
{
    "messageTypes": ["1", "3"],  // 只接收必要消息类型
    "includeSelfMessage": false   // 不包含自己消息
}
```

#### 3. 网络优化
```json
{
    "url": "https://your-webhook-url.com/webhook",  // 使用HTTPS
    "useDirectStream": true,                        // 启用直接消息流
    "independentMode": true                         // 启用独立模式
}
```

## 🚨 紧急修复

### 快速修复步骤

#### 1. 服务重启
```bash
# 停止服务
pkill wechatpadpro

# 启动服务
./wechatpadpro.exe
```

#### 2. 配置重置
```bash
# 备份当前配置
cp webhook_config.json webhook_config.json.backup

# 重置为默认配置
cat > webhook_config.json << EOF
{
    "url": "http://localhost:3000/webhook",
    "secret": "emergency_key_2024",
    "enabled": true,
    "timeout": 10,
    "retryCount": 3,
    "messageTypes": ["*"],
    "includeSelfMessage": true,
    "useDirectStream": true,
    "useRedisSync": false,
    "independentMode": true,
    "wxId": "",
    "lastSendTime": 0,
    "lastSendStatus": false,
    "totalSent": 0,
    "totalFailed": 0
}
EOF
```

#### 3. 连接重置
```bash
# 重置所有Webhook连接
curl -X POST "http://localhost:1238/webhook/ReinitializeAll"

# 强制启用直接消息流
curl -X POST "http://localhost:1238/webhook/ForceEnableDirectStream?key=your_key"
```

### 紧急联系信息

#### 技术支持
- **GitHub Issues**: https://github.com/WeChatPadPro/WeChatPadPro/issues
- **官方论坛**: https://bbs.knowhub.cloud/
- **技术文档**: https://docs.knowhub.cloud/

#### 日志收集
```bash
# 收集诊断信息
curl "http://localhost:1238/webhook/Diagnostics?key=your_key" > diagnostics.json

# 收集日志文件
tar -czf webhook_logs_$(date +%Y%m%d_%H%M%S).tar.gz logs/

# 收集配置信息
cp webhook_config.json config_backup_$(date +%Y%m%d_%H%M%S).json
```

---

**版本**: v1.0  
**更新时间**: 2025-07-21  
**适用版本**: WeChatPadPro v8.6.1+