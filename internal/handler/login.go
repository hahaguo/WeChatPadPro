package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/WeChatPadPro/WeChatPadPro/internal/config"
	"github.com/WeChatPadPro/WeChatPadPro/internal/model"
	"github.com/WeChatPadPro/WeChatPadPro/internal/repository"
	"github.com/WeChatPadPro/WeChatPadPro/internal/service"
	"github.com/WeChatPadPro/WeChatPadPro/pkg/logger"
	"github.com/gin-gonic/gin"
)

// LoginHandler 登录处理器
type LoginHandler struct {
	authService *service.AuthService
	userRepo    *repository.UserRepository
	deviceRepo  *repository.DeviceRepository
	cfg         *config.Config
}

// NewLoginHandler 创建登录处理器
func NewLoginHandler(authService *service.AuthService, userRepo *repository.UserRepository, cfg *config.Config) *LoginHandler {
	return &LoginHandler{
		authService: authService,
		userRepo:    userRepo,
		cfg:         cfg,
	}
}

// GenAuthKey 生成授权密钥
// GET /api/login/GenAuthKey2?key=ADMIN_KEY&count=1&days=365
func (h *LoginHandler) GenAuthKey(c *gin.Context) {
	key := c.Query("key")
	count := 1
	days := 365

	if key == "" {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "missing admin key",
		})
		return
	}

	if !h.authService.VerifyAdminKey(key) {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "invalid admin key",
		})
		return
	}

	keys, err := h.authService.GenerateLicenseKeys(count, days)
	if err != nil {
		logger.Error("Failed to generate license keys: %v", err)
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, model.GenerateAuthKeyResponse{
		Code: model.StatusSuccess,
		Data: keys,
	})
}

// GetQRCode 获取登录二维码
func (h *LoginHandler) GetQRCode(c *gin.Context) {
	var req model.QRCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "invalid request",
		})
		return
	}

	// 生成UUID
	uuid, err := model.NewUUID()
	if err != nil {
		logger.Error("Failed to generate UUID: %v", err)
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "failed to generate UUID",
		})
		return
	}

	// 生成ticket
	ticket, err := model.GenerateTicket()
	if err != nil {
		logger.Error("Failed to generate ticket: %v", err)
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "failed to generate ticket",
		})
		return
	}

	// 生成base64二维码图片 (这里模拟一个简单的二维码)
	qrImage := generateQRCodeBase64(uuid)

	c.JSON(200, model.QRCodeResponse{
		Code:   model.StatusSuccess,
		Key:    uuid,
		UUID:   uuid,
		QRCode: qrImage,
		Ticket: ticket,
		QRURL:  fmt.Sprintf("https://login.weixin.qq.com/qrcode/%s", uuid),
	})
}

// GetQRCodeNewX 获取新版本登录二维码 (绕过验证码)
func (h *LoginHandler) GetQRCodeNewX(c *gin.Context) {
	h.GetQRCode(c)
}

// CheckLoginStatus 检查扫码登录状态
func (h *LoginHandler) CheckLoginStatus(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "missing key",
		})
		return
	}

	// 模拟检查Redis中的登录状态
	// 实际应该从Redis获取状态
	status := "success"

	user, err := h.userRepo.FindByUUID(key)
	if err == nil {
		c.JSON(200, model.CheckLoginStatusResponse{
			Code:     model.StatusLoggedIn,
			Status:   status,
			UUID:     user.UUID,
			WxID:     user.WxID,
			Nickname: user.Nickname,
			Avatar:   user.HeadURL,
		})
		return
	}

	c.JSON(200, model.CheckLoginStatusResponse{
		Code:   model.StatusPending,
		Status: "waiting",
	})
}

// GetLoginStatus 获取在线状态
func (h *LoginHandler) GetLoginStatus(c *gin.Context) {
	key := c.Query("key")

	user, err := h.userRepo.FindByUUID(key)
	if err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "user not found",
		})
		return
	}

	online := user.State == 1
	c.JSON(200, gin.H{
		"code":   model.StatusSuccess,
		"online": online,
		"wxid":   user.WxID,
	})
}

