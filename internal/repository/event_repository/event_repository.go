package event_repository

import (
	"calendar-server/internal/domain"
	"context"
)

// EventRepository определяет контракт для работы с хранилищем событий
type EventRepository interface {
	Create(ctx context.Context, event domain.Event) error
	Update(ctx context.Context, event domain.Event) error
	Delete(ctx context.Context, eventID string) error
	GetByUserIDAndDate(ctx context.Context, userID, date string) ([]domain.Event, error)
	GetByUserIDAndWeek(ctx context.Context, userID, date string) ([]domain.Event, error)
	GetByUserIDAndMonth(ctx context.Context, userID, date string) ([]domain.Event, error)
}
