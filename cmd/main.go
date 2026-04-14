package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WeChatPadPro/WeChatPadPro/internal/config"
	"github.com/WeChatPadPro/WeChatPadPro/internal/database"
	"github.com/WeChatPadPro/WeChatPadPro/internal/handler"
	"github.com/WeChatPadPro/WeChatPadPro/internal/middleware"
	"github.com/WeChatPadPro/WeChatPadPro/internal/model"
	"github.com/WeChatPadPro/WeChatPadPro/internal/repository"
	"github.com/WeChatPadPro/WeChatPadPro/internal/service"
	"github.com/WeChatPadPro/WeChatPadPro/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Debug); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	logger.Info("Starting WeChatPadPro v%s", cfg.Version)

	// 初始化数据库
	db, err := database.InitMySQL(cfg.MySQL)
	if err != nil {
		logger.Error("Failed to init MySQL: %v", err)
		return
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(db); err != nil {
		logger.Error("Failed to migrate database: %v", err)
		return
	}

	// 初始化Redis
	redisClient, err := database.InitRedis(cfg.Redis)
	if err != nil {
		logger.Error("Failed to init Redis: %v", err)
		return
	}

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	licenseRepo := repository.NewLicenseRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	webhookRepo := repository.NewWebhookRepository(db, redisClient)

	// 初始化服务
	authService := service.NewAuthService(cfg.AdminKey, licenseRepo, deviceRepo)
	msgService := service.NewMessageService(userRepo)
	webhookService := service.NewWebhookService(webhookRepo, cfg.Webhook)

	// 初始化处理器
	loginHandler := handler.NewLoginHandler(authService, userRepo, cfg)
	webhookHandler := handler.NewWebhookHandler(webhookService)
	messageHandler := handler.NewMessageHandler(msgService)
	healthHandler := handler.NewHealthHandler(db, redisClient)

	// 设置 Gin 模式
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.Auth(cfg.AdminKey))

	// 注册路由
	registerRoutes(router, loginHandler, webhookHandler, messageHandler, healthHandler)
	registerStaticRoutes(router)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动服务器
	go func() {
		logger.Info("Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

// registerRoutes 注册所有路由
func registerRoutes(r *gin.Engine, login *handler.LoginHandler, webhook *handler.WebhookHandler,
	msg *handler.MessageHandler, health *handler.HealthHandler) {

	// 健康检查
	r.GET("/health", health.Check)
	r.GET("/ping", health.Ping)

	// 登录相关
	loginGroup := r.Group("/api/login")
	{
		loginGroup.GET("/GenAuthKey2", login.GenAuthKey)
		loginGroup.POST("/qr/new", login.GetQRCode)
		loginGroup.POST("/qr/newx", login.GetQRCodeNewX)
		loginGroup.GET("/CheckLoginStatus", login.CheckLoginStatus)
		loginGroup.GET("/GetLoginStatus", login.GetLoginStatus)
		loginGroup.GET("/GetInItStatus", login.GetInitStatus)
		loginGroup.GET("/CheckCanSetAlias", login.CheckCanSetAlias)
		loginGroup.POST("/AutoVerificationcode", login.AutoVerificationCode)
		loginGroup.POST("/verify/auto", login.VerifyCodeAuto)
		loginGroup.POST("/verify/manual", login.VerifyCodeManual)
		loginGroup.POST("/DeviceLogin", login.DeviceLogin)
		loginGroup.POST("/A16Login", login.A16Login)
		loginGroup.POST("/SmsLogin", login.SmsLogin)
		loginGroup.POST("/logout", login.Logout)
	}

	// 消息相关
	msgGroup := r.Group("/api/message")
	{
		msgGroup.POST("/sendText", msg.SendText)
		msgGroup.POST("/sendImage", msg.SendImage)
		msgGroup.POST("/sendFile", msg.SendFile)
		msgGroup.POST("/sendAppMessage", msg.SendAppMessage)
		msgGroup.GET("/getHistory", msg.GetHistory)
		msgGroup.POST("/revoke", msg.RevokeMessage)
	}

	// 好友相关
	friendGroup := r.Group("/api/friend")
	{
		friendGroup.GET("/list", msg.GetFriendList)
		friendGroup.POST("/add", msg.AddFriend)
		friendGroup.POST("/delete", msg.DeleteFriend)
		friendGroup.POST("/updateRemark", msg.UpdateRemark)
	}

	// 群组相关
	groupGroup := r.Group("/api/group")
	{
		groupGroup.GET("/list", msg.GetGroupList)
		groupGroup.POST("/create", msg.CreateGroup)
		groupGroup.POST("/invite", msg.InviteMember)
		groupGroup.POST("/kick", msg.KickMember)
		groupGroup.POST("/quit", msg.QuitGroup)
		groupGroup.POST("/updateName", msg.UpdateGroupName)
		groupGroup.POST("/sendAtMsg", msg.SendAtMessage)
	}

	// Webhook 相关
	webhookGroup := r.Group("/webhook")
	{
		webhookGroup.GET("/List", webhook.List)
		webhookGroup.GET("/Status", webhook.Status)
		webhookGroup.GET("/Test", webhook.Test)
		webhookGroup.POST("/Config", webhook.Config)
		webhookGroup.PUT("/Update", webhook.Update)
		webhookGroup.DELETE("/Delete", webhook.Delete)
		webhookGroup.POST("/Diagnostics", webhook.Diagnostics)
		webhookGroup.POST("/ResetConnection", webhook.ResetConnection)
	}

	// V1 API
	v1Group := r.Group("/v1")
	{
		v1Group.POST("/webhook/Config", webhook.Config)
		v1Group.PUT("/webhook/Update", webhook.Update)
		v1Group.GET("/webhook/List", webhook.List)
		v1Group.GET("/webhook/Status", webhook.Status)
		v1Group.GET("/webhook/Test", webhook.Test)
		v1Group.GET("/webhook/Diagnostics", webhook.Diagnostics)
	}

	// SSE 事件流
	r.GET("/sse", handler.SSEHandler)
}

// registerStaticRoutes 注册静态文件路由
func registerStaticRoutes(r *gin.Engine) {
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "WeChatPadPro 登录",
		})
	})
}
