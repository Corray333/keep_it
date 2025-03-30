package service

import (
	"context"
	"log/slog"

	"github.com/Corray333/keep_it/internal/helpers"
)

type vkIdSetter interface {
	SetVKID(ctx context.Context, userID int64, vkID int64) error
}

func (s *UserService) SetVKID(ctx context.Context, userID int64, vkIDToken string) error {
	vkID, err := helpers.ExtractVKIDCredentials(vkIDToken)
	if err != nil {
		slog.Error("Error extracting VK ID", "error", err)
		return err
	}
	return s.vkIdSetter.SetVKID(ctx, userID, vkID)
}
