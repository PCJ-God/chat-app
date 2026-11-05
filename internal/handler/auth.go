package handler

import (
	"errors"
	"net/http"
	"time"

	"chat-app/internal/model"
	"chat-app/internal/repository"
	"chat-app/pkg/auth"
	"chat-app/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userRepo  repository.UserRepository
	jwtExpire time.Duration
}

func NewAuthHandler(userRepo repository.UserRepository, jwtExpire time.Duration) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		jwtExpire: jwtExpire,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数无效"})
		return
	}

	if _, err := h.userRepo.GetByUsername(req.Username); err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "用户名已存在"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}

	if _, err := h.userRepo.GetByEmail(req.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "邮箱已被注册"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "密码加密失败"})
		return
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
		Nickname: req.Nickname,
		Status:   "offline",
	}

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "注册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "注册成功", "data": user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数无效"})
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户名或密码错误"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户名或密码错误"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username, h.jwtExpire)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}
