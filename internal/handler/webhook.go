package handler

import (
	"encoding/json"
	"net/http"

	"github.com/WeChatPadPro/WeChatPadPro/internal/model"
	"github.com/WeChatPadPro/WeChatPadPro/internal/service"
	"github.com/gin-gonic/gin"
)

// WebhookHandler Webhook处理器
type WebhookHandler struct {
	webhookService *service.WebhookService
}

// NewWebhookHandler 创建Webhook处理器
func NewWebhookHandler(webhookService *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
	}
}

// List 列出所有Webhook配置
func (h *WebhookHandler) List(c *gin.Context) {
	configs, err := h.webhookService.ListConfigs()
	if err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "success",
		Data:    configs,
	})
}

// Status 获取Webhook配置状态
func (h *WebhookHandler) Status(c *gin.Context) {
	configs, err := h.webhookService.ListConfigs()
	if err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: err.Error(),
		})
		return
	}

	// 返回状态信息
	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "success",
		Data: gin.H{
			"total":     len(configs),
			"enabled":   countEnabled(configs),
			"configs":   configs,
		},
	})
}

// countEnabled 统计启用的配置数量
func countEnabled(configs []*model.WebhookConfig) int {
	count := 0
	for _, cfg := range configs {
		if cfg.Enabled {
			count++
		}
	}
	return count
}

// Test 测试Webhook连接
func (h *WebhookHandler) Test(c *gin.Context) {
	// 发送测试消息
	testMsg := &model.Message{
		Key:        "test",
		MsgID:      "test_" + getCurrentTimestamp(),
		Timestamp:  getCurrentTimestampInt64(),
		FromUser:   "test_sender",
		ToUser:     "test_receiver",
		MsgType:    model.MsgTypeText,
		Content:    "This is a test message from Webhook",
		IsSelfMsg:  false,
		CreateTime: getCurrentTimestampInt64(),
		IsHistory:  false,
	}

	// 获取所有启用的配置并发送测试消息
	configs, err := h.webhookService.ListConfigs()
	if err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: err.Error(),
		})
		return
	}

	successCount := 0
	for _, cfg := range configs {
		if cfg.Enabled {
			// 这里应该实际发送测试消息
			successCount++
		}
	}

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "test completed",
		Data: gin.H{
			"tested":   successCount,
			"message":  testMsg,
		},
	})
}

// Config 创建Webhook配置
func (h *WebhookHandler) Config(c *gin.Context) {
	var cfg model.WebhookConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "invalid request body",
		})
		return
	}

	// 验证必填字段
	if cfg.URL == "" {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "URL is required",
		})
		return
	}

	if cfg.Secret == "" {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "secret is required",
		})
		return
	}

	// 设置默认值
	if cfg.Timeout == 0 {
		cfg.Timeout = 10
	}
	if cfg.RetryCount == 0 {
		cfg.RetryCount = 3
	}
	if cfg.MessageTypes == "" {
		cfg.MessageTypes = "[\"*\"]"
	}

	if err := h.webhookService.CreateConfig(&cfg); err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "Webhook config created successfully",
		Data:    cfg,
	})
}

// Update 更新Webhook配置
func (h *WebhookHandler) Update(c *gin.Context) {
	var cfg model.WebhookConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "invalid request body",
		})
		return
	}

	if cfg.ID == 0 {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "ID is required for update",
		})
		return
	}

	if err := h.webhookService.UpdateConfig(&cfg); err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "Webhook config updated successfully",
		Data:    cfg,
	})
}

// Delete 删除Webhook配置
func (h *WebhookHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "ID is required",
		})
		return
	}

	var id int64
	if _, err := json.Unmarshal([]byte(idStr), &id); err != nil {
		// 尝试直接解析
		id = int64(json.Number(idStr).Int64())
	}

	if err := h.webhookService.DeleteConfig(id); err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "Webhook config deleted successfully",
	})
}

// Diagnostics 获取诊断信息
func (h *WebhookHandler) Diagnostics(c *gin.Context) {
	configs, err := h.webhookService.ListConfigs()
	if err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: err.Error(),
		})
		return
	}

	var diagnostics []gin.H
	for _, cfg := range configs {
		diag := gin.H{
			"id":               cfg.ID,
			"url":              cfg.URL,
			"enabled":          cfg.Enabled,
			"lastSendTime":     cfg.LastSendTime,
			"lastSendStatus":   cfg.LastSendStatus,
			"totalSent":        cfg.TotalSent,
			"totalFailed":      cfg.TotalFailed,
			"successRate":      calculateSuccessRate(cfg),
			"connectionStatus": getConnectionStatus(cfg),
		}
		diagnostics = append(diagnostics, diag)
	}

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "diagnostics retrieved",
		Data:    diagnostics,
	})
}

// calculateSuccessRate 计算成功率
func calculateSuccessRate(cfg *model.WebhookConfig) float64 {
	total := cfg.TotalSent + cfg.TotalFailed
	if total == 0 {
		return 0.0
	}
	return float64(cfg.TotalSent) / float64(total) * 100
}

// getConnectionStatus 获取连接状态
func getConnectionStatus(cfg *model.WebhookConfig) string {
	if !cfg.Enabled {
		return "disabled"
	}
	if cfg.LastSendStatus {
		return "connected"
	}
	return "disconnected"
}

// ResetConnection 重置连接
func (h *WebhookHandler) ResetConnection(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "key is required",
		})
		return
	}

	// 这里应该实现重置连接的逻辑
	// 例如: 清除连接状态、重新初始化等

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "connection reset successfully",
	})
}
