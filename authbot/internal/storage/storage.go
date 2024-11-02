package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/corray333/keep_it_authbot/internal/types"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	PHONE_LIFETIME = 60 * 15
)

type Storage struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

func New() *Storage {
	dbName := os.Getenv("POSTGRES_DB_NAME")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbHost := os.Getenv("POSTGRES_HOST")

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if res := redisClient.Ping(); res.Err() != nil {
		panic(res.Err())
	}

	return &Storage{
		DB:    db,
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
