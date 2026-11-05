package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"chat-app/internal/repository"
	"chat-app/pkg/websocket"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

const onlineUsersKey = "online:users"

type PresenceService interface {
	SetOnline(userID uint) error
	SetOffline(userID uint) error
	SetAway(userID uint) error
	GetStatus(userID uint) (string, error)
	GetOnlineUsers() ([]uint, error)
	UpdateLastSeen(userID uint) error
}

type PresenceServiceImpl struct {
	userRepo repository.UserRepository
	redis    *redis.Client
	hub      *websocket.Hub
}

func NewPresenceService(userRepo repository.UserRepository, rdb *redis.Client, hub *websocket.Hub) PresenceService {
	return &PresenceServiceImpl{
		userRepo: userRepo,
		redis:    rdb,
		hub:      hub,
	}
}

func (s *PresenceServiceImpl) setStatus(userID uint, status string, online bool) error {
	// 更新数据库状态
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.Status = status
	now := time.Now()
	user.LastSeen = &now

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// 更新Redis缓存
	key := fmt.Sprintf("user:status:%d", userID)
	if online {
		s.redis.Set(ctx, key, status, 5*time.Minute)
		s.redis.SAdd(ctx, onlineUsersKey, userID)
	} else {
		s.redis.Del(ctx, key)
		s.redis.SRem(ctx, onlineUsersKey, userID)
	}

	// 广播状态变更
	s.broadcastStatusChange(userID, status)

	return nil
}

func (s *PresenceServiceImpl) SetOnline(userID uint) error {
	return s.setStatus(userID, "online", true)
}

func (s *PresenceServiceImpl) SetOffline(userID uint) error {
	return s.setStatus(userID, "offline", false)
}

func (s *PresenceServiceImpl) SetAway(userID uint) error {
	return s.setStatus(userID, "away", true)
}

func (s *PresenceServiceImpl) GetStatus(userID uint) (string, error) {
	key := fmt.Sprintf("user:status:%d", userID)
	status, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return "offline", nil
	}
	if err != nil {
		return "", err
	}
	return status, nil
}

func (s *PresenceServiceImpl) GetOnlineUsers() ([]uint, error) {
	members, err := s.redis.SMembers(ctx, onlineUsersKey).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]uint, 0, len(members))
	for _, m := range members {
		if id, err := strconv.ParseUint(m, 10, 64); err == nil {
			ids = append(ids, uint(id))
		}
	}
	return ids, nil
}

func (s *PresenceServiceImpl) UpdateLastSeen(userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	now := time.Now()
	user.LastSeen = &now
	return s.userRepo.Update(user)
}

func (s *PresenceServiceImpl) broadcastStatusChange(userID uint, status string) {
	s.hub.BroadcastAll(map[string]interface{}{
		"type":    "presence",
		"user_id": userID,
		"status":  status,
	})
}
