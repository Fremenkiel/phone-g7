package action

import "context"

type repository interface {
	CreateOutgoingEvent(ctx context.Context, payload string) (int64, error)
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

func (s *Service) CreateEvent(ctx context.Context, payload string) (int64, error) {
	return s.repository.CreateOutgoingEvent(ctx, payload)
}
