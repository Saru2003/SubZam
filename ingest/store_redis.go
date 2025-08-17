package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"os"
)

var redisClient *redis.Client
var ctx = context.Background()

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

var redisCfg RedisConfig

func init() {
	data, err := os.ReadFile("../config/config.yaml")
	if err != nil {
		log.Fatal("Could not read config:", err)
	}

	tempCfg := struct {
		Redis RedisConfig `yaml:"redis"`
	}{}

	if err := yaml.Unmarshal(data, &tempCfg); err != nil {
		log.Fatal("Could not parse config:", err)
	}

	redisCfg = tempCfg.Redis

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Could not connect to Redis:", err)
	}

	log.Println("Connected to Redis at", redisCfg.Addr)
}


func StoreRedis(closestHash uint64, phoneticHash string, chunk Chunk, title, year string) {
    // Store SimHash
    closestKey := fmt.Sprintf("closest:%d", closestHash)
    redisClient.HSet(ctx, closestKey, map[string]interface{}{
        "title":    title,
        "year":     year,
        "chunk":    chunk.Cleaned,
        "raw":      chunk.Original,
        "phonetic": phoneticHash,
    })

    // Store phonetic 
    phoneticKey := fmt.Sprintf("phonetic:%s", phoneticHash)
    redisClient.HSet(ctx, phoneticKey, map[string]interface{}{
        "title": title,
        "year":  year,
        "chunk": chunk.Cleaned,
        "raw":   chunk.Original,
    })
}
