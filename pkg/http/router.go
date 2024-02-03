package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lz1marine/notification-service/pkg/handlers"
	ih "github.com/lz1marine/notification-service/pkg/handlers/inter"
	temp "github.com/lz1marine/notification-service/pkg/handlers/temp"
)

func NotificationServiceHandlers(c *gin.Engine) {
	c.GET("/v1/notifications", handlers.GetChannels)
	c.GET("/v1/notifications/sub/:id", handlers.GetSubNotifications)
	c.PATCH("/v1/notifications/sub/:id", handlers.PatchSubNotifications)
}

func NotificationHandlers(c *gin.Engine) {
	c.POST("/v1/internal/notifications/:id", ih.PostNotification)
}

// TODO: remove
func TempNotificationWorkerHandlers(c *gin.Engine, t *temp.ChannelHandler) {
	c.POST("/v1/temp/notifications/:id", t.PostNotification)
}
