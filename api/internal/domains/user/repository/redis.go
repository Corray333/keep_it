package repository

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Corray333/keep_it/internal/domains/user/entities"
	"github.com/Corray333/keep_it/internal/helpers"
	"github.com/go-redis/redis"
)

func (s *UserRepository) GetCodeRequest(ctx context.Context, username string) (*entities.CodeQuery, error) {
	res := s.Redis.Get(username)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return nil, helpers.NewError(http.StatusUnauthorized, "code not found")
		}
		return nil, err
	}

	query := &entities.CodeQuery{}

	if err := json.Unmarshal([]byte(res.Val()), query); err != nil {
		return nil, err
	}

	return query, nil
}
