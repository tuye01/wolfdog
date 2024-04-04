package cmd

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"wolfdog/internal/consts"
)

func InitDBEnv() {
	// 连接数据库
	dsn := "root:112411@tcp(localhost:3306)/wolfdog?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	consts.GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

func InitRedisEnv() {
	// 创建 Redis 客户端
	consts.RedisDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "",               // Redis 密码
		DB:       0,                // 使用默认数据库
	})

	// 检查连接是否成功
	pong, err := consts.RedisDB.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)
}
