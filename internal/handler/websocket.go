package handler

import (
	"fmt"
	"net/http"
	"time"

	"chat-app/internal/service"
	"chat-app/pkg/websocket"

	"github.com/gin-gonic/gin"
)

type WebSocketHandler struct {
	hub      *websocket.Hub
	presence service.PresenceService
}

func NewWebSocketHandler(hub *websocket.Hub, presence service.PresenceService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:      hub,
		presence: presence,
	}
}

func (h *WebSocketHandler) Handle(c *gin.Context) {
	// 获取用户ID（从JWT token中）
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 升级HTTP连接为WebSocket
	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "WebSocket升级失败"})
		return
	}

	// 创建客户端
	client := &websocket.Client{
		ID:     generateClientID(),
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    h.hub,
		Done:   make(chan struct{}),
	}

	// 注册客户端
	h.hub.RegisterClient(client)

	if h.presence != nil {
		_ = h.presence.SetOnline(userID)
	}

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump()

	// 断开后更新在线状态
	go func() {
		<-client.Done
		if h.presence != nil {
			_ = h.presence.SetOffline(userID)
		}
	}()
}

func generateClientID() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
}
