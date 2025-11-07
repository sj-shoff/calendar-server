package event_handler

import (
	"bytes"
	"calendar-server/internal/domain"
	"calendar-server/pkg/errors"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	uc "calendar-server/internal/usecase/event_usecase"

	"go.uber.org/zap"
)

type mockEventUseCase struct {
	events map[string]domain.Event
}

func newMockEventUseCase() uc.EventUseCaseContract {
	return &mockEventUseCase{
		events: make(map[string]domain.Event),
	}
}

func (m *mockEventUseCase) CreateEvent(ctx context.Context, event domain.Event) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if event.ID == "" {
		return errors.ErrEmptyEventID
	}
	if event.UserID == "" {
		return errors.ErrEmptyUserID
	}
	if event.Title == "" {
		return errors.ErrEmptyTitle
	}
	if event.Date == "" {
		return errors.ErrInvalidDate
	}

	if _, exists := m.events[event.ID]; exists {
		return errors.ErrEventConflict
	}

	m.events[event.ID] = event
	return nil
}

func (m *mockEventUseCase) UpdateEvent(ctx context.Context, event domain.Event) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, exists := m.events[event.ID]; !exists {
		return errors.ErrEventNotFound
	}

	m.events[event.ID] = event
	return nil
}

func (m *mockEventUseCase) DeleteEvent(ctx context.Context, eventID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, exists := m.events[eventID]; !exists {
		return errors.ErrEventNotFound
	}

	delete(m.events, eventID)
	return nil
}

func (m *mockEventUseCase) GetEventsForDay(ctx context.Context, userID, date string) ([]domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var result []domain.Event
	for _, event := range m.events {
		if event.UserID == userID && event.Date == date {
			result = append(result, event)
		}
	}
	return result, nil
}

func (m *mockEventUseCase) GetEventsForWeek(ctx context.Context, userID, date string) ([]domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var result []domain.Event
	for _, event := range m.events {
		if event.UserID == userID {
			result = append(result, event)
		}
	}
	return result, nil
}

func (m *mockEventUseCase) GetEventsForMonth(ctx context.Context, userID, date string) ([]domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var result []domain.Event
	for _, event := range m.events {
		if event.UserID == userID {
			result = append(result, event)
		}
	}
	return result, nil
}

func setupTestHandler() *EventHandler {
	logger, _ := zap.NewDevelopment()
	eventUseCase := newMockEventUseCase()
	return NewEventHandler(eventUseCase, logger)
}

