package event

import (
	"context"
)

type repository interface {
	Querier
}

type Service struct {
	repository repository
}

var _ service = (*Service)(nil)

func NewService(r repository) *Service {
	return &Service{
		repository: r,
	}
}

func (s *Service) GetEvents(ctx context.Context) ([]Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	events, err := s.repository.ListEvents(ctx)
	if err != nil {
		return nil, err
	}

	return events, nil
}
