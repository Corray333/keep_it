package service

import (
	"context"
)

func (s *UserService) CodeExists(ctx context.Context, username string, syn int64) (bool, error) {
	query, err := s.verificationCodeGetter.GetCodeRequest(ctx, username)
	if err != nil {
		return false, err
	}

	if query == nil {
		return false, nil
	}

	if query.Syn != syn {
		return false, nil
	}

	return true, nil
}
