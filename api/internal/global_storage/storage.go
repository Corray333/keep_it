package global_storage

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type storage struct {
	DB *sqlx.DB
}

type Storage interface {
	GetDB() *sqlx.DB
}

func New() *storage {
	db, err := sqlx.Open("postgres", os.Getenv("DB_CONN_STR"))
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &storage{
		DB: db,
	}
}

func (s *storage) GetDB() *sqlx.DB {
	return s.DB
}
