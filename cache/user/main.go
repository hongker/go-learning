package main

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func init() {
	var err error
	// connect database
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("connect database: %v", err)
	}

	// connect redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})
}
func main() {
	// get user
	user, err := GetUser(1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user)
}
