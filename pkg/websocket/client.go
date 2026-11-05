package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client WebSocket客户端
type Client struct {
	ID     string
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
	Done   chan struct{}
	mu     sync.Mutex
}

// Hub 管理所有客户端连接
type Hub struct {
	clients    map[*Client]bool
	rooms      map[uint]map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[uint]map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			for roomID, members := range h.rooms {
				if _, ok := members[client]; ok {
					delete(members, client)
					if len(members) == 0 {
						delete(h.rooms, roomID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// ReadPump 从WebSocket连接读取消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		close(c.Done)
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}

		c.handleMessage(message)
	}
}

// WritePump 向WebSocket连接写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(message []byte) {
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("解析消息失败: %v", err)
		return
	}

	switch msg.Type {
	case "ping":
		c.sendPong()
	case "join":
		if msg.RoomID != nil {
			c.Hub.JoinRoom(*msg.RoomID, c)
		}
	case "leave":
		if msg.RoomID != nil {
			c.Hub.LeaveRoom(*msg.RoomID, c)
		}
	case "message":
		c.Hub.RouteMessage(c, &msg)
	default:
		log.Printf("未知消息类型: %s", msg.Type)
	}
}

func (c *Client) sendPong() {
	c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err := c.Conn.WriteMessage(websocket.PongMessage, nil); err != nil {
		log.Printf("发送pong失败: %v", err)
	}
}
