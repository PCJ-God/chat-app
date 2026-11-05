package handler

import (
	"net/http"

	"chat-app/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// Me 返回当前登录用户
func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

// List 返回用户列表（供开始会话时选择）
func (h *UserHandler) List(c *gin.Context) {
	users, _, err := h.userRepo.List(1, 200)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取用户列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}
