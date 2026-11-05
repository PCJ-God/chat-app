package websocket

import (
	"encoding/json"
	"log"
	"time"
)

// RegisterClient 注册客户端
func (h *Hub) RegisterClient(c *Client) {
	h.register <- c
}

// UnregisterClient 注销客户端
func (h *Hub) UnregisterClient(c *Client) {
	h.unregister <- c
}

// JoinRoom 客户端加入房间
func (h *Hub) JoinRoom(roomID uint, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][c] = true
}

// LeaveRoom 客户端离开房间
func (h *Hub) LeaveRoom(roomID uint, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if members, ok := h.rooms[roomID]; ok {
		delete(members, c)
		if len(members) == 0 {
			delete(h.rooms, roomID)
		}
	}
}

// RouteMessage 根据消息目标路由到用户或房间
func (h *Hub) RouteMessage(sender *Client, msg *Message) {
	msg.From = sender.UserID
	msg.Timestamp = time.Now().Unix()

	if msg.RoomID != nil {
		h.BroadcastToRoom(*msg.RoomID, msg)
		return
	}
	if msg.To != nil {
		h.SendToUser(*msg.To, msg)
		return
	}
	h.BroadcastAll(msg)
}

// SendToUser 发送消息给指定用户的所有连接
func (h *Hub) SendToUser(userID uint, data interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("序列化消息失败: %v", err)
		return
	}

	h.mu.RLock()
	var targets []*Client
	for c := range h.clients {
		if c.UserID == userID {
			targets = append(targets, c)
		}
	}
	h.mu.RUnlock()

	for _, c := range targets {
		h.sendToClient(c, payload)
	}
}

// BroadcastToRoom 向房间内所有客户端广播
func (h *Hub) BroadcastToRoom(roomID uint, data interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("序列化消息失败: %v", err)
		return
	}

	h.mu.RLock()
	var targets []*Client
	for c := range h.rooms[roomID] {
		targets = append(targets, c)
	}
	h.mu.RUnlock()

	for _, c := range targets {
		h.sendToClient(c, payload)
	}
}

// BroadcastAll 向所有客户端广播
func (h *Hub) BroadcastAll(data interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("序列化消息失败: %v", err)
		return
	}

	h.mu.RLock()
	var targets []*Client
	for c := range h.clients {
		targets = append(targets, c)
	}
	h.mu.RUnlock()

	for _, c := range targets {
		h.sendToClient(c, payload)
	}
}

func (h *Hub) sendToClient(c *Client, payload []byte) {
	select {
	case c.Send <- payload:
	default:
		log.Printf("客户端 %s 发送队列已满，消息丢弃", c.ID)
	}
}
