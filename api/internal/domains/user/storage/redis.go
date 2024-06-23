package storage

import (
	"encoding/json"

	"github.com/Corray333/keep_it/internal/domains/user/types"
)

func (s *UserStorage) GetCodeRequest(username string) (*types.CodeQuery, error) {
	res := s.redis.Get(username)
	if err := res.Err(); err != nil {
		return nil, err
	}

	query := &types.CodeQuery{}

	if err := json.Unmarshal([]byte(res.Val()), query); err != nil {
		return nil, err
	}

	return query, nil
}
