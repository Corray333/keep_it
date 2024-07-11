package storage

import (
	"database/sql"
	"encoding/json"
	"os"
	"time"

	"github.com/Corray333/authbot/internal/types"
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

type user struct {
	Username string `db:"username"`
}

func (s *Storage) UsernameIsAppropriate(username, tg_username string) (bool, error) {
	found := ""
	row := s.DB.QueryRow("SELECT username FROM users WHERE tg_username = $1", tg_username)
	if err := row.Scan(&found); err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	return found == username, nil
}
