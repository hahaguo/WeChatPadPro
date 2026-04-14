package model

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// UserInfo 用户信息实体
type UserInfo struct {
	TargetIP         string    `json:"targetIp" gorm:"column:targetIp"`
	UUID             string    `json:"uuid" gorm:"column:uuid;primaryKey"`
	Uin              uint32    `json:"uin" gorm:"column:uin"`
	WxID             string    `json:"wxId" gorm:"column:wxId;primaryKey"`
	Nickname         string    `json:"nickname" gorm:"column:nickname"`
	Username         string    `json:"userName" gorm:"column:userName"`
	Password         string    `json:"password" gorm:"column:password"`
	HeadURL          string    `json:"headurl" gorm:"column:headurl"`
	Cookie           []byte    `json:"-" gorm:"column:cookie"`
	SessionKey       []byte    `json:"-" gorm:"column:sessionKey"`
	ShortHost        string    `json:"shorthost" gorm:"column:shorthost"`
	LongHost         string    `json:"longhost" gorm:"column:longhost"`
	ECPUKey          []byte    `json:"-" gorm:"column:ecpukey"`
	ECPRKey          []byte    `json:"-" gorm:"column:ecprkey"`
	ChecksumKey      []byte    `json:"-" gorm:"column:checksumkey"`
	AutoAuthKey      string    `json:"autoauthkey" gorm:"column:autoauthkey"`
	State            int       `json:"state" gorm:"column:state"`
	SyncKey          string    `json:"synckey" gorm:"column:synckey"`
	FavSyncKey       string    `json:"favsynckey" gorm:"column:favsynckey"`
	LoginRSAVer      uint32    `json:"login_rsa_ver" gorm:"column:login_rsa_ver"`
	ErrMsg           string    `json:"err_msg" gorm:"column:err_msg"`
	DeviceCreateTime time.Time `json:"device_create_time" gorm:"column:device_create_time;default:CURRENT_TIMESTAMP"`
	LastLoginTime    time.Time `json:"last_login_time" gorm:"column:last_login_time;default:CURRENT_TIMESTAMP"`
	LastAuthTime     time.Time `json:"last_auth_time" gorm:"column:last_auth_time;default:CURRENT_TIMESTAMP"`
}

// TableName 表名
func (UserInfo) TableName() string {
	return "user_info_entity"
}

// DeviceInfo 设备信息实体
type DeviceInfo struct {
	WxID         string `json:"wxid" gorm:"column:wxid;primaryKey"`
	UUIDOne      string `json:"uuidone" gorm:"column:uuidone"`
	UUIDTwo      string `json:"uuidtwo" gorm:"column:uuidtwo"`
	IMEI         string `json:"imei" gorm:"column:imei"`
	DeviceID     []byte `json:"deviceid" gorm:"column:deviceid"`
	DeviceName   string `json:"devicename" gorm:"column:devicename"`
	Timezone     string `json:"timezone" gorm:"column:timezone"`
	Language     string `json:"language" gorm:"column:language"`
	DeviceBrand  string `json:"devicebrand" gorm:"column:devicebrand"`
	RealCountry  string `json:"realcountry" gorm:"column:realcountry"`
	IPhoneVer    string `json:"iphonever" gorm:"column:iphonever"`
	BoudleID     string `json:"boudleid" gorm:"column:boudleid"`
	OSType       string `json:"ostype" gorm:"column:ostype"`
	AdSource     string `json:"adsource" gorm:"column:adsource"`
	OSTypeNumber string `json:"ostypenumber" gorm:"column:ostypenumber"`
	CoreCount    uint   `json:"corecount" gorm:"column:corecount"`
	CarrierName  string `json:"carriername" gorm:"column:carriername"`
	SoftTypeXML  string `json:"softtypexml" gorm:"column:softtypexml"`
	ClientCheck  string `json:"clientcheckdataxml" gorm:"column:clientcheckdataxml"`
	GUID2        string `json:"guid2" gorm:"column:guid2"`
}

// TableName 表名
func (DeviceInfo) TableName() string {
	return "device_info_entity"
}

// LicenseKey 授权密钥
type LicenseKey struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	DeviceToken string `json:"deviceToken" gorm:"column:device_token"`
	Status      int    `json:"status" gorm:"column:status"`
	License     string `json:"license" gorm:"column:license"`
	ExpiryDate  string `json:"expiryDate" gorm:"column:expiry_date"`
	WxID        string `json:"wx_id" gorm:"column:wx_id"`
	Nickname    string `json:"nick_name" gorm:"column:nick_name"`
}

// TableName 表名
func (LicenseKey) TableName() string {
	return "license_key"
}

// UserBusinessLog 用户业务日志
type UserBusinessLog struct {
	ID           uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID         string `json:"uuid" gorm:"column:uuid"`
	UserName     string `json:"userName" gorm:"column:user_name"`
	BusinessType string `json:"businessType" gorm:"column:business_type"`
	ExResult     string `json:"exResult" gorm:"column:ex_result"`
}

// TableName 表名
func (UserBusinessLog) TableName() string {
	return "user_business_log"
}

