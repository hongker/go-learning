package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

var (
	db          *gorm.DB
	redisClient redis.UniversalClient
	localCache  = cache.New(LocalCacheExpire, LocalCacheCleanup)
)

const (
	LocalCacheExpire  = time.Minute * 10
	LocalCacheCleanup = time.Minute * 30
	RedisCacheExpire  = time.Hour * 24
)

// User 用户信息
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (User) TableName() string {
	return "users"
}

// GetUserFromDB returns a User from the database
func getUserFromDB(id int64) (*User, error) {
	user := &User{}
	err := db.First(user, id).Error
	return user, err
}

// getUserFromRedis returns a User from the redis
func getUserFromRedis(ctx context.Context, id int64) (*User, error) {
	cacheKey := buildCacheKey(id)
	res, err := redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil { // 缓存数据存在
		user := &User{}
		err = json.Unmarshal(res, user)
		return user, err
	}

	if err == redis.Nil { // 缓存数据不存在
		user, err := getUserFromDB(id)
		if err != nil {
			return nil, err
		}

		// 更新redis缓存
		res, _ = json.Marshal(user)
		redisClient.Set(ctx, cacheKey, res, RedisCacheExpire)
		return user, nil
	}

	// 查询错误
	return nil, err

}

// getFromLocalCache returns user from local cache
func getFromLocalCache(id int64) (*User, error) {
	cacheKey := buildCacheKey(id)
	val, found := localCache.Get(cacheKey)
	if found {
		return val.(*User), nil
	}

	user, err := getUserFromRedis(context.Background(), id)
	if err != nil {
		return nil, err
	}

	localCache.SetDefault(cacheKey, user)
	return user, nil
}

// GetUser returns user info
func GetUser(id int64) (*User, error) {
	return getFromLocalCache(id)
}

func buildCacheKey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}
