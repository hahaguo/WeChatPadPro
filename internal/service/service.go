package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/WeChatPadPro/WeChatPadPro/internal/config"
	"github.com/WeChatPadPro/WeChatPadPro/internal/model"
	"github.com/WeChatPadPro/WeChatPadPro/internal/repository"
	"github.com/WeChatPadPro/WeChatPadPro/pkg/logger"
	"golang.org/x/crypto/hmac"
	"golang.org/x/crypto/sha256"
)

// AuthService 认证服务
type AuthService struct {
	adminKey     string
	licenseRepo  *repository.LicenseRepository
	deviceRepo   *repository.DeviceRepository
	licenseCache sync.Map
}

// NewAuthService 创建认证服务
func NewAuthService(adminKey string, licenseRepo *repository.LicenseRepository, deviceRepo *repository.DeviceRepository) *AuthService {
	return &AuthService{
		adminKey:    adminKey,
		licenseRepo: licenseRepo,
		deviceRepo:  deviceRepo,
	}
}

// VerifyAdminKey 验证管理员密钥
func (s *AuthService) VerifyAdminKey(key string) bool {
	return key == s.adminKey
}

// GenerateLicenseKeys 生成授权密钥
func (s *AuthService) GenerateLicenseKeys(count, days int) ([]string, error) {
	var keys []string
	expiryDate := time.Now().AddDate(0, 0, days).Format("2006-01-02 15:04:05")

	for i := 0; i < count; i++ {
		uuid, err := model.NewUUID()
		if err != nil {
			return nil, err
		}

		ticket, err := model.GenerateTicket()
		if err != nil {
			return nil, err
		}

		license := fmt.Sprintf("expaL9%s", ticket)

		lk := &model.LicenseKey{
			Status:     1,
			License:    license,
			ExpiryDate: expiryDate,
		}

		if err := s.licenseRepo.Create(lk); err != nil {
			return nil, err
		}

		keys = append(keys, license)
	}

	return keys, nil
}

// VerifyLicense 验证授权
func (s *AuthService) VerifyLicense(license string) (*model.LicenseKey, error) {
	// 先查缓存
	if cached, ok := s.licenseCache.Load(license); ok {
		if lk, ok := cached.(*model.LicenseKey); ok {
			return lk, nil
		}
	}

	// 查数据库
	lk, err := s.licenseRepo.FindByLicense(license)
	if err != nil {
		return nil, fmt.Errorf("invalid license: %w", err)
	}

	// 检查状态和过期时间
	if lk.Status != 1 {
		return nil, fmt.Errorf("license is disabled")
	}

	expiryTime, err := time.Parse("2006-01-02 15:04:05", lk.ExpiryDate)
	if err != nil {
		return nil, fmt.Errorf("invalid expiry date")
	}

	if time.Now().After(expiryTime) {
		return nil, fmt.Errorf("license has expired")
	}

	// 缓存
	s.licenseCache.Store(license, lk)

	return lk, nil
}

// BindLicense 绑定授权到微信账号
func (s *AuthService) BindLicense(license, wxID, nickname string) error {
	lk, err := s.VerifyLicense(license)
	if err != nil {
		return err
	}

	if lk.WxID != "" {
		return fmt.Errorf("license already bound to another account")
	}

	lk.WxID = wxID
	lk.Nickname = nickname

	return s.licenseRepo.Update(lk)
}

// MessageService 消息服务
type MessageService struct {
	userRepo *repository.UserRepository
}

// NewMessageService 创建消息服务
func NewMessageService(userRepo *repository.UserRepository) *MessageService {
	return &MessageService{
		userRepo: userRepo,
	}
}

// SendMessage 发送消息
func (s *MessageService) SendMessage(key, toUser, content string) error {
	user, err := s.userRepo.FindByUUID(key)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// 这里应该是实际的微信协议发送逻辑
	logger.Info("Sending message from %s to %s: %s", user.WxID, toUser, content)
	return nil
}

// WebhookService Webhook服务
type WebhookService struct {
	webhookRepo *repository.WebhookRepository
	cfg         config.WebhookConfig
	queue       chan *model.Message
	workers     int
}