func TestEventHandler_CreateEvent(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		payload        interface{}
		contentType    string
		expectedStatus int
	}{
		{
			name: "valid event creation",
			payload: map[string]string{
				"id":      "test-1",
				"user_id": "user-1",
				"date":    "2025-01-15",
				"title":   "Test Event",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid content type",
			payload: map[string]string{
				"id":      "test-2",
				"user_id": "user-1",
				"date":    "2025-01-15",
				"title":   "Test Event",
			},
			contentType:    "text/plain",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty event ID",
			payload: map[string]string{
				"id":      "",
				"user_id": "user-1",
				"date":    "2025-01-15",
				"title":   "Test Event",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			payload:        "invalid json",
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			switch v := tt.payload.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", tt.contentType)
			rr := httptest.NewRecorder()

			handler.CreateEvent(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestEventHandler_CreateEvent_ErrorCases(t *testing.T) {
	handler := setupTestHandler()

	// Test with wrong content type
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()
	handler.CreateEvent(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for wrong content type, got %d", rr.Code)
	}

	// Test with invalid JSON
	req = httptest.NewRequest("POST", "/create_event", bytes.NewBufferString(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler.CreateEvent(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", rr.Code)
	}
}

func TestEventHandler_UpdateEvent(t *testing.T) {
	handler := setupTestHandler()

	createPayload := map[string]string{
		"id":      "test-1",
		"user_id": "user-1",
		"date":    "2025-01-15",
		"title":   "Original Title",
	}

	body, _ := json.Marshal(createPayload)
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.CreateEvent(rr, req)

	updatePayload := map[string]string{
		"id":      "test-1",
		"user_id": "user-1",
		"date":    "2025-01-15",
		"title":   "Updated Title",
	}

	body, _ = json.Marshal(updatePayload)
	req = httptest.NewRequest("POST", "/update_event", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	handler.UpdateEvent(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response Response
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Result != "event updated" {
		t.Errorf("Expected result 'event updated', got %v", response.Result)
	}
}

func TestEventHandler_UpdateEvent_NotFound(t *testing.T) {
	handler := setupTestHandler()

	payload := map[string]string{
		"id":      "non-existent",
		"user_id": "user-1",
		"date":    "2025-01-15",
		"title":   "Test Event",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/update_event", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.UpdateEvent(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 503 for non-existent event, got %d", rr.Code)
	}
}

func TestEventHandler_DeleteEvent(t *testing.T) {
	handler := setupTestHandler()

	createPayload := map[string]string{
		"id":      "test-1",
		"user_id": "user-1",
		"date":    "2025-01-15",
		"title":   "Test Event",
	}

	body, _ := json.Marshal(createPayload)
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.CreateEvent(rr, req)

	deletePayload := map[string]string{
		"id": "test-1",
	}

	body, _ = json.Marshal(deletePayload)
	req = httptest.NewRequest("POST", "/delete_event", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	handler.DeleteEvent(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response Response
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Result != "event deleted" {
		t.Errorf("Expected result 'event deleted', got %v", response.Result)
	}
}

func TestEventHandler_DeleteEvent_NotFound(t *testing.T) {
	handler := setupTestHandler()

	payload := map[string]string{"id": "non-existent"}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/delete_event", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.DeleteEvent(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 503 for non-existent event, got %d", rr.Code)
	}
}

func TestEventHandler_GetEvents(t *testing.T) {
	handler := setupTestHandler()

	events := []map[string]string{
		{"id": "1", "user_id": "user-1", "date": "2025-01-15", "title": "Event 1"},
		{"id": "2", "user_id": "user-1", "date": "2025-01-15", "title": "Event 2"},
	}

	for _, event := range events {
		body, _ := json.Marshal(event)
		req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.CreateEvent(rr, req)
	}

	req := httptest.NewRequest("GET", "/events_for_day?user_id=user-1&date=2025-01-15", nil)
	rr := httptest.NewRecorder()

	handler.EventsForDay(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response Response
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	eventsResult, ok := response.Result.([]interface{})
	if !ok {
		t.Fatal("Expected result to be an array")
	}

	if len(eventsResult) != 2 {
		t.Errorf("Expected 2 events, got %d", len(eventsResult))
	}
}

func TestEventHandler_GetEvents_MissingParameters(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "missing user_id",
			url:            "/events_for_day?date=2025-01-15",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing date",
			url:            "/events_for_day?user_id=user-1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "both parameters present",
			url:            "/events_for_day?user_id=user-1&date=2025-01-15",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			rr := httptest.NewRecorder()

			handler.EventsForDay(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestEventHandler_ErrorHandling(t *testing.T) {
	handler := setupTestHandler()

	payload := map[string]string{
		"id":      "non-existent",
		"user_id": "user-1",
		"date":    "2025-01-15",
		"title":   "Test Event",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/update_event", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.UpdateEvent(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 503 for non-existent event, got %d", rr.Code)
	}
}

func TestEventHandler_GetEventsForWeek(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("GET", "/events_for_week?user_id=user-1&date=2025-01-15", nil)
	rr := httptest.NewRecorder()

	handler.EventsForWeek(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestEventHandler_GetEventsForMonth(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("GET", "/events_for_month?user_id=user-1&date=2025-01-15", nil)
	rr := httptest.NewRecorder()

	handler.EventsForMonth(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}
