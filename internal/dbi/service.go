package dbi

import (
	"context"

	"github.com/DevHeaven/db/domain/models"
)

type Service interface {
	// Get a single document from the database
	FindInPep1(ctx context.Context, linkedInID string) (models.Payload, error)
}

type Storage interface {
	// Get a single document from the database
	FindInPep1(ctx context.Context, linkedInID string) (models.Payload, error)
}

type service struct {
	repo Storage
}

func NewService(repo Storage) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) FindInPep1(ctx context.Context, linkedInID string) (models.Payload, error) {
	resp, err := s.repo.FindInPep1(ctx, linkedInID)
	return resp, err
}
