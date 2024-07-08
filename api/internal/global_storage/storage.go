package global_storage

import (
	"os"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type storage struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

type Storage interface {
	GetDB() *sqlx.DB
	GetRedis() *redis.Client
}

func New() *storage {
	db, err := sqlx.Open("postgres", os.Getenv("DB_CONN_STR"))
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

	return &storage{
		DB:    db,
		Redis: redisClient,
	}
}

func (s *storage) GetDB() *sqlx.DB {
	return s.DB
}

func (s *storage) GetRedis() *redis.Client {
	return s.Redis
}