// UserLoginLog 用户登录日志
type UserLoginLog struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	TargetIP   string    `json:"targetIp" gorm:"column:targetIp"`
	UUID       string    `json:"u_uid" gorm:"column:u_uid"`
	UserName   string    `json:"userName" gorm:"column:user_name"`
	Nickname   string    `json:"nickName" gorm:"column:nick_name"`
	LoginType  string    `json:"loginType" gorm:"column:login_type"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
	RetCode    int       `json:"retCode" gorm:"column:ret_code"`
	ErrMsg     string    `json:"err_msg" gorm:"column:err_msg"`
}

// TableName 表名
func (UserLoginLog) TableName() string {
	return "user_login_log"
}

// WebhookConfig Webhook 配置
type WebhookConfig struct {
	ID                int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	URL               string    `json:"url" gorm:"column:url;not null;size:500"`
	Secret            string    `json:"secret" gorm:"column:secret;size:255"`
	Enabled           bool      `json:"enabled" gorm:"column:enabled;default:true"`
	Timeout           int       `json:"timeout" gorm:"column:timeout;default:10"`
	RetryCount        int       `json:"retryCount" gorm:"column:retry_count;default:3"`
	MessageTypes      string    `json:"messageTypes" gorm:"column:message_types;type:text"` // JSON array as string
	IncludeSelfMsg    bool      `json:"includeSelfMessage" gorm:"column:include_self_msg;default:true"`
	WxID              string    `json:"wxId" gorm:"column:wxid;size:128"`
	UseDirectStream   bool      `json:"useDirectStream" gorm:"column:use_direct_stream;default:true"`
	UseRedisSync      bool      `json:"useRedisSync" gorm:"column:use_redis_sync;default:false"`
	IndependentMode   bool      `json:"independentMode" gorm:"column:independent_mode;default:true"`
	LastSendTime      int64     `json:"lastSendTime" gorm:"column:last_send_time;default:0"`
	LastSendStatus    bool      `json:"lastSendStatus" gorm:"column:last_send_status;default:false"`
	TotalSent         int64     `json:"totalSent" gorm:"column:total_sent;default:0"`
	TotalFailed       int64     `json:"totalFailed" gorm:"column:total_failed;default:0"`
	CreatedAt         time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 表名
func (WebhookConfig) TableName() string {
	return "webhook_config"
}

// QRCodeRequest 二维码登录请求
type QRCodeRequest struct {
	Proxy     string `json:"proxy" binding:"omitempty"`
	DeviceName string `json:"deviceName" binding:"omitempty"`
	DeviceID   string `json:"deviceId" binding:"omitempty"`
}

// QRCodeResponse 二维码登录响应
type QRCodeResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Key      string `json:"key"`
	UUID     string `json:"uuid"`
	QRCode   string `json:"qrcode"`   // base64 encoded image
	Data62   string `json:"data62"`   // auto generated
	Ticket   string `json:"ticket"`
	QRURL    string `json:"qrcodeUrl"`
}

// CheckLoginStatusRequest 检查登录状态请求
type CheckLoginStatusRequest struct {
	Key string `form:"key" binding:"required"`
}

// CheckLoginStatusResponse 检查登录状态响应
type CheckLoginStatusResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Status   string `json:"status"`
	UUID     string `json:"uuid"`
	WxID     string `json:"wxId"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// VerifyCodeRequest 验证码验证请求
type VerifyCodeRequest struct {
	Key     string `form:"key" binding:"required"`
	Code    string `json:"code" binding:"required"`
	UUID    string `json:"uuid"`
	Ticket  string `json:"ticket"`
	Data62  string `json:"data62"`
}

// GenerateAuthKeyRequest 生成授权密钥请求
type GenerateAuthKeyRequest struct {
	Key   string `form:"key" binding:"required"`
	Count int    `form:"count" binding:"required,min=1,max=100"`
	Days  int    `form:"days" binding:"required,min=1,max=3650"`
}

// GenerateAuthKeyResponse 生成授权密钥响应
type GenerateAuthKeyResponse struct {
	Code int      `json:"code"`
	Data []string `json:"data"`
}

// Message 消息结构
type Message struct {
	Key        string `json:"key"`
	MsgID      string `json:"msgId"`
	Timestamp  int64  `json:"timestamp"`
	FromUser   string `json:"fromUser"`
	ToUser     string `json:"toUser"`
	MsgType    int    `json:"msgType"`
	Content    string `json:"content"`  // JSON string for complex types
	IsSelfMsg  bool   `json:"isSelfMsg"`
	CreateTime int64  `json:"createTime"`
	IsHistory  bool   `json:"isHistory"`
	Seq        int64  `json:"seq"`
}

// SendTextRequest 发送文本消息请求
type SendTextRequest struct {
	Key      string `json:"key" binding:"required"`
	ToUser   string `json:"toUser" binding:"required"`
	Content  string `json:"content" binding:"required"`
	AtUsers  string `json:"atUsers" binding:"omitempty"`
}

// SendFileRequest 发送文件消息请求
type SendFileRequest struct {
	Key      string `form:"key" binding:"required"`
	ToUser   string `form:"toUser" binding:"required"`
	FilePath string `form:"filePath" binding:"required"`
	FileType string `form:"fileType"`
	FileName string `form:"fileName"`
}

// Response 通用响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	// 消息类型常量
	MsgTypeText       = 1    // 文本消息
	MsgTypeImage      = 3    // 图片消息
	MsgTypeVoice      = 34   // 语音消息
	MsgTypeVideo      = 43   // 视频消息
	MsgTypeEmoji      = 47   // 表情消息
	MsgTypeLink       = 49   // 链接消息
	MsgTypeApp        = 50   // APP消息
	MsgTypeVOIP       = 62   // VOIP消息
	MsgTypeSystem     = 10000 // 系统消息
	MsgTypeRecall     = 10002 // 撤回消息
)

// Status 状态常量
const (
	StatusSuccess = 0
	StatusError   = -1
	StatusNeedQR  = -2
	StatusPending = 1
	StatusLoggedIn = 200
)

// NewUUID 生成新的UUID
func NewUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GenerateTicket 生成ticket
func GenerateTicket() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
