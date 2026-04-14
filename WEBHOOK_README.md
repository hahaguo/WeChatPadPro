# Flask Webhook Receiver

一个轻量级的 Flask Webhook 接收服务，支持签名验证、日志记录、消息内容格式化及过滤自身消息等功能，适用于各种内部系统或第三方服务的消息回调处理。

## 📦 特性

- 支持 Webhook 签名验证（HMAC-SHA256）
- 日志输出到文件和控制台
- 可配置接收消息类型及是否接收自己发送的消息
- 自动格式化显示接收到的消息内容
- 支持失败重试配置和超时限制

## 🧾 示例配置

```json
{
  "URL": "http://192.168.0.101:8000/webhook",
  "Secret": "your_secret_key",
  "Enabled": true,
  "Timeout": 10,
  "RetryCount": 3,
  "MessageTypes": ["*"],
  "IncludeSelfMessage": true
}
```

## 🚀 快速开始

### 1. 安装依赖

```bash
pip install flask
```

### 2. 启动服务

```bash
python webhook.py
```

服务默认监听在 `http://0.0.0.0:8000/webhook`

## 💻 示例代码

```python
from flask import Flask, request, jsonify
import hmac
import hashlib
import logging
from datetime import datetime

app = Flask(__name__)

WEBHOOK_SECRET = "your_secret_key"
INCLUDE_SELF_MESSAGE = True

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    handlers=[
        logging.FileHandler("webhook.log", encoding='utf-8'),
        logging.StreamHandler()
    ]
)

def verify_signature(data, signature, secret, timestamp):
    mac = hmac.new(secret.encode('utf-8'), digestmod=hashlib.sha256)
    mac.update(timestamp.encode('utf-8'))
    mac.update(data)
    expected = mac.hexdigest()
    return hmac.compare_digest(expected, signature)

def format_message(data):
    msg_id = data.get("msgId")
    from_user = data.get("fromUser")
    to_user = data.get("toUser")
    msg_type = data.get("msgType")
    timestamp = data.get("timestamp")
    content = data.get("content", "")

    if isinstance(content, dict):
        content = content.get("str", "")
    elif not isinstance(content, str):
        content = str(content)

    try:
        time_str = datetime.fromtimestamp(timestamp).strftime("%Y-%m-%d %H:%M:%S")
    except:
        time_str = str(timestamp)

    return f"""✅ Received message:
🆔 MsgID: {msg_id}
👤 From: {from_user}
🎯 To: {to_user}
🕒 Time: {time_str}
📨 Type: {msg_type}
💬 Content: {content}
"""

@app.route("/webhook", methods=["POST"])
def webhook():
    raw_data = request.data
    headers = request.headers

    signature = headers.get('X-Webhook-Signature')
    timestamp = headers.get('X-Webhook-Timestamp')

    if WEBHOOK_SECRET:
        if not signature or not timestamp:
            return jsonify({"error": "Missing signature or timestamp"}), 400
        if not verify_signature(raw_data, signature, WEBHOOK_SECRET, timestamp):
            logging.warning("❌ Signature verification failed")
            return jsonify({"error": "Invalid signature"}), 403

    try:
        data = request.get_json(force=True)

        if not INCLUDE_SELF_MESSAGE and data.get('isSelfMsg'):
            return jsonify({"status": "ignored", "reason": "self message skipped"}), 200

        log_msg = format_message(data)
        logging.info(log_msg)

        return jsonify({"status": "ok"}), 200

    except Exception as e:
        logging.exception("❌ Error processing webhook:")
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    logging.info("🚀 Webhook server is running on port 8000...")
    app.run(host="0.0.0.0", port=8000)
```

## 📝 日志输出格式

每条接收的消息会记录类似格式：

```
✅ Received message:
🆔 MsgID: 123456
👤 From: user_a
🎯 To: user_b
🕒 Time: 2025-07-12 12:00:00
📨 Type: text
💬 Content: Hello world!
```

## ⚙️ 配置说明

| 配置项               | 描述                                      |
|----------------------|-------------------------------------------|
| `Secret`             | 用于签名验证的密钥（可选但推荐）           |
| `IncludeSelfMessage` | 是否处理自己发送的消息                     |
| `Timeout`            | 推送请求超时时间（1~30秒）                |
| `RetryCount`         | 推送失败的最大重试次数                     |
| `MessageTypes`       | 接收的消息类型，`["*"]` 表示全部            |

## 📄 License

本项目可自由修改与部署，遵循 MIT 开源协议。
