package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type TxKey struct{}

// Errors
var (
	ErrInvalidTxType = fmt.Errorf("invalid transaction type")
)

type Storage struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

type Transactioner interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	GetTx(ctx context.Context) (tx *sqlx.Tx, isNew bool, err error)
}

func New() (*Storage, error) {
	dbName := os.Getenv("POSTGRES_DB_NAME")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbHost := os.Getenv("POSTGRES_HOST")

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
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
	}, nil
}

func (r *Storage) Begin(ctx context.Context) (context.Context, error) {
	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, TxKey{}, tx), nil
}

func (r *Storage) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx.Commit()
}

func (r *Storage) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx.Rollback()
}

func (r *Storage) GetTx(ctx context.Context) (tx *sqlx.Tx, isNew bool, err error) {
	txRaw := ctx.Value(TxKey{})
	if txRaw != nil {
		var ok bool
		tx, ok = txRaw.(*sqlx.Tx)
		if !ok {
			slog.Error("invalid transaction type")
			return nil, false, ErrInvalidTxType
		}
	}
	if tx == nil {
		tx, err = r.DB.BeginTxx(ctx, nil)
		if err != nil {
			slog.Error("failed to begin transaction: " + err.Error())
			return nil, false, err
		}

		return tx, true, nil
	}

	return tx, false, nil
}
