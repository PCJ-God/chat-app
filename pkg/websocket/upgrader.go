package websocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 生产环境应该验证Origin
		return true
	},
	HandshakeTimeout: 10 * time.Second,
}

// 配置WebSocket连接参数
func ConfigureUpgrader(readBuffer, writeBuffer int, checkOrigin func(*http.Request) bool) {
	Upgrader.ReadBufferSize = readBuffer
	Upgrader.WriteBufferSize = writeBuffer
	if checkOrigin != nil {
		Upgrader.CheckOrigin = checkOrigin
	}
}
