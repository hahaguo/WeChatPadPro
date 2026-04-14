package handler

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/WeChatPadPro/WeChatPadPro/pkg/logger"
	"github.com/gin-gonic/gin"
)

// SSEHandler Server-Sent Events 处理器
type SSEHandler struct {
	clients map[string]chan SSEMessage
	mu      sync.RWMutex
}

// SSEMessage SSE消息
type SSEMessage struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	Data  string `json:"data"`
}

// NewSSEHandler 创建SSE处理器
func NewSSEHandler() *SSEHandler {
	return &SSEHandler{
		clients: make(map[string]chan SSEMessage),
	}
}

// SSEHandlerFunc 处理SSE连接
func SSEHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(400, gin.H{"error": "key is required"})
		return
	}

	logger.Info("New SSE connection established for key: %s", key)

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// 获取响应写入器
	clientGone := c.Request.Context().Done()

	// 发送心跳消息
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	// 模拟发送一些初始消息
	sendSSE(c, "connected", `{"status":"connected","time":"`+time.Now().Format(time.RFC3339)+`"}`)

	for {
		select {
		case <-clientGone:
			logger.Info("SSE client disconnected for key: %s", key)
			return
		case <-ticker.C:
			// 发送心跳
			sendSSE(c, "ping", `{"time":"`+time.Now().Format(time.RFC3339)+`"}`)
		}
	}
}

// sendSSE 发送SSE消息
func sendSSE(c *gin.Context, event, data string) {
	fmt.Fprintf(c.Writer, "event: %s\n", event)
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)

	if f, ok := c.Writer.(http.Flusher); ok {
		f.Flush()
	}
}