// GetInitStatus 获取初始化状态
func (h *LoginHandler) GetInitStatus(c *gin.Context) {
	c.JSON(200,gin.H{
		"code":   model.StatusSuccess,
		"status": "ready",
	})
}

// CheckCanSetAlias 检测微信登录环境
func (h *LoginHandler) CheckCanSetAlias(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":   model.StatusSuccess,
		"canSet": true,
	})
}

// AutoVerificationCode 自动验证码提交
func (h *LoginHandler) AutoVerificationCode(c *gin.Context) {
	var req model.VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "invalid request",
		})
		return
	}

	logger.Info("Auto verification code: key=%s, code=%s", req.Key, req.Code)

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "verification code submitted",
	})
}

// VerifyCodeAuto 自动处理验证码 (推荐)
func (h *LoginHandler) VerifyCodeAuto(c *gin.Context) {
	h.AutoVerificationCode(c)
}

// VerifyCodeManual 手动处理验证码
func (h *LoginHandler) VerifyCodeManual(c *gin.Context) {
	var req model.VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "invalid request",
		})
		return
	}

	logger.Info("Manual verification code: key=%s, code=%s, ticket=%s", req.Key, req.Code, req.Ticket)

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "verification code submitted",
	})
}

// DeviceLogin 设备登录
func (h *LoginHandler) DeviceLogin(c *gin.Context) {
	var req struct {
		Key      string `json:"key" binding:"required"`
		DeviceID string `json:"deviceId"`
		Data     string `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "invalid request",
		})
		return
	}

	logger.Info("Device login: key=%s", req.Key)

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "device login success",
	})
}

// A16Login A16数据登录
func (h *LoginHandler) A16Login(c *gin.Context) {
	var req struct {
		Key  string `json:"key" binding:"required"`
		A16  string `json:"a16" binding:"required"`
		UUID string `json:"uuid"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "invalid request",
		})
		return
	}

	logger.Info("A16 login: key=%s", req.Key)

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "A16 login success",
	})
}

// SmsLogin 短信登录
func (h *LoginHandler) SmsLogin(c *gin.Context) {
	var req struct {
		Key      string `json:"key" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Code     string `json:"code" binding:"required"`
		DeviceID string `json:"deviceId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "invalid request",
		})
		return
	}

	logger.Info("SMS login: key=%s, phone=%s", req.Key, req.Phone)

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "SMS login success",
	})
}

// Logout 登出
func (h *LoginHandler) Logout(c *gin.Context) {
	key := c.Query("key")

	if err := h.userRepo.UpdateState(key, 0); err != nil {
		logger.Error("Failed to logout: %v", err)
		c.JSON(200, model.Response{
			Code:    model.StatusError,
			Message: "logout failed",
		})
		return
	}

	c.JSON(200, model.Response{
		Code:    model.StatusSuccess,
		Message: "logout success",
	})
}

// generateQRCodeBase64 生成二维码图片的base64编码
// 这里只是一个模拟实现，实际应使用真实的二维码生成库
func generateQRCodeBase64(uuid string) string {
	// SVG格式的简单二维码模拟
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200">
		<rect width="200" height="200" fill="white"/>
		<text x="100" y="100" font-size="12" text-anchor="middle">WeChat: %s</text>
	</svg>`, uuid)

	return fmt.Sprintf("data:image/svg+xml;base64,%s", base64.StdEncoding.EncodeToString([]byte(svg)))
}

// getCurrentTimestamp 获取当前时间戳字符串
func getCurrentTimestamp() string {
	return time.Now().Format("20060102150405")
}

// getCurrentTimestampInt64 获取当前时间戳(int64)
func getCurrentTimestampInt64() int64 {
	return time.Now().Unix()
}
