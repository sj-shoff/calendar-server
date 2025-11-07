package event_usecase

import (
	"calendar-server/internal/domain"
	"calendar-server/pkg/errors"
	"context"
	stdErrors "errors"
	"testing"
	"time"

	"go.uber.org/zap"
)

type mockEventRepository struct {
	events map[string]domain.Event
}

func newMockEventRepository() *mockEventRepository {
	return &mockEventRepository{
		events: make(map[string]domain.Event),
	}
}

func (m *mockEventRepository) Create(ctx context.Context, event domain.Event) error {
	if _, exists := m.events[event.ID]; exists {
		return errors.ErrEventConflict
	}
	m.events[event.ID] = event
	return nil
}

func (m *mockEventRepository) Update(ctx context.Context, event domain.Event) error {
	if _, exists := m.events[event.ID]; !exists {
		return errors.ErrEventNotFound
	}
	m.events[event.ID] = event
	return nil
}

func (m *mockEventRepository) Delete(ctx context.Context, eventID string) error {
	if _, exists := m.events[eventID]; !exists {
		return errors.ErrEventNotFound
	}
	delete(m.events, eventID)
	return nil
}

func (m *mockEventRepository) GetByUserIDAndDate(ctx context.Context, userID, date string) ([]domain.Event, error) {
	var result []domain.Event
	for _, event := range m.events {
		if event.UserID == userID && event.Date == date {
			result = append(result, event)
		}
	}
	return result, nil
}

func (m *mockEventRepository) GetByUserIDAndWeek(ctx context.Context, userID, date string) ([]domain.Event, error) {
	targetDate, _ := time.Parse("2006-01-02", date)
	year, week := targetDate.ISOWeek()

	var result []domain.Event
	for _, event := range m.events {
		if event.UserID != userID {
			continue
		}

		eventDate, err := time.Parse("2006-01-02", event.Date)
		if err != nil {
			continue
		}

		eventYear, eventWeek := eventDate.ISOWeek()
		if eventYear == year && eventWeek == week {
			result = append(result, event)
		}
	}
	return result, nil
}

func (m *mockEventRepository) GetByUserIDAndMonth(ctx context.Context, userID, date string) ([]domain.Event, error) {
	targetDate, _ := time.Parse("2006-01-02", date)
	targetYear, targetMonth := targetDate.Year(), targetDate.Month()

	var result []domain.Event
	for _, event := range m.events {
		if event.UserID != userID {
			continue
		}

		eventDate, err := time.Parse("2006-01-02", event.Date)
		if err != nil {
			continue
		}

		if eventDate.Year() == targetYear && eventDate.Month() == targetMonth {
			result = append(result, event)
		}
	}
	return result, nil
}

func setupTestUseCase() (*EventUseCase, context.Context) {
	logger, _ := zap.NewDevelopment()
	repo := newMockEventRepository()
	uc := NewEventUseCase(repo, logger)
	ctx := context.Background()
	return uc, ctx
}

func TestEventUseCase_CreateEvent(t *testing.T) {
	uc, ctx := setupTestUseCase()

	validEvent := domain.Event{
		ID:     "test-1",
		UserID: "user-1",
		Date:   "2025-01-15",
		Title:  "Valid Event",
	}

	err := uc.CreateEvent(ctx, validEvent)
	if err != nil {
		t.Fatalf("Failed to create valid event: %v", err)
	}

	err = uc.CreateEvent(ctx, validEvent)
	if !stdErrors.Is(err, errors.ErrEventConflict) {
		t.Errorf("Expected ErrEventConflict for duplicate event, got %v", err)
	}
}

func TestEventUseCase_Validation(t *testing.T) {
	uc, ctx := setupTestUseCase()

	testCases := []struct {
		name      string
		event     domain.Event
		expectErr error
	}{
		{
			name:      "empty event ID",
			event:     domain.Event{UserID: "user-1", Date: "2025-01-15", Title: "Title"},
			expectErr: errors.ErrEmptyEventID,
		},
		{
			name:      "empty user ID",
			event:     domain.Event{ID: "test-1", Date: "2025-01-15", Title: "Title"},
			expectErr: errors.ErrEmptyUserID,
		},
		{
			name:      "empty title",
			event:     domain.Event{ID: "test-1", UserID: "user-1", Date: "2025-01-15", Title: ""},
			expectErr: errors.ErrEmptyTitle,
		},
		{
			name:      "invalid date format",
			event:     domain.Event{ID: "test-1", UserID: "user-1", Date: "2025/01/15", Title: "Title"},
			expectErr: errors.ErrInvalidDate,
		},
		{
			name:      "valid event",
			event:     domain.Event{ID: "test-1", UserID: "user-1", Date: "2025-01-15", Title: "Title"},
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := uc.CreateEvent(ctx, tc.event)
			if !stdErrors.Is(err, tc.expectErr) {
				t.Errorf("Expected error %v, got %v", tc.expectErr, err)
			}
		})
	}
}

