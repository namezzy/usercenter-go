package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user-center/internal/model"

	"github.com/redis/go-redis/v9"
)

// CacheService Redis缓存服务
type CacheService struct {
	client *redis.Client
}

// NewCacheService 创建缓存服务
func NewCacheService(client *redis.Client) *CacheService {
	return &CacheService{client: client}
}

// SetUser 缓存用户信息
func (s *CacheService) SetUser(ctx context.Context, userId int64, user *model.SafetyUser, expiration time.Duration) error {
	key := fmt.Sprintf("user:%d", userId)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, data, expiration).Err()
}

// GetUser 获取缓存的用户信息
func (s *CacheService) GetUser(ctx context.Context, userId int64) (*model.SafetyUser, error) {
	key := fmt.Sprintf("user:%d", userId)
	data, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	
	var user model.SafetyUser
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser 删除用户缓存
func (s *CacheService) DeleteUser(ctx context.Context, userId int64) error {
	key := fmt.Sprintf("user:%d", userId)
	return s.client.Del(ctx, key).Err()
}

// SetUserList 缓存用户列表
func (s *CacheService) SetUserList(ctx context.Context, cacheKey string, users []model.SafetyUser, expiration time.Duration) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, cacheKey, data, expiration).Err()
}

// GetUserList 获取缓存的用户列表
func (s *CacheService) GetUserList(ctx context.Context, cacheKey string) ([]model.SafetyUser, error) {
	data, err := s.client.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}
	
	var users []model.SafetyUser
	err = json.Unmarshal([]byte(data), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// SetLoginAttempts 设置登录失败次数
func (s *CacheService) SetLoginAttempts(ctx context.Context, account string, attempts int, expiration time.Duration) error {
	key := fmt.Sprintf("login:attempts:%s", account)
	return s.client.Set(ctx, key, attempts, expiration).Err()
}

// GetLoginAttempts 获取登录失败次数
func (s *CacheService) GetLoginAttempts(ctx context.Context, account string) (int, error) {
	key := fmt.Sprintf("login:attempts:%s", account)
	result, err := s.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return result, err
}

// IncrLoginAttempts 增加登录失败次数
func (s *CacheService) IncrLoginAttempts(ctx context.Context, account string, expiration time.Duration) error {
	key := fmt.Sprintf("login:attempts:%s", account)
	pipe := s.client.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expiration)
	_, err := pipe.Exec(ctx)
	return err
}

// DeleteLoginAttempts 删除登录失败次数
func (s *CacheService) DeleteLoginAttempts(ctx context.Context, account string) error {
	key := fmt.Sprintf("login:attempts:%s", account)
	return s.client.Del(ctx, key).Err()
}

// Ping 检查Redis连接
func (s *CacheService) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

// Close 关闭Redis连接
func (s *CacheService) Close() error {
	return s.client.Close()
}
