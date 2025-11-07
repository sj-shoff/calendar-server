package inmemory

import (
	"calendar-server/internal/domain"
	"context"
	"sort"
	"sync"
	"time"

	"calendar-server/pkg/errors"
	"calendar-server/pkg/logger/zappretty"

	"go.uber.org/zap"
)

// EventRepository - реализация хранилища событий в памяти
type EventRepository struct {
	mu     sync.RWMutex
	events map[string]domain.Event
	logger *zap.Logger
}

// NewEventRepository - конструктор хранилища событий в памяти
func NewEventRepository(logger *zap.Logger) *EventRepository {
	return &EventRepository{
		events: make(map[string]domain.Event),
		logger: logger,
	}
}

// Create - создание события
func (r *EventRepository) Create(ctx context.Context, event domain.Event) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Debug("Creating event in repository",
		zappretty.Field("event_id", event.ID),
		zappretty.Field("user_id", event.UserID),
	)

	if _, exists := r.events[event.ID]; exists {
		r.logger.Warn("Event conflict - ID already exists",
			zappretty.Field("event_id", event.ID),
		)
		return errors.ErrEventConflict
	}

	r.events[event.ID] = event
	r.logger.Debug("Event created successfully in repository",
		zappretty.Field("event_id", event.ID),
	)
	return nil
}

// Update - обновление события
func (r *EventRepository) Update(ctx context.Context, event domain.Event) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Debug("Updating event in repository",
		zappretty.Field("event_id", event.ID),
	)

	if _, exists := r.events[event.ID]; !exists {
		r.logger.Warn("Event not found for update",
			zappretty.Field("event_id", event.ID),
		)
		return errors.ErrEventNotFound
	}

	r.events[event.ID] = event
	r.logger.Debug("Event updated successfully in repository",
		zappretty.Field("event_id", event.ID),
	)
	return nil
}

// Delete - удаление события
func (r *EventRepository) Delete(ctx context.Context, eventID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Debug("Deleting event in repository",
		zappretty.Field("event_id", eventID),
	)

	if _, exists := r.events[eventID]; !exists {
		r.logger.Warn("Event not found for deletion",
			zappretty.Field("event_id", eventID),
		)
		return errors.ErrEventNotFound
	}

	delete(r.events, eventID)
	r.logger.Debug("Event deleted successfully from repository",
		zappretty.Field("event_id", eventID),
	)
	return nil
}

// GetByUserIDAndDate - получение событий по ID пользователя и дате
func (r *EventRepository) GetByUserIDAndDate(ctx context.Context, userID, date string) ([]domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var events []domain.Event
	for _, event := range r.events {
		if event.UserID == userID && event.Date == date {
			events = append(events, event)
		}
	}

	sortEvents(events)
	return events, nil
}

// GetByUserIDAndWeek - получение событий по ID пользователя и неделе
func (r *EventRepository) GetByUserIDAndWeek(ctx context.Context, userID, date string) ([]domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	targetDate, _ := time.Parse("2006-01-02", date)
	year, week := targetDate.ISOWeek()

	r.mu.RLock()
	defer r.mu.RUnlock()

	var events []domain.Event
	for _, event := range r.events {
		if event.UserID != userID {
			continue
		}

		eventDate, err := time.Parse("2006-01-02", event.Date)
		if err != nil {
			continue
		}

		eventYear, eventWeek := eventDate.ISOWeek()
		if eventYear == year && eventWeek == week {
			events = append(events, event)
		}
	}

	sortEvents(events)
	return events, nil
}

// GetByUserIDAndMonth - получение событий по ID пользователя и месяцу
func (r *EventRepository) GetByUserIDAndMonth(ctx context.Context, userID, date string) ([]domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	targetDate, _ := time.Parse("2006-01-02", date)
	targetYear, targetMonth := targetDate.Year(), targetDate.Month()

	r.mu.RLock()
	defer r.mu.RUnlock()

	var events []domain.Event
	for _, event := range r.events {
		if event.UserID != userID {
			continue
		}

		eventDate, err := time.Parse("2006-01-02", event.Date)
		if err != nil {
			continue
		}

		if eventDate.Year() == targetYear && eventDate.Month() == targetMonth {
			events = append(events, event)
		}
	}

	sortEvents(events)
	return events, nil
}

// sortEvents сортирует события по дате и названию
func sortEvents(events []domain.Event) {
	sort.Slice(events, func(i, j int) bool {
		if events[i].Date == events[j].Date {
			return events[i].Title < events[j].Title
		}
		return events[i].Date < events[j].Date
	})
}
