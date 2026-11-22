package admin

import (
	"context"

	"github.com/fremenkiel/phony-g7/internal/event"
)

type eventService interface {
	ListEvents(context.Context) ([]event.Event, error)
}

type Service struct {
	eventService eventService
}

func NewService(es eventService) *Service {
	return &Service{
		eventService: es,
	}
}

func (s *Service) GetDashboardData(ctx context.Context) (DashboardData, error) {
	events, err := s.eventService.ListEvents(ctx)
	if err != nil {
		return DashboardData{}, err
	}

	return DashboardData{
		events: events,
	}, nil
}
