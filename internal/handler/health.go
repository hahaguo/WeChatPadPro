package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/WeChatPadPro/WeChatPadPro/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(db *gorm.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

// Check 全面健康检查
func (h *HealthHandler) Check(c *gin.Context) {
	status := gin.H{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"version": "1.0.0",
	}

	// 检查数据库
	dbStatus := h.checkDatabase()
	status["database"] = dbStatus

	// 检查Redis
	redisStatus := h.checkRedis()
	status["redis"] = redisStatus

	// 判断整体状态
	if dbStatus["status"] != "ok" || redisStatus["status"] != "ok" {
		status["status"] = "unhealthy"
		c.JSON(503, status)
		return
	}

	c.JSON(200, status)
}

// Ping 简单ping检查
func (h *HealthHandler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// checkDatabase 检查数据库连接
func (h *HealthHandler) checkDatabase() gin.H {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sqlDB, err := h.db.DB()
	if err != nil {
		logger.Error("Failed to get SQL DB: %v", err)
		return gin.H{
			"status": "error",
			"error":  "failed to get database instance",
		}
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Error("Database ping failed: %v", err)
		return gin.H{
			"status": "error",
			"error":  fmt.Sprintf("ping failed: %v", err),
		}
	}

	// 获取数据库统计信息
	stats := sqlDB.Stats()
	return gin.H{
		"status":       "ok",
		"open":         stats.OpenConnections,
		"in_use":       stats.InUse,
		"idle":         stats.Idle,
		"wait_count":   stats.WaitCount,
		"wait_duration": stats.WaitDuration.String(),
		"max_idle_closed": stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}
}

// checkRedis 检查Redis连接
func (h *HealthHandler) checkRedis() gin.H {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.redis.Ping(ctx).Err(); err != nil {
		logger.Error("Redis ping failed: %v", err)
		return gin.H{
			"status": "error",
			"error":  fmt.Sprintf("ping failed: %v", err),
		}
	}

	// 获取Redis统计信息
	info, err := h.redis.Info(ctx, "stats").Result()
	if err != nil {
		logger.Error("Failed to get Redis info: %v", err)
	}

	poolStats := h.redis.PoolStats()

	return gin.H{
		"status":                "ok",
		"hits":                  poolStats.Hits,
		"misses":                poolStats.Misses,
		"timeouts":              poolStats.Timeouts,
		"total_conns":           poolStats.TotalConns,
		"idle_conns":            poolStats.IdleConns,
		"stale_conns":           poolStats.StaleConns,
		"info":                  info,
	}
}
