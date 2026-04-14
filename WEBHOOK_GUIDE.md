# Webhook 連接與訊息處理指南

## 目錄

1. [簡介](#簡介)
2. [Webhook 配置](#webhook-配置)
3. [訊息接收](#訊息接收)
4. [訊息格式](#訊息格式)
5. [訊息處理](#訊息處理)
6. [安全性](#安全性)
7. [錯誤處理](#錯誤處理)
8. [最佳實踐](#最佳實踐)
9. [示例代碼](#示例代碼)
10. [常見問題](#常見問題)

## 簡介

Webhook 是一種基於 HTTP 的回調機制，允許我們的系統在特定事件發生時，向您預先配置的 URL 發送實時通知。通過 Webhook，您可以接收微信訊息、狀態變更等事件，實現系統間的實時數據同步。

### 主要功能

- 實時接收微信訊息（文字、圖片、語音、視頻等）
- 接收群組訊息和私聊訊息
- 接收狀態變更通知（登入狀態、好友請求等）
- 支持訊息過濾和自定義配置

## Webhook 配置

### 配置步驟

1. **創建 Webhook 配置**

   通過 API 創建 Webhook 配置：

   ```http
   POST /v1/webhook/Config
   ```

   請求體：
   ```json
   {
     "URL": "https://your-server.com/webhook/receiver",
     "Secret": "your_secret_key",
     "Enabled": true,
     "Timeout": 10,
     "RetryCount": 3,
     "MessageTypes": ["1", "3", "34", "43", "47", "49"],
     "IncludeSelfMessage": false
   }
   ```

   參數說明：
   - `URL`: 接收訊息的伺服器地址
   - `Secret`: 用於訊息簽名驗證的密鑰
   - `Enabled`: 是否啟用該 Webhook
   - `Timeout`: 請求超時時間（秒）
   - `RetryCount`: 失敗重試次數
   - `MessageTypes`: 接收的訊息類型（1=文字, 3=圖片, 34=語音, 43=視頻, 47=表情, 49=鏈接）
   - `IncludeSelfMessage`: 是否接收自己發送的訊息

2. **查看 Webhook 配置**

   ```http
   GET /v1/webhook/List
   ```

3. **測試 Webhook 連接**

   ```http
   GET /v1/webhook/Test
   ```

4. **更新 Webhook 配置**

   ```http
   PUT /v1/webhook/Update
   ```

## 訊息接收

### 伺服器要求

1. **公網可訪問**：您的伺服器必須能夠從公網訪問
2. **支持 HTTPS**：建議使用 HTTPS 確保安全性
3. **響應及時**：伺服器應在 5 秒內響應請求，否則會被視為失敗

### 接收流程

1. 當事件發生時，系統會向您配置的 URL 發送 POST 請求
2. 您的伺服器需要返回 HTTP 狀態碼 200 表示成功接收
3. 如果返回其他狀態碼或超時，系統會根據配置的重試次數進行重試

## 訊息格式

### 請求頭

```
Content-Type: application/json
X-Webhook-Timestamp: 1625123456
X-Webhook-Signature: abcdef1234567890...
```

### 請求體

```json
{
  "key": "wxid_abcdefg123456",
  "msgId": "7654321098765432123",
  "timestamp": 1625123456,
  "fromUser": "wxid_sender12345",
  "toUser": "wxid_receiver12345",
  "msgType": 1,
  "content": "你好，這是一條測試訊息",
  "isSelfMsg": false,
  "createTime": 1625123450,
  "isHistory": false
}
```

### 字段說明

| 字段 | 類型 | 說明 |
|------|------|------|
| key | string | 微信帳號唯一標識 |
| msgId | string | 訊息唯一ID |
| timestamp | int64 | 訊息推送時間戳 |
| fromUser | string | 發送者ID |
| toUser | string | 接收者ID |
| msgType | int | 訊息類型 |
| content | object | 訊息內容，根據訊息類型不同而不同 |
| isSelfMsg | boolean | 是否為自己發送的訊息 |
| createTime | int64 | 訊息創建時間戳 |
| isHistory | boolean | 是否為歷史訊息 |

### 訊息類型說明

| msgType | 說明 | content 格式 |
|---------|------|--------------|
| 1 | 文字訊息 | 字符串 |
| 3 | 圖片訊息 | 圖片URL或Base64 |
| 34 | 語音訊息 | 語音URL或Base64 |
| 43 | 視頻訊息 | 視頻URL或Base64 |
| 47 | 表情訊息 | 表情URL或Base64 |
| 49 | 鏈接訊息 | 包含標題、描述、URL等的對象 |
| 10000 | 系統通知 | 字符串 |
| 10002 | 撤回訊息 | 字符串 |

## 訊息處理

### 驗證簽名

每個 Webhook 請求都包含簽名，您應該驗證簽名以確保請求的合法性：

```python
import hmac
import hashlib

def verify_signature(request_body, timestamp, signature, secret):
    # 組合待簽名字符串
    string_to_sign = timestamp + ":" + request_body
    
    # 計算簽名
    computed_signature = hmac.new(
        secret.encode('utf-8'),
        string_to_sign.encode('utf-8'),
        hashlib.sha256
    ).hexdigest()
    
    # 比較簽名
    return hmac.compare_digest(computed_signature, signature)
```

### 處理不同類型的訊息

根據 `msgType` 字段處理不同類型的訊息：

```python
def process_message(message):
    msg_type = message.get("msgType")
    
    if msg_type == 1:  # 文字訊息
        process_text_message(message)
    elif msg_type == 3:  # 圖片訊息
        process_image_message(message)
    elif msg_type == 34:  # 語音訊息
        process_voice_message(message)
    elif msg_type == 43:  # 視頻訊息
        process_video_message(message)
    elif msg_type == 47:  # 表情訊息
        process_emoji_message(message)
    elif msg_type == 49:  # 鏈接訊息
        process_link_message(message)
    elif msg_type == 10000:  # 系統通知
        process_system_notification(message)
    elif msg_type == 10002:  # 撤回訊息
        process_recall_message(message)
    else:
        process_unknown_message(message)
```

### 訊息去重

系統已實現訊息去重機制，但您仍應在自己的服務中實現訊息去重，以防止重複處理：

```python
processed_messages = set()

def is_duplicate(message_id):
    if message_id in processed_messages:
        return True
    
    processed_messages.add(message_id)
    # 維護一個固定大小的集合，避免內存泄漏
    if len(processed_messages) > 10000:
        processed_messages.pop()
    
    return False

def process_message(message):
    message_id = message.get("msgId")
    
    if is_duplicate(message_id):
        print(f"跳過重複訊息: {message_id}")
        return
    
    # 處理訊息...
```

## 安全性

### 簽名驗證

始終驗證請求簽名，確保請求來自合法來源。

### HTTPS

強烈建議使用 HTTPS 協議保護數據傳輸安全。

### IP 白名單

可以設置 IP 白名單，只接受來自特定 IP 地址的請求。

### 敏感數據處理

妥善處理和存儲接收到的敏感數據，遵守相關的數據保護法規。

## 錯誤處理

### 重試機制

系統會根據配置的重試次數，在請求失敗時進行重試。重試間隔會隨著重試次數增加而增加（指數退避）。

### 日誌記錄

記錄所有 Webhook 請求和響應，以便排查問題：

```python
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("webhook")

def handle_webhook(request):
    try:
        # 記錄請求
        logger.info(f"Received webhook: {request.body}")
        
        # 處理請求
        process_message(request.json())
        
        # 記錄成功
        logger.info("Webhook processed successfully")
        
        return {"status": "success"}, 200
    except Exception as e:
        # 記錄錯誤
        logger.error(f"Error processing webhook: {str(e)}", exc_info=True)
        
        return {"status": "error", "message": str(e)}, 500
```

## 最佳實踐

1. **快速響應**：儘快返回 HTTP 200 狀態碼，將耗時操作放入後台任務隊列
2. **訊息去重**：實現訊息去重機制，避免重複處理
3. **異步處理**：使用異步處理框架提高並發能力
4. **監控告警**：設置監控和告警，及時發現問題
5. **定期清理**：定期清理歷史數據，避免存儲空間耗盡
6. **容錯處理**：妥善處理各種異常情況，確保服務穩定性

## 示例代碼

### Python (Flask) 示例

```python
from flask import Flask, request, jsonify
import hmac
import hashlib
import json
import time
import threading
import queue

app = Flask(__name__)

# 訊息隊列
message_queue = queue.Queue()

# 配置
WEBHOOK_SECRET = "your_secret_key"
PROCESSED_MESSAGES = set()
MAX_QUEUE_SIZE = 1000

# 訊息處理線程
def message_processor():
    while True:
        try:
            message = message_queue.get(timeout=1)
            # 處理訊息
            process_message(message)
            message_queue.task_done()
        except queue.Empty:
            time.sleep(0.1)
        except Exception as e:
            print(f"Error processing message: {e}")

# 啟動處理線程
threading.Thread(target=message_processor, daemon=True).start()

@app.route('/webhook/receiver', methods=['POST'])
def webhook_receiver():
    # 獲取請求數據
    timestamp = request.headers.get('X-Webhook-Timestamp')
    signature = request.headers.get('X-Webhook-Signature')
    body = request.get_data(as_text=True)
    
    # 驗證簽名
    if not verify_signature(body, timestamp, signature, WEBHOOK_SECRET):
        return jsonify({"status": "error", "message": "Invalid signature"}), 401
    
    # 解析訊息
    try:
        message = json.loads(body)
    except json.JSONDecodeError:
        return jsonify({"status": "error", "message": "Invalid JSON"}), 400
    
    # 檢查訊息ID
    message_id = message.get("msgId")
    if not message_id:
        return jsonify({"status": "error", "message": "Missing message ID"}), 400
    
    # 檢查是否重複訊息
    if message_id in PROCESSED_MESSAGES:
        return jsonify({"status": "success", "message": "Duplicate message"}), 200
    
    # 添加到已處理集合
    PROCESSED_MESSAGES.add(message_id)
    if len(PROCESSED_MESSAGES) > MAX_QUEUE_SIZE:
        PROCESSED_MESSAGES.pop()
    
    # 將訊息加入隊列異步處理
    try:
        message_queue.put(message, block=False)
    except queue.Full:
        return jsonify({"status": "error", "message": "Server busy"}), 503
    
    # 立即返回成功
    return jsonify({"status": "success"}), 200

def process_message(message):
    # 根據訊息類型處理
    msg_type = message.get("msgType")
    
    if msg_type == 1:  # 文字訊息
        handle_text_message(message)
    elif msg_type == 3:  # 圖片訊息
        handle_image_message(message)
    # 處理其他訊息類型...

def handle_text_message(message):
    content = message.get("content")
    from_user = message.get("fromUser")
    print(f"收到來自 {from_user} 的文字訊息: {content}")
    # 處理文字訊息的業務邏輯...

def handle_image_message(message):
    content = message.get("content")
    from_user = message.get("fromUser")
    print(f"收到來自 {from_user} 的圖片訊息: {content}")
    # 處理圖片訊息的業務邏輯...

def verify_signature(body, timestamp, signature, secret):
    if not all([body, timestamp, signature, secret]):
        return False
    
    string_to_sign = f"{timestamp}:{body}"
    computed_signature = hmac.new(
        secret.encode('utf-8'),
        string_to_sign.encode('utf-8'),
        hashlib.sha256
    ).hexdigest()
    
    return hmac.compare_digest(computed_signature, signature)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
```

### Node.js (Express) 示例

```javascript
const express = require('express');
const crypto = require('crypto');
const bodyParser = require('body-parser');

const app = express();
app.use(bodyParser.json());

// 配置
const WEBHOOK_SECRET = 'your_secret_key';
const PROCESSED_MESSAGES = new Set();
const MAX_SET_SIZE = 10000;

// 驗證簽名
function verifySignature(body, timestamp, signature, secret) {
  if (!body || !timestamp || !signature || !secret) {
    return false;
  }
  
  const stringToSign = `${timestamp}:${JSON.stringify(body)}`;
  const computedSignature = crypto
    .createHmac('sha256', secret)
    .update(stringToSign)
    .digest('hex');
  
  return crypto.timingSafeEqual(
    Buffer.from(computedSignature),
    Buffer.from(signature)
  );
}

// 訊息處理
function processMessage(message) {
  const msgType = message.msgType;
  
  switch (msgType) {
    case 1:
      handleTextMessage(message);
      break;
    case 3:
      handleImageMessage(message);
      break;
    // 處理其他訊息類型...
    default:
      console.log(`未知訊息類型: ${msgType}`);
  }
}

function handleTextMessage(message) {
  const content = message.content;
  const fromUser = message.fromUser;
  console.log(`收到來自 ${fromUser} 的文字訊息: ${content}`);
  // 處理文字訊息的業務邏輯...
}

function handleImageMessage(message) {
  const content = message.content;
  const fromUser = message.fromUser;
  console.log(`收到來自 ${fromUser} 的圖片訊息: ${content}`);
  // 處理圖片訊息的業務邏輯...
}

// Webhook 接收端點
app.post('/webhook/receiver', (req, res) => {
  // 獲取請求數據
  const timestamp = req.headers['x-webhook-timestamp'];
  const signature = req.headers['x-webhook-signature'];
  const body = req.body;
  
  // 驗證簽名
  if (!verifySignature(body, timestamp, signature, WEBHOOK_SECRET)) {
    return res.status(401).json({ status: 'error', message: 'Invalid signature' });
  }
  
  // 檢查訊息ID
  const messageId = body.msgId;
  if (!messageId) {
    return res.status(400).json({ status: 'error', message: 'Missing message ID' });
  }
  
  // 檢查是否重複訊息
  if (PROCESSED_MESSAGES.has(messageId)) {
    return res.status(200).json({ status: 'success', message: 'Duplicate message' });
  }
  
  // 添加到已處理集合
  PROCESSED_MESSAGES.add(messageId);
  if (PROCESSED_MESSAGES.size > MAX_SET_SIZE) {
    const firstItem = PROCESSED_MESSAGES.values().next().value;
    PROCESSED_MESSAGES.delete(firstItem);
  }
  
  // 異步處理訊息
  setImmediate(() => {
    try {
      processMessage(body);
    } catch (error) {
      console.error(`處理訊息錯誤: ${error}`);
    }
  });
  
  // 立即返回成功
  return res.status(200).json({ status: 'success' });
});

// 啟動伺服器
const PORT = process.env.PORT || 5000;
app.listen(PORT, () => {
  console.log(`Webhook 伺服器運行在端口 ${PORT}`);
});
```

## 常見問題

### 1. 為什麼我沒有收到 Webhook 訊息？

可能的原因：
- Webhook URL 不可訪問（檢查防火牆、網絡設置）
- 伺服器響應超時（確保在 5 秒內返回響應）
- 簽名驗證失敗（檢查密鑰是否正確）
- 訊息類型過濾（檢查 MessageTypes 配置）

### 2. 如何處理大量並發訊息？

- 使用異步處理框架（如 Celery、Bull 等）
- 實現訊息隊列（RabbitMQ、Redis 等）
- 水平擴展接收伺服器
- 優化數據庫操作

### 3. 如何確保訊息不丟失？

- 實現訊息持久化存儲
- 使用可靠的訊息隊列
- 實現確認機制
- 定期檢查訊息處理狀態

### 4. 如何測試 Webhook 接收功能？

- 使用 `/v1/webhook/Test` API 發送測試訊息
- 使用工具如 Postman、curl 模擬請求
- 設置本地測試環境（使用 ngrok 等工具暴露本地服務）

### 5. 為什麼會收到重複訊息？

- 系統重試機制（請求失敗時會重試）
- 網絡問題導致重複發送
- 系統重啟後重新處理未確認訊息

### 6. 如何優化 Webhook 性能？

- 快速響應，異步處理
- 使用連接池
- 實現緩存機制
- 定期清理過期數據
- 使用高效的數據結構和算法 