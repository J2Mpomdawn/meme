package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"meme/service"
)

func RecordUaAndTime(c *gin.Context) {
	logger, err := zap.NewProduction()

	if err != nil {
		service.LogPrint("red", "NewProduction", err)
	}

	oldTime := time.Now()
	ua := c.GetHeader("User-Agent")
	fci := c.GetHeader("Fly-Client-Ip")

	path := service.StrJoin(64, c.GetHeader("Fly-Forwarded-Proto"), "://", c.Request.Host, c.Request.URL.Path)

	c.Set("FullPath", path)

	c.Next()

	logger.Info(
		"incoming request",
		zap.String("path", c.Request.URL.Path),
		zap.String("ra", ua),
		zap.Int("status", c.Writer.Status()),
		zap.Duration("elapsed", time.Since(oldTime)),
		zap.String("fci", fci),
	)
}
