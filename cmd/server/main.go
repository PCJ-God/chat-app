package main

import (
	"fmt"
	"log"

	"chat-app/internal/handler"
	"chat-app/internal/repository"
	"chat-app/internal/service"
	"chat-app/pkg/auth"
	"chat-app/pkg/config"
	"chat-app/pkg/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 初始化数据库
	dsn := repository.BuildDSN(cfg.Database)
	db, err := repository.InitDB(dsn)
	if err != nil {
		log.Fatal("初始化数据库失败:", err)
	}

	// 初始化 Redis
	rdb := repository.NewRedisClient(cfg.Redis)

	// 初始化 JWT
	auth.InitJWT(cfg.JWT.Secret)

	// 配置WebSocket
	websocket.ConfigureUpgrader(
		cfg.WebSocket.ReadBufferSize,
		cfg.WebSocket.WriteBufferSize,
		nil,
	)

	// 依赖装配
	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	roomMemberRepo := repository.NewRoomMemberRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	hub := websocket.NewHub()
	go hub.Run()

	messageService := service.NewMessageService(messageRepo, hub)
	roomService := service.NewRoomService(roomRepo, roomMemberRepo, messageService)
	presenceService := service.NewPresenceService(userRepo, rdb, hub)
	notificationService := service.NewNotificationService(notificationRepo, messageRepo)

	authHandler := handler.NewAuthHandler(userRepo, cfg.JWT.Expire)
	userHandler := handler.NewUserHandler(userRepo)
	messageHandler := handler.NewMessageHandler(messageService)
	roomHandler := handler.NewRoomHandler(roomService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	wsHandler := handler.NewWebSocketHandler(hub, presenceService)

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	r := gin.Default()
	setupRoutes(r, authHandler, userHandler, messageHandler, roomHandler, notificationHandler, wsHandler)

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("聊天服务器启动在 http://%s, WebSocket: ws://%s/ws\n", addr, addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}

func setupRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	messageHandler *handler.MessageHandler,
	roomHandler *handler.RoomHandler,
	notificationHandler *handler.NotificationHandler,
	wsHandler *handler.WebSocketHandler,
) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})

	// 认证（无需登录）
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// 需要登录的接口
	api := r.Group("/api", auth.AuthMiddleware())
	{
		// 用户
		api.GET("/users", userHandler.List)
		api.GET("/users/me", userHandler.Me)

		// 消息
		api.POST("/messages", messageHandler.SendMessage)
		api.GET("/messages", messageHandler.GetMessages)
		api.POST("/messages/:id/read", messageHandler.MarkAsRead)
		api.GET("/messages/unread/count", messageHandler.GetUnreadCount)

		// 房间
		api.POST("/rooms", roomHandler.CreateRoom)
		api.GET("/rooms", roomHandler.GetUserRooms)
		api.GET("/rooms/:id", roomHandler.GetRoom)
		api.POST("/rooms/:id/join", roomHandler.JoinRoom)
		api.POST("/rooms/:id/leave", roomHandler.LeaveRoom)
		api.POST("/rooms/:id/messages", roomHandler.SendRoomMessage)

		// 通知
		api.GET("/notifications", notificationHandler.GetNotifications)
		api.POST("/notifications/:id/read", notificationHandler.MarkAsRead)
	}

	// WebSocket（需要登录）
	r.GET("/ws", auth.AuthMiddleware(), wsHandler.Handle)
}
