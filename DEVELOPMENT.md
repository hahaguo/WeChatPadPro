# WeChatPadPro 开发指南

## 项目结构

```
WeChatPadPro/
├── cmd/                    # 主程序入口
│   └── main.go            # 主程序文件
├── internal/              # 内部代码
│   ├── config/           # 配置管理
│   │   └── config.go
│   ├── database/         # 数据库初始化
│   │   └── database.go
│   ├── model/            # 数据模型
│   │   └── model.go
│   ├── repository/       # 数据访问层
│   │   └── repository.go
│   ├── service/          # 业务逻辑层
│   │   └── service.go
│   ├── handler/          # HTTP处理器
│   │   ├── login.go
│   │   ├── webhook.go
│   │   ├── message.go
│   │   ├── health.go
│   │   └── sse.go
│   └── middleware/       # 中间件
│       └── middleware.go
├── pkg/                  # 公共包
│   └── logger/          # 日志工具
│       └── logger.go
├── static/              # 静态文件
├── assets/              # 资源文件
├── deploy/              # 部署配置
├── go.mod              # Go模块定义
├── go.sum              # Go模块校验
└── Dockerfile          # Docker构建文件
```

## 快速开始

### 安装依赖

```bash
go mod download
```

### 配置环境变量

复制 `.env.example` 到 `.env` 并修改配置：

```bash
cp .env.example .env
```

### 运行开发服务器

```bash
go run cmd/main.go
```

### 构建生产版本

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o wechatpadpro-linux-amd64 cmd/main.go

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o wechatpadpro-macos-amd64 cmd/main.go

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o wechatpadpro-windows-amd64.exe cmd/main.go
```

## API 文档

### 登录相关

#### 生成授权密钥
```
GET /api/login/GenAuthKey2?key=ADMIN_KEY&count=1&days=365
```

#### 获取登录二维码
```
POST /api/login/qr/newx
Content-Type: application/json

{
  "proxy": "http://proxy:port",
  "deviceName": "iPhone",
  "deviceId": "device_id"
}
```

#### 检查扫码状态
```
GET /api/login/CheckLoginStatus?key=your-uuid
```

#### 验证码处理
```
POST /api/login/verify/auto?key=your-uuid
Content-Type: application/json

{
  "uuid": "your-uuid",
  "code": "123456"
}
```

### Webhook 相关

#### 创建配置
```
POST /v1/webhook/Config
Content-Type: application/json

{
  "URL": "https://your-server.com/webhook",
  "Secret": "your_secret_key",
  "Enabled": true,
  "Timeout": 10,
  "RetryCount": 3,
  "MessageTypes": ["*"],
  "IncludeSelfMessage": false
}
```

#### 查看配置
```
GET /webhook/List
```

#### 测试连接
```
GET /webhook/Test
```

### 消息相关

#### 发送文本消息
```
POST /api/message/sendText
Content-Type: application/json

{
  "key": "your_key",
  "toUser": "wxid_receiver",
  "content": "Hello, World!"
}
```

## 开发规范

### 代码风格

- 遵循 [Effective Go](https://go.dev/doc/effective_go) 指南
- 使用 `gofmt` 格式化代码
- 添加必要的注释
- 使用有意义的变量名

### 提交规范

```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式
refactor: 重构
test: 测试相关
chore: 构建/工具
```

### 分支管理

- `main` - 主分支，稳定版本
- `develop` - 开发分支
- `feature/*` - 功能分支
- `bugfix/*` - 修复分支