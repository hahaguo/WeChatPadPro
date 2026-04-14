 # Webhook配置详细说明

## 配置文件：webhook_config.json

### 基础配置

```json
{
    "url": "http://192.168.0.101:8000/webhook",           // Webhook接收端点URL
    "secret": "your_secret_key_2024_enhanced_v1",          // 用于消息签名验证的密钥（建议32字符以上）
    "enabled": true,                                        // 是否启用Webhook功能
    "timeout": 10,                                          // HTTP请求超时时间（秒）
    "retryCount": 3,                                        // 失败重试次数
    "messageTypes": ["*"],                                  // 要接收的消息类型
    "includeSelfMessage": true,                             // 是否包含自己发送的消息
    "wxId": "",                                             // 微信ID，用于消息过滤
    "useDirectStream": true,                                // 启用直接消息流（零延迟推送）
    "useRedisSync": false,                                  // 是否使用Redis同步（独立模式建议关闭）
    "independentMode": true,                                // 启用完全独立模式（与WebSocket消息流分离）
    "lastSendTime": 0,                                      // 最后发送消息的时间戳（系统自动维护）
    "lastSendStatus": false,                                // 最后一次发送的状态（系统自动维护）
    "totalSent": 0,                                         // 总发送消息数量（系统自动维护）
    "totalFailed": 0                                        // 总失败消息数量（系统自动维护）
}
```

## 配置项详解

### 🔧 基础配置

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `url` | string | ✅ | - | Webhook接收端点URL |
| `secret` | string | ✅ | - | 用于消息签名验证的密钥 |
| `enabled` | boolean | ✅ | true | 是否启用Webhook功能 |

### 🌐 网络配置

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `timeout` | integer | ✅ | 10 | HTTP请求超时时间（秒） |
| `retryCount` | integer | ✅ | 3 | 失败重试次数 |

### 📨 消息过滤配置

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `messageTypes` | array | ✅ | ["*"] | 要接收的消息类型 |
| `includeSelfMessage` | boolean | ✅ | true | 是否包含自己发送的消息 |
| `wxId` | string | ❌ | "" | 微信ID，用于消息过滤 |

### ⚡ 独立模式配置

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `useDirectStream` | boolean | ✅ | true | 启用直接消息流（零延迟推送） |
| `useRedisSync` | boolean | ✅ | false | 是否使用Redis同步（独立模式建议关闭） |
| `independentMode` | boolean | ✅ | true | 启用完全独立模式（与WebSocket消息流分离） |

### 📊 统计信息字段（系统自动维护）

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `lastSendTime` | int64 | ❌ | 0 | 最后发送消息的时间戳 |
| `lastSendStatus` | boolean | ❌ | false | 最后一次发送的状态 |
| `totalSent` | int64 | ❌ | 0 | 总发送消息数量 |
| `totalFailed` | int64 | ❌ | 0 | 总失败消息数量 |

## 消息类型说明

### 支持的消息类型

```json
{
    "messageTypes": ["*"]           // 接收所有消息类型
}
```

或者指定特定类型：

```json
{
    "messageTypes": ["1", "3", "47", "49"]  // 只接收文本、图片、表情、链接消息
}
```

### 消息类型对照表

| 类型值 | 说明 | 示例 |
|--------|------|------|
| `"*"` | 所有消息类型 | - |
| `"1"` | 文本消息 | 聊天文字 |
| `"3"` | 图片消息 | 照片、截图 |
| `"47"` | 表情消息 | 微信表情 |
| `"49"` | 链接消息 | 分享链接 |
| `"10000"` | 系统消息 | 群通知、撤回消息 |

## 配置示例

### 示例1：客服系统配置
```json
{
    "url": "https://customer-service.company.com/webhook",
    "secret": "cs_webhook_secure_key_2024",
    "enabled": true,
    "timeout": 8,
    "retryCount": 3,
    "messageTypes": ["1", "3", "47", "49"],
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
```

### 示例2：营销自动化配置
```json
{
    "url": "https://marketing-bot.company.com/webhook",
    "secret": "marketing_automation_key_2024",
    "enabled": true,
    "timeout": 5,
    "retryCount": 2,
    "messageTypes": ["*"],
    "includeSelfMessage": false,
    "useDirectStream": true,
    "useRedisSync": false,
    "independentMode": true,
    "wxId": "",
    "lastSendTime": 0,
    "lastSendStatus": false,
    "totalSent": 0,
    "totalFailed": 0
}
```

### 示例3：多账号管理配置
```json
{
    "url": "https://multi-account.company.com/webhook",
    "secret": "multi_account_key_2024",
    "enabled": true,
    "timeout": 10,
    "retryCount": 3,
    "messageTypes": ["*"],
    "includeSelfMessage": true,
    "useDirectStream": true,
    "useRedisSync": false,
    "independentMode": true,
    "wxId": "wxid_abc123def456",
    "lastSendTime": 0,
    "lastSendStatus": false,
    "totalSent": 0,
    "totalFailed": 0
}
```

### 示例4：高性能生产环境配置
```json
{
    "url": "https://high-performance.company.com/webhook",
    "secret": "hp_production_key_2024",
    "enabled": true,
    "timeout": 3,
    "retryCount": 1,
    "messageTypes": ["1", "3"],
    "includeSelfMessage": false,
    "useDirectStream": true,
    "useRedisSync": false,
    "independentMode": true,
    "wxId": "",
    "lastSendTime": 0,
    "lastSendStatus": false,
    "totalSent": 0,
    "totalFailed": 0
}
```

## 最佳实践

### 🔒 安全配置
- 使用HTTPS协议（生产环境）
- 使用强密钥（≥32字符）
- 定期更换密钥
- 实现签名验证

### ⚡ 性能优化
- 合理设置超时时间（5-10秒）
- 适当的重试次数（1-3次）
- 启用独立模式
- 过滤不需要的消息类型

### 🛡️ 稳定性保证
- 实现幂等性处理
- 添加消息去重机制
- 设置健康检查
- 监控Webhook状态

## 故障排查

### 常见问题

1. **Webhook无法接收消息**
   - 检查URL是否可访问
   - 确认配置是否启用（enabled: true）
   - 验证消息类型过滤是否正确

2. **消息延迟或丢失**
   - 确保useDirectStream为true
   - 确保useRedisSync为false
   - 确保independentMode为true

3. **频繁重试失败**
   - 检查服务端是否正常响应
   - 调整超时时间
   - 减少重试次数

4. **签名验证失败**
   - 确认密钥是否正确
   - 检查签名算法是否一致
   - 验证时间戳是否在有效范围内

## 版本信息

- **配置文件版本**: v2.0
- **更新时间**: 2025-07-21
- **适用版本**: WeChatPadPro v8.6.1+
- **配置格式**: JSON
