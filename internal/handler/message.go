package handler

import (
	"net/http"
	"strconv"

	"chat-app/internal/model"
	"chat-app/internal/service"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService service.MessageService
}

func NewMessageHandler(messageService service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效",
		})
		return
	}

	if req.Type == "" {
		req.Type = "text"
	}

	message := &model.Message{
		SenderID:   userID,
		ReceiverID: req.ReceiverID,
		RoomID:     req.RoomID,
		Content:    req.Content,
		Type:       req.Type,
		Status:     "sent",
	}

	if err := h.messageService.SendMessage(message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "发送消息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "消息发送成功",
		"data":    message,
	})
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var targetID *uint
	var roomID *uint

	if targetIDStr := c.Query("target_id"); targetIDStr != "" {
		if id, err := strconv.ParseUint(targetIDStr, 10, 32); err == nil {
			targetID = new(uint)
			*targetID = uint(id)
		}
	}

	if roomIDStr := c.Query("room_id"); roomIDStr != "" {
		if id, err := strconv.ParseUint(roomIDStr, 10, 32); err == nil {
			roomID = new(uint)
			*roomID = uint(id)
		}
	}

	messages, total, err := h.messageService.GetMessages(userID, targetID, roomID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取消息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items":     messages,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	messageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "消息ID无效"})
		return
	}

	if err := h.messageService.MarkAsRead(uint(messageID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "标记已读失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "标记已读成功"})
}

func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")
	count, err := h.messageService.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取未读数失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"count": count}})
}
