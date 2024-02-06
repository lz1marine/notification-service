package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lz1marine/notification-service/pkg/handler"
	ih "github.com/lz1marine/notification-service/pkg/handler/inter"
)

// NotificationServiceHandlers adds all notification service routes to the gin engine. They are exposed externally
func NotificationServiceHandlers(c *gin.Engine, eh *handler.ExternalHandler) {
	c.GET("/api/v1/notifications", eh.GetChannels)
	c.GET("/api/v1/notifications/sub/:id", eh.GetSubNotifications)
	c.PATCH("/api/v1/notifications/sub/:id", eh.PatchSubNotifications)
}

// NotificationHandlers adds all notification service routes to the gin engine. They are exposed internally
func NotificationHandlers(c *gin.Engine, nh *ih.InternalHandler) {
	c.POST("/api/v1/internal/notifications/:id", nh.PostNotification)
}
