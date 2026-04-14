package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/WeChatPadPro/WeChatPadPro/internal/model"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// UserRepository 用户仓储
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByUUID 根据 UUID 查找用户
func (r *UserRepository) FindByUUID(uuid string) (*model.UserInfo, error) {
	var user model.UserInfo
	err := r.db.Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByWxID 根据 WxID 查找用户
func (r *UserRepository) FindByWxID(wxID string) (*model.UserInfo, error) {
	var user model.UserInfo
	err := r.db.Where("wx_id = ?", wxID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *model.UserInfo) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.UserInfo) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(wxID string) error {
	return r.db.Where("wx_id = ?", wxID).Delete(&model.UserInfo{}).Error
}

// List 列出所有用户
func (r *UserRepository) List() ([]*model.UserInfo, error) {
	var users []*model.UserInfo
	err := r.db.Find(&users).Error
	return users, err
}

// UpdateState 更新用户状态
func (r *UserRepository) UpdateState(wxID string, state int) error {
	return r.db.Model(&model.UserInfo{}).
		Where("wx_id = ?", wxID).
		Update("state", state).Error
}

// UpdateLastAuthTime 更新最后认证时间
func (r *UserRepository) UpdateLastAuthTime(wxID string) error {
	return r.db.Model(&model.UserInfo{}).
		Where("wx_id = ?", wxID).
		Update("last_auth_time", time.Now()).Error
}

// LicenseRepository 授权仓储
type LicenseRepository struct {
	db *gorm.DB
}

// NewLicenseRepository 创建授权仓储
func NewLicenseRepository(db *gorm.DB) *LicenseRepository {
	return &LicenseRepository{db: db}
}

// FindByLicense 根据授权码查找
func (r *LicenseRepository) FindByLicense(license string) (*model.LicenseKey, error) {
	var lk model.LicenseKey
	err := r.db.Where("license = ?", license).First(&lk).Error
	if err != nil {
		return nil, err
	}
	return &lk, nil
}

// Create 创建授权
func (r *LicenseRepository) Create(lk *model.LicenseKey) error {
	return r.db.Create(lk).Error
}

// Update 更新授权
func (r *LicenseRepository) Update(lk *model.LicenseKey) error {
	return r.db.Save(lk).Error
}

// FindByWxID 根据微信ID查找授权
func (r *LicenseRepository) FindByWxID(wxID string) (*model.LicenseKey, error) {
	var lk model.LicenseKey
	err := r.db.Where("wx_id = ?", wxID).First(&lk).Error
	if err != nil {
		return nil, err
	}
	return &lk, nil
}

// DeviceRepository 设备仓储
type DeviceRepository struct {
	db *gorm.DB
}

// NewDeviceRepository 创建设备仓储
func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// FindByWxID 根据微信ID查找设备
func (r *DeviceRepository) FindByWxID(wxID string) (*model.DeviceInfo, error) {
	var device model.DeviceInfo
	err := r.db.Where("wxid = ?", wxID).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// Create 创建设备信息
func (r *DeviceRepository) Create(device *model.DeviceInfo) error {
	return r.db.Create(device).Error
}

// Update 更新设备信息
func (r *DeviceRepository) Update(device *model.DeviceInfo) error {
	return r.db.Save(device).Error
}

// WebhookRepository Webhook仓储
type WebhookRepository struct {
	db    *gorm.DB
	redis *redis.Client
	mu    sync.RWMutex
}

// NewWebhookRepository 创建Webhook仓储
func NewWebhookRepository(db *gorm.DB, redis *redis.Client) *WebhookRepository {
	return &WebhookRepository{
		db:    db,
		redis: redis,
	}
}

// Create 创建Webhook配置
func (r *WebhookRepository) Create(cfg *model.WebhookConfig) error {
	return r.db.Create(cfg).Error
}

// Update 更新Webhook配置
func (r *WebhookRepository) Update(cfg *model.WebhookConfig) error {
	return r.db.Save(cfg).Error
}

// Delete 删除Webhook配置
func (r *WebhookRepository) Delete(id int64) error {
	return r.db.Delete(&model.WebhookConfig{}, id).Error
}

// FindByID 根据ID查找
func (r *WebhookRepository) FindByID(id int64) (*model.WebhookConfig, error) {
	var cfg model.WebhookConfig
	err := r.db.First(&cfg, id).Error
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// List 列出所有Webhook配置
func (r *WebhookRepository) List() ([]*model.WebhookConfig, error) {
	var configs []*model.WebhookConfig
	err := r.db.Find(&configs).Error
	return configs, err
}

// FindEnabled 查找所有启用的配置
func (r *WebhookRepository) FindEnabled() ([]*model.WebhookConfig, error) {
	var configs []*model.WebhookConfig
	err := r.db.Where("enabled = ?", true).Find(&configs).Error
	return configs, err
}

// UpdateStats 更新统计信息
func (r *WebhookRepository) UpdateStats(id int64, success bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	updates := map[string]interface{}{
		"last_send_time": time.Now().Unix(),
		"last_send_status": success,
	}

	if success {
		updates["total_sent"] = gorm.Expr("total_sent + ?", 1)
	} else {
		updates["total_failed"] = gorm.Expr("total_failed + ?", 1)
	}

	return r.db.Model(&model.WebhookConfig{}).Where("id = ?", id).Updates(updates).Error
}

// GetTicketForKey 从Redis获取ticket
func (r *WebhookRepository) GetTicketForKey(ctx context.Context, key string) (string, error) {
	ticket, err := r.redis.Get(ctx, fmt.Sprintf("ticket:%s", key)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return ticket, err
}

// SetTicketForKey 设置ticket
func (r *WebhookRepository) SetTicketForKey(ctx context.Context, key, ticket string, expiration time.Duration) error {
	return r.redis.Set(ctx, fmt.Sprintf("ticket:%s", key), ticket, expiration).Err()
}

// GetCheckStatusCache 获取检查状态缓存
func (r *WebhookRepository) GetCheckStatusCache(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, fmt.Sprintf("check_status:%s", key)).Result()
}

// SetCheckStatusCache 设置检查状态缓存
func (r *WebhookRepository) SetCheckStatusCache(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.redis.Set(ctx, fmt.Sprintf("check_status:%s", key), value, expiration).Err()
}

// SetLoginStatus 设置登录状态
func (r *WebhookRepository) SetLoginStatus(ctx context.Context, key, status string, expiration time.Duration) error {
	return r.redis.Set(ctx, fmt.Sprintf("login_status:%s", key), status, expiration).Err()
}

// GetLoginStatus 获取登录状态
func (r *WebhookRepository) GetLoginStatus(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, fmt.Sprintf("login_status:%s", key)).Result()
}