package storage

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Corray333/authbot/internal/types"
	"github.com/go-redis/redis"
)

type Storage struct {
	Redis *redis.Client
}

const (
	PHONE_LIFETIME = 60 * 15
)

func NewStorage() *Storage {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if res := redisClient.Ping(); res.Err() != nil {
		panic(res.Err())
	}

	return &Storage{
		Redis: redisClient,
	}
}

func (s *Storage) SetUserRequest(query *types.CodeQuery) error {
	serialized, err := json.Marshal(query)
	if err != nil {
		return err
	}

	if res := s.Redis.Set(query.Username, string(serialized), PHONE_LIFETIME*time.Second); res.Err() != nil {
		return res.Err()
	}
	return nil
}
