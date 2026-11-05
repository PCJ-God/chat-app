package handler

import (
	"net/http"
	"strconv"

	"chat-app/internal/model"
	"chat-app/internal/service"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomService service.RoomService
}

func NewRoomHandler(roomService service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数无效"})
		return
	}

	room := &model.Room{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Avatar:      req.Avatar,
	}

	if err := h.roomService.CreateRoom(room, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建房间失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "房间创建成功", "data": room})
}

func (h *RoomHandler) JoinRoom(c *gin.Context) {
	userID := c.GetUint("user_id")
	roomID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "房间ID无效"})
		return
	}

	if err := h.roomService.JoinRoom(roomID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "加入房间失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "加入房间成功"})
}

func (h *RoomHandler) LeaveRoom(c *gin.Context) {
	userID := c.GetUint("user_id")
	roomID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "房间ID无效"})
		return
	}

	if err := h.roomService.LeaveRoom(roomID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "退出房间失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "退出房间成功"})
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	roomID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "房间ID无效"})
		return
	}

	room, err := h.roomService.GetRoom(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取房间失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": room})
}

func (h *RoomHandler) GetUserRooms(c *gin.Context) {
	userID := c.GetUint("user_id")
	rooms, err := h.roomService.GetUserRooms(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取房间列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": rooms})
}

func (h *RoomHandler) SendRoomMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	roomID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "房间ID无效"})
		return
	}

	var req SendRoomMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数无效"})
		return
	}

	if err := h.roomService.SendRoomMessage(roomID, userID, req.Content); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "消息发送成功"})
}

func parseID(s string) (uint, error) {
	id, err := strconv.ParseUint(s, 10, 32)
	return uint(id), err
}
