package handler

// MessageRequest 发送消息请求
type MessageRequest struct {
	ReceiverID *uint  `json:"receiver_id"`
	RoomID     *uint  `json:"room_id"`
	Content    string `json:"content" binding:"required"`
	Type       string `json:"type"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Avatar      string `json:"avatar"`
}

// SendRoomMessageRequest 房间内发送消息请求
type SendRoomMessageRequest struct {
	Content string `json:"content" binding:"required"`
}