// NewWebhookService 创建Webhook服务
func NewWebhookService(webhookRepo *repository.WebhookRepository, cfg config.WebhookConfig) *WebhookService {
	svc := &WebhookService{
		webhookRepo: webhookRepo,
		cfg:         cfg,
		queue:       make(chan *model.Message, cfg.QueueSize),
		workers:     10,
	}

	// 启动工作协程
	for i := 0; i < svc.workers; i++ {
		go svc.worker()
	}

	return svc
}

// PushMessage 推送消息到队列
func (s *WebhookService) PushMessage(msg *model.Message) error {
	select {
	case s.queue <- msg:
		return nil
	default:
		return fmt.Errorf("webhook queue is full")
	}
}

// worker 工作协程
func (s *WebhookService) worker() {
	for msg := range s.queue {
		s.processMessage(msg)
	}
}

// processMessage 处理消息
func (s *WebhookService) processMessage(msg *model.Message) {
	// 获取所有启用的Webhook配置
	configs, err := s.webhookRepo.FindEnabled()
	if err != nil {
		logger.Error("Failed to get webhook configs: %v", err)
		return
	}

	for _, cfg := range configs {
		// 检查消息类型过滤
		if !s.shouldSend(cfg, msg) {
			continue
		}

		// 检查是否过滤自己发送的消息
		if !cfg.IncludeSelfMsg && msg.IsSelfMsg {
			continue
		}

		// 检查微信ID过滤
		if cfg.WxID != "" && cfg.WxID != msg.Key {
			continue
		}

		// 发送Webhook
		go s.sendWebhook(cfg, msg)
	}
}

// shouldSend 判断是否应该发送消息
func (s *WebhookService) shouldSend(cfg *model.WebhookConfig, msg *model.Message) bool {
	if cfg.MessageTypes == "*" {
		return true
	}

	var types []string
	if err := json.Unmarshal([]byte(cfg.MessageTypes), &types); err != nil {
		logger.Error("Failed to parse message types: %v", err)
		return false
	}

	msgTypeStr := fmt.Sprintf("%d", msg.MsgType)
	for _, t := range types {
		if t == "*" || t == msgTypeStr {
			return true
		}
	}

	return false
}

// sendWebhook 发送Webhook请求
func (s *WebhookService) sendWebhook(cfg *model.WebhookConfig, msg *model.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	success := false

	for i := 0; i < cfg.RetryCount; i++ {
		if err := s.doSend(ctx, cfg, msg); err == nil {
			success = true
			break
		}
		// 指数退避
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	// 更新统计
	if err := s.webhookRepo.UpdateStats(cfg.ID, success); err != nil {
		logger.Error("Failed to update webhook stats: %v", err)
	}
}

// doSend 执行发送
func (s *WebhookService) doSend(ctx context.Context, cfg *model.WebhookConfig, msg *model.Message) error {
	// 构建请求体
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 生成签名
	timestamp := time.Now().Unix()
	signature := s.generateSignature(string(body), timestamp, cfg.Secret)

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", cfg.URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Webhook-Signature", signature)

	// 发送请求
	client := &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	logger.Info("Webhook sent successfully to %s", cfg.URL)
	return nil
}

// generateSignature 生成签名
func (s *WebhookService) generateSignature(body string, timestamp int64, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(fmt.Sprintf("%d:%s", timestamp, body)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ListConfigs 列出所有配置
func (s *WebhookService) ListConfigs() ([]*model.WebhookConfig, error) {
	return s.webhookRepo.List()
}

// CreateConfig 创建配置
func (s *WebhookService) CreateConfig(cfg *model.WebhookConfig) error {
	return s.webhookRepo.Create(cfg)
}

// UpdateConfig 更新配置
func (s *WebhookService) UpdateConfig(cfg *model.WebhookConfig) error {
	return s.webhookRepo.Update(cfg)
}

// DeleteConfig 删除配置
func (s *WebhookService) DeleteConfig(id int64) error {
	return s.webhookRepo.Delete(id)
}
