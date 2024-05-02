package dbi

import (
	"context"

	"github.com/DevHeaven/db/domain/models"
)

type Service interface {
	// Get a single document from the database
	FindInPep1(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error)
	GetPersonalEmail(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error)
	GetProfessionalEmails(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error)
}

type Storage interface {
	// Get a single document from the database
	FindInPep1(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error)
	GetPersonalEmail(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error)
	GetProfessionalEmails(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error)
}

type service struct {
	repo Storage
}

func NewService(repo Storage) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) FindInPep1(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error) {
	resp, err := s.repo.FindInPep1(ctx, linkedInID, firstname, lastname, domain)
	return resp, err
}

func (s *service) GetPersonalEmail(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error) {
	resp, err := s.repo.GetPersonalEmail(ctx, linkedInID, firstname, lastname, domain)
	return resp, err
}

func (s *service) GetProfessionalEmails(ctx context.Context, linkedInID string, firstname string, lastname string, domain string) (models.Payload, error) {
	resp, err := s.repo.GetProfessionalEmails(ctx, linkedInID, firstname, lastname, domain)
	return resp, err
}

