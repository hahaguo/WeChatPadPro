from flask import Flask, request, jsonify
import hmac
import hashlib
import time
import json
import logging
from datetime import datetime

app = Flask(__name__)

WEBHOOK_SECRET = "wh_sk_2024_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0"
INCLUDE_SELF_MESSAGE = True

# 日志配置
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    handlers=[
        logging.FileHandler("webhook.log", encoding='utf-8'),
        logging.StreamHandler()
    ]
)

# 签名验证
def verify_signature(data, signature, secret, timestamp):
    mac = hmac.new(secret.encode('utf-8'), digestmod=hashlib.sha256)
    mac.update(timestamp.encode('utf-8'))
    mac.update(data)
    expected_signature = mac.hexdigest()
    return hmac.compare_digest(expected_signature, signature)

# 自动格式化所有字段
def pretty_format(data, indent=0):
    spacing = '  ' * indent
    if isinstance(data, dict):
        result = ""
        for key, value in data.items():
            result += f"{spacing}- {key}:"
            if isinstance(value, (dict, list)):
                result += "\n" + pretty_format(value, indent + 1)
            else:
                result += f" {value}\n"
        return result
    elif isinstance(data, list):
        result = ""
        for idx, item in enumerate(data):
            result += f"{spacing}- [{idx}]:\n" + pretty_format(item, indent + 1)
        return result
    else:
        return f"{spacing}{data}\n"

# 格式化日志输出
def format_message(data):
    timestamp = data.get("timestamp", None)
    try:
        if isinstance(timestamp, (int, float)):
            time_str = datetime.fromtimestamp(timestamp).strftime("%Y-%m-%d %H:%M:%S")
        else:
            time_str = str(timestamp)
    except Exception:
        time_str = str(timestamp)

    return (
        f"\n✅ Received webhook message at {time_str}:\n"
        f"{pretty_format(data)}"
    )

# Webhook 接口
@app.route('/webhook', methods=['POST', 'HEAD'])
def webhook():
    if request.method == 'HEAD':
        # 健康检查，直接返回200
        return '', 200

    raw_data = request.data
    headers = request.headers

    # 记录请求头和部分body（可选，生产可注释）
    logging.info(f"[Request Headers] {dict(headers)}")
    try:
        logging.info(f"[Request Body] {raw_data[:500].decode('utf-8', errors='ignore')}")
    except Exception as e:
        logging.warning(f"[Request Body decode error]: {e}")

    signature = headers.get('X-Webhook-Signature')
    timestamp = headers.get('X-Webhook-Timestamp')

    if WEBHOOK_SECRET:
        if not signature or not timestamp:
            logging.warning("❌ Missing signature or timestamp")
            return jsonify({"status": "error", "message": "Missing signature or timestamp"}), 400
        try:
            if not verify_signature(raw_data, signature, WEBHOOK_SECRET, timestamp):
                logging.warning("❌ Signature verification failed")
                return jsonify({"status": "error", "message": "Invalid signature"}), 403
        except Exception as e:
            logging.error(f"❌ Signature verify exception: {e}")
            return jsonify({"status": "error", "message": "Signature verify exception"}), 400

    try:
        data = request.get_json(force=True)

        if not INCLUDE_SELF_MESSAGE and data.get('isSelfMsg'):
            return jsonify({"status": "ignored", "reason": "self message skipped"}), 200

        formatted = format_message(data)
        logging.info(formatted)

        return jsonify({"status": "ok"}), 200

    except Exception as e:
        logging.exception("❌ Error processing webhook:")
        return jsonify({"status": "error", "message": str(e)}), 500

if __name__ == '__main__':
    logging.info("🚀 Webhook server is running on port 8000...")
    app.run(host='0.0.0.0', port=8000)

