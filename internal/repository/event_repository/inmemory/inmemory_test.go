package inmemory

import (
	"calendar-server/internal/domain"
	"calendar-server/pkg/errors"
	"context"
	stdErrors "errors"
	"testing"

	"go.uber.org/zap"
)

func setupTest() (*EventRepository, context.Context) {
	logger, _ := zap.NewDevelopment()
	repo := NewEventRepository(logger)
	ctx := context.Background()
	return repo, ctx
}

func TestEventRepository_Create(t *testing.T) {
	repo, ctx := setupTest()

	event := domain.Event{
		ID:     "test-1",
		UserID: "user-1",
		Date:   "2025-01-15",
		Title:  "Test Event",
	}

	err := repo.Create(ctx, event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	err = repo.Create(ctx, event)
	if !stdErrors.Is(err, errors.ErrEventConflict) {
		t.Errorf("Expected ErrEventConflict for duplicate event ID, got %v", err)
	}
}

func TestEventRepository_Update(t *testing.T) {
	repo, ctx := setupTest()

	event := domain.Event{
		ID:     "test-1",
		UserID: "user-1",
		Date:   "2025-01-15",
		Title:  "Original Title",
	}

	err := repo.Create(ctx, event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	event.Title = "Updated Title"
	err = repo.Update(ctx, event)
	if err != nil {
		t.Fatalf("Failed to update event: %v", err)
	}

	nonExistentEvent := domain.Event{
		ID:     "non-existent",
		UserID: "user-1",
		Date:   "2025-01-15",
		Title:  "Test",
	}
	err = repo.Update(ctx, nonExistentEvent)
	if !stdErrors.Is(err, errors.ErrEventNotFound) {
		t.Errorf("Expected ErrEventNotFound when updating non-existent event, got %v", err)
	}
}

func TestEventRepository_Delete(t *testing.T) {
	repo, ctx := setupTest()

	event := domain.Event{
		ID:     "test-1",
		UserID: "user-1",
		Date:   "2025-01-15",
		Title:  "Test Event",
	}

	err := repo.Create(ctx, event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	err = repo.Delete(ctx, "test-1")
	if err != nil {
		t.Fatalf("Failed to delete event: %v", err)
	}

	err = repo.Delete(ctx, "non-existent")
	if !stdErrors.Is(err, errors.ErrEventNotFound) {
		t.Errorf("Expected ErrEventNotFound when deleting non-existent event, got %v", err)
	}
}

func TestEventRepository_GetByUserIDAndDate(t *testing.T) {
	repo, ctx := setupTest()

	events := []domain.Event{
		{ID: "1", UserID: "user-1", Date: "2025-01-15", Title: "Event 1"},
		{ID: "2", UserID: "user-1", Date: "2025-01-15", Title: "Event 2"},
		{ID: "3", UserID: "user-2", Date: "2025-01-15", Title: "Event 3"},
		{ID: "4", UserID: "user-1", Date: "2025-01-16", Title: "Event 4"},
	}

	for _, event := range events {
		err := repo.Create(ctx, event)
		if err != nil {
			t.Fatalf("Failed to create event: %v", err)
		}
	}

	result, err := repo.GetByUserIDAndDate(ctx, "user-1", "2025-01-15")
	if err != nil {
		t.Fatalf("Failed to get events: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(result))
	}

	if result[0].ID != "1" || result[1].ID != "2" {
		t.Error("Events are not sorted correctly")
	}
}

func TestEventRepository_GetByUserIDAndWeek(t *testing.T) {
	repo, ctx := setupTest()

	events := []domain.Event{
		{ID: "1", UserID: "user-1", Date: "2025-01-15", Title: "Monday Event"},
		{ID: "2", UserID: "user-1", Date: "2025-01-16", Title: "Tuesday Event"},
		{ID: "3", UserID: "user-1", Date: "2025-01-22", Title: "Next Week Event"},
	}

	for _, event := range events {
		err := repo.Create(ctx, event)
		if err != nil {
			t.Fatalf("Failed to create event: %v", err)
		}
	}

	result, err := repo.GetByUserIDAndWeek(ctx, "user-1", "2025-01-15")
	if err != nil {
		t.Fatalf("Failed to get events for week: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 events for the week, got %d", len(result))
	}
}

func TestEventRepository_GetByUserIDAndMonth(t *testing.T) {
	repo, ctx := setupTest()

	events := []domain.Event{
		{ID: "1", UserID: "user-1", Date: "2025-01-15", Title: "Mid Month Event"},
		{ID: "2", UserID: "user-1", Date: "2025-01-31", Title: "End Month Event"},
		{ID: "3", UserID: "user-1", Date: "2025-02-01", Title: "Next Month Event"},
	}

	for _, event := range events {
		err := repo.Create(ctx, event)
		if err != nil {
			t.Fatalf("Failed to create event: %v", err)
		}
	}

	result, err := repo.GetByUserIDAndMonth(ctx, "user-1", "2025-01-15")
	if err != nil {
		t.Fatalf("Failed to get events for month: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 events for the month, got %d", len(result))
	}
}

func TestEventRepository_GetByUserIDAndDate_Empty(t *testing.T) {
	repo, ctx := setupTest()

	result, err := repo.GetByUserIDAndDate(ctx, "non-existent", "2025-01-15")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 events, got %d", len(result))
	}
}

func TestEventRepository_GetByUserIDAndWeek_Empty(t *testing.T) {
	repo, ctx := setupTest()

	result, err := repo.GetByUserIDAndWeek(ctx, "non-existent", "2025-01-15")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 events, got %d", len(result))
	}
}

func TestEventRepository_GetByUserIDAndMonth_Empty(t *testing.T) {
	repo, ctx := setupTest()

	result, err := repo.GetByUserIDAndMonth(ctx, "non-existent", "2025-01-15")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 events, got %d", len(result))
	}
}

func TestEventRepository_ContextCancellation(t *testing.T) {
	repo, ctx := setupTest()

	cancelledCtx, cancel := context.WithCancel(ctx)
	cancel()

	event := domain.Event{
		ID:     "test-1",
		UserID: "user-1",
		Date:   "2025-01-15",
		Title:  "Test Event",
	}

	err := repo.Create(cancelledCtx, event)
	if err == nil {
		t.Error("Expected error when context is cancelled")
	}

	err = repo.Update(cancelledCtx, event)
	if err == nil {
		t.Error("Expected error when context is cancelled")
	}

	err = repo.Delete(cancelledCtx, "test-1")
	if err == nil {
		t.Error("Expected error when context is cancelled")
	}

	_, err = repo.GetByUserIDAndDate(cancelledCtx, "user-1", "2025-01-15")
	if err == nil {
		t.Error("Expected error when context is cancelled")
	}
}

func TestEventRepository_ConcurrentAccess(t *testing.T) {
	repo, ctx := setupTest()

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			event := domain.Event{
				ID:     string(rune('A' + i)),
				UserID: "user-1",
				Date:   "2025-01-15",
				Title:  "Concurrent Event",
			}
			err := repo.Create(ctx, event)
			if err != nil {
				t.Errorf("Failed to create event in goroutine: %v", err)
			}
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	result, err := repo.GetByUserIDAndDate(ctx, "user-1", "2025-01-15")
	if err != nil {
		t.Fatalf("Failed to get events: %v", err)
	}

	if len(result) != 10 {
		t.Fatalf("Expected 10 events, got %d", len(result))
	}
}
