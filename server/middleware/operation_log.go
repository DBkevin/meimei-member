package middleware

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OperationLog 操作日志记录中间件
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()
		userID := utils.GetUserID(c)

		// 记录请求开始
		global.GVA_LOG.Debug("API请求开始",
			zap.String("path", path),
			zap.String("method", method),
			zap.String("ip", ip),
			zap.Uint("userId", userID),
		)

		// 继续处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()

		// 记录响应结果
		global.GVA_LOG.Info("API请求完成",
			zap.String("path", path),
			zap.String("method", method),
			zap.Int("statusCode", statusCode),
			zap.Duration("duration", duration),
			zap.String("ip", ip),
			zap.Uint("userId", userID),
		)

		// 如果请求耗时过长，记录警告
		if duration > 5*time.Second {
			global.GVA_LOG.Warn("API请求耗时过长",
				zap.String("path", path),
				zap.String("method", method),
				zap.Duration("duration", duration),
			)
		}
	}
}
