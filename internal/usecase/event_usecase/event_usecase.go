package event_usecase

import (
	repo "calendar-server/internal/repository/event_repository"
	"context"

	"calendar-server/internal/domain"
	"calendar-server/pkg/logger/zappretty"

	"go.uber.org/zap"
)

// EventUseCaseContract - контракт для работы с событиями
type EventUseCaseContract interface {
	CreateEvent(ctx context.Context, event domain.Event) error
	UpdateEvent(ctx context.Context, event domain.Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetEventsForDay(ctx context.Context, userID, date string) ([]domain.Event, error)
	GetEventsForWeek(ctx context.Context, userID, date string) ([]domain.Event, error)
	GetEventsForMonth(ctx context.Context, userID, date string) ([]domain.Event, error)
}

// EventUseCase - реализация EventUseCaseContract
type EventUseCase struct {
	repo   repo.EventRepository
	logger *zap.Logger
}

// NewEventUseCase - конструктор EventUseCase
func NewEventUseCase(repo repo.EventRepository, logger *zap.Logger) *EventUseCase {
	return &EventUseCase{
		repo:   repo,
		logger: logger,
	}
}

// CreateEvent - метод создания события
func (uc *EventUseCase) CreateEvent(ctx context.Context, event domain.Event) error {
	uc.logger.Debug("Creating event in usecase",
		zappretty.Field("event_id", event.ID),
		zappretty.Field("user_id", event.UserID),
	)

	if err := ctx.Err(); err != nil {
		uc.logger.Warn("Context cancelled before creating event")
		return err
	}

	if err := uc.validateEvent(event); err != nil {
		uc.logger.Warn("Event validation failed",
			zappretty.Field("error", err),
			zappretty.Field("event_id", event.ID),
		)
		return err
	}

	return uc.repo.Create(ctx, event)
}

// UpdateEvent - метод обновления события
func (uc *EventUseCase) UpdateEvent(ctx context.Context, event domain.Event) error {
	uc.logger.Debug("Updating event in usecase",
		zappretty.Field("event_id", event.ID),
	)

	if err := ctx.Err(); err != nil {
		uc.logger.Warn("Context cancelled before updating event")
		return err
	}

	if err := uc.validateEvent(event); err != nil {
		uc.logger.Warn("Event validation failed during update",
			zappretty.Field("error", err),
			zappretty.Field("event_id", event.ID),
		)
		return err
	}

	return uc.repo.Update(ctx, event)
}

// DeleteEvent - метод удаления события
func (uc *EventUseCase) DeleteEvent(ctx context.Context, eventID string) error {
	uc.logger.Debug("Deleting event in usecase",
		zappretty.Field("event_id", eventID),
	)

	if err := ctx.Err(); err != nil {
		uc.logger.Warn("Context cancelled before deleting event")
		return err
	}

	if err := uc.validateEventID(eventID); err != nil {
		uc.logger.Warn("Empty event ID provided for deletion")
		return err
	}

	return uc.repo.Delete(ctx, eventID)
}

// GetEventsForDay - метод получения событий для конкретной даты
func (uc *EventUseCase) GetEventsForDay(ctx context.Context, userID, date string) ([]domain.Event, error) {
	uc.logger.Debug("Getting events for day in usecase",
		zappretty.Field("user_id", userID),
		zappretty.Field("date", date),
	)

	if err := ctx.Err(); err != nil {
		uc.logger.Warn("Context cancelled before getting events")
		return nil, err
	}

	if err := uc.validateUserID(userID); err != nil {
		uc.logger.Warn("Empty user ID provided for events query")
		return nil, err
	}
	if err := uc.validateDate(date); err != nil {
		uc.logger.Warn("Invalid date provided for events query")
		return nil, err
	}

	return uc.repo.GetByUserIDAndDate(ctx, userID, date)
}

// GetEventsForWeek - метод получения событий за неделю
func (uc *EventUseCase) GetEventsForWeek(ctx context.Context, userID, date string) ([]domain.Event, error) {
	uc.logger.Debug("Getting events for week in usecase",
		zappretty.Field("user_id", userID),
		zappretty.Field("date", date),
	)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := uc.validateUserID(userID); err != nil {
		return nil, err
	}
	if err := uc.validateDate(date); err != nil {
		return nil, err
	}

	return uc.repo.GetByUserIDAndWeek(ctx, userID, date)
}

// GetEventsForMonth - метод получения событий за месяц
func (uc *EventUseCase) GetEventsForMonth(ctx context.Context, userID, date string) ([]domain.Event, error) {
	uc.logger.Debug("Getting events for month in usecase",
		zappretty.Field("user_id", userID),
		zappretty.Field("date", date),
	)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := uc.validateUserID(userID); err != nil {
		return nil, err
	}
	if err := uc.validateDate(date); err != nil {
		return nil, err
	}

	return uc.repo.GetByUserIDAndMonth(ctx, userID, date)
}
