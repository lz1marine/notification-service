package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lz1marine/notification-service/pkg/handler"
	ih "github.com/lz1marine/notification-service/pkg/handler/inter"
)

func NotificationServiceHandlers(c *gin.Engine, eh *handler.ExternalHandler) {
	c.GET("/v1/notifications", eh.GetChannels)
	c.GET("/v1/notifications/sub/:id", eh.GetSubNotifications)
	c.PATCH("/v1/notifications/sub/:id", eh.PatchSubNotifications)
}

func NotificationHandlers(c *gin.Engine, nh *ih.InternalHandler) {
	c.POST("/v1/internal/notifications/:id", nh.PostNotification)
}
