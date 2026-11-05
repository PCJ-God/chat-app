package handler

import (
	"net/http"
	"strconv"

	"chat-app/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := h.notificationService.GetNotifications(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items":     list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "通知ID无效"})
		return
	}

	if err := h.notificationService.MarkAsRead(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "标记已读失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "标记已读成功"})
}
