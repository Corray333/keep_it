package repository

import (
	"context"
	"encoding/json"

	"github.com/Corray333/keep_it/internal/domains/user/entities"
)

func (s *UserRepository) GetCodeRequest(ctx context.Context, username string) (*entities.CodeQuery, error) {
	res := s.Redis.Get(username)
	if err := res.Err(); err != nil {
		return nil, err
	}

	query := &entities.CodeQuery{}

	if err := json.Unmarshal([]byte(res.Val()), query); err != nil {
		return nil, err
	}

	return query, nil
}