func TestEventUseCase_UpdateEvent(t *testing.T) {
	uc, ctx := setupTestUseCase()

	event := domain.Event{
		ID:     "test-1",
		UserID: "user-1",
		Date:   "2025-01-15",
		Title:  "Original Title",
	}

	err := uc.CreateEvent(ctx, event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	event.Title = "Updated Title"
	err = uc.UpdateEvent(ctx, event)
	if err != nil {
		t.Fatalf("Failed to update event: %v", err)
	}

	nonExistentEvent := domain.Event{
		ID:     "non-existent",
		UserID: "user-1",
		Date:   "2025-01-15",
		Title:  "Title",
	}
	err = uc.UpdateEvent(ctx, nonExistentEvent)
	if !stdErrors.Is(err, errors.ErrEventNotFound) {
		t.Errorf("Expected ErrEventNotFound when updating non-existent event, got %v", err)
	}
}

func TestEventUseCase_DeleteEvent(t *testing.T) {
	uc, ctx := setupTestUseCase()

	event := domain.Event{
		ID:     "test-1",
		UserID: "user-1",
		Date:   "2025-01-15",
		Title:  "Test Event",
	}

	err := uc.CreateEvent(ctx, event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	err = uc.DeleteEvent(ctx, "test-1")
	if err != nil {
		t.Fatalf("Failed to delete event: %v", err)
	}

	err = uc.DeleteEvent(ctx, "non-existent")
	if !stdErrors.Is(err, errors.ErrEventNotFound) {
		t.Errorf("Expected ErrEventNotFound when deleting non-existent event, got %v", err)
	}

	err = uc.DeleteEvent(ctx, "")
	if !stdErrors.Is(err, errors.ErrEmptyEventID) {
		t.Errorf("Expected ErrEmptyEventID when deleting with empty ID, got %v", err)
	}
}

func TestEventUseCase_GetEvents(t *testing.T) {
	uc, ctx := setupTestUseCase()

	events := []domain.Event{
		{ID: "1", UserID: "user-1", Date: "2025-01-15", Title: "Event 1"},
		{ID: "2", UserID: "user-1", Date: "2025-01-15", Title: "Event 2"},
		{ID: "3", UserID: "user-1", Date: "2025-01-16", Title: "Event 3"},
		{ID: "4", UserID: "user-1", Date: "2025-01-22", Title: "Event 4"},
	}

	for _, event := range events {
		err := uc.CreateEvent(ctx, event)
		if err != nil {
			t.Fatalf("Failed to create event: %v", err)
		}
	}

	dayEvents, err := uc.GetEventsForDay(ctx, "user-1", "2025-01-15")
	if err != nil {
		t.Fatalf("Failed to get events for day: %v", err)
	}
	if len(dayEvents) != 2 {
		t.Errorf("Expected 2 events for day, got %d", len(dayEvents))
	}

	weekEvents, err := uc.GetEventsForWeek(ctx, "user-1", "2025-01-15")
	if err != nil {
		t.Fatalf("Failed to get events for week: %v", err)
	}
	if len(weekEvents) != 3 {
		t.Errorf("Expected 3 events for week, got %d", len(weekEvents))
	}

	monthEvents, err := uc.GetEventsForMonth(ctx, "user-1", "2025-01-15")
	if err != nil {
		t.Fatalf("Failed to get events for month: %v", err)
	}
	if len(monthEvents) != 4 {
		t.Errorf("Expected 4 events for month, got %d", len(monthEvents))
	}
}

func TestEventUseCase_GetEvents_EmptyResults(t *testing.T) {
	uc, ctx := setupTestUseCase()

	events, err := uc.GetEventsForDay(ctx, "non-existent-user", "2025-01-15")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}

	events, err = uc.GetEventsForWeek(ctx, "non-existent-user", "2025-01-15")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}

	events, err = uc.GetEventsForMonth(ctx, "non-existent-user", "2025-01-15")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}
}

func TestEventUseCase_Validation_EdgeCases(t *testing.T) {
	uc, ctx := setupTestUseCase()

	testCases := []struct {
		name  string
		event domain.Event
	}{
		{
			name:  "empty user ID",
			event: domain.Event{ID: "test-1", UserID: "", Date: "2025-01-15", Title: "Title"},
		},
		{
			name:  "empty title",
			event: domain.Event{ID: "test-1", UserID: "user-1", Date: "2025-01-15", Title: ""},
		},
		{
			name:  "invalid date format",
			event: domain.Event{ID: "test-1", UserID: "user-1", Date: "2025/01/15", Title: "Title"},
		},
		{
			name:  "empty date",
			event: domain.Event{ID: "test-1", UserID: "user-1", Date: "", Title: "Title"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := uc.CreateEvent(ctx, tc.event)
			if err == nil {
				t.Error("Expected validation error")
			}
		})
	}
}

func TestEventUseCase_DeleteEvent_Validation(t *testing.T) {
	uc, ctx := setupTestUseCase()

	err := uc.DeleteEvent(ctx, "")
	if !stdErrors.Is(err, errors.ErrEmptyEventID) {
		t.Errorf("Expected ErrEmptyEventID, got %v", err)
	}
}

func TestEventUseCase_ContextCancellation(t *testing.T) {
	uc, ctx := setupTestUseCase()

	cancelledCtx, cancel := context.WithCancel(ctx)
	cancel()

	event := domain.Event{
		ID:     "test-1",
		UserID: "user-1",
		Date:   "2025-01-15",
		Title:  "Test Event",
	}

	err := uc.CreateEvent(cancelledCtx, event)
	if err == nil {
		t.Error("Expected error when context is cancelled")
	}
}
