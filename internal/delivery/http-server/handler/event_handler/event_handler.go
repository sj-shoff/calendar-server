package event_handler

import (
	"calendar-server/internal/domain"
	"calendar-server/pkg/errors"
	"encoding/json"
	stdErrors "errors"
	"net/http"

	uc "calendar-server/internal/usecase/event_usecase"
	"calendar-server/pkg/logger/zappretty"

	"go.uber.org/zap"
)

// Response - структура ответа
type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// EventHandler - обработчик событий
type EventHandler struct {
	eventUseCase uc.EventUseCaseContract
	logger       *zap.Logger
}

// NewEventHandler - конструктор обработчика событий
func NewEventHandler(eventUseCase uc.EventUseCaseContract, logger *zap.Logger) *EventHandler {
	return &EventHandler{
		eventUseCase: eventUseCase,
		logger:       logger,
	}
}

// CreateEvent - метод создания события
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Debug("Creating event",
		zappretty.Field("method", r.Method),
		zappretty.Field("path", r.URL.Path),
	)

	if r.Header.Get("Content-Type") != "application/json" {
		h.logger.Warn("Unsupported media type",
			zappretty.Field("content_type", r.Header.Get("Content-Type")),
		)
		h.writeError(w, errors.ErrUnsupportedMedia.Error(), http.StatusBadRequest)
		return
	}

	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.logger.Warn("Invalid JSON format",
			zappretty.Field("error", err),
		)
		h.writeError(w, errors.ErrInvalidJSON.Error(), http.StatusBadRequest)
		return
	}

	if err := h.eventUseCase.CreateEvent(ctx, event); err != nil {
		h.logger.Error("Failed to create event",
			zappretty.Field("error", err),
			zappretty.Field("event_id", event.ID),
			zappretty.Field("user_id", event.UserID),
		)
		h.handleCalendarError(w, err)
		return
	}

	h.logger.Info("Event created successfully",
		zappretty.Field("event_id", event.ID),
		zappretty.Field("user_id", event.UserID),
	)
	h.writeResponse(w, Response{Result: "event created"})
}

// UpdateEvent - метод обновления события
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Debug("Updating event",
		zappretty.Field("method", r.Method),
		zappretty.Field("path", r.URL.Path),
	)

	if r.Header.Get("Content-Type") != "application/json" {
		h.logger.Warn("Unsupported media type")
		h.writeError(w, errors.ErrUnsupportedMedia.Error(), http.StatusBadRequest)
		return
	}

	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.logger.Warn("Invalid JSON format", zappretty.Field("error", err))
		h.writeError(w, errors.ErrInvalidJSON.Error(), http.StatusBadRequest)
		return
	}

	if err := h.eventUseCase.UpdateEvent(ctx, event); err != nil {
		h.logger.Error("Failed to update event",
			zappretty.Field("error", err),
			zappretty.Field("event_id", event.ID),
		)
		h.handleCalendarError(w, err)
		return
	}

	h.logger.Info("Event updated successfully",
		zappretty.Field("event_id", event.ID),
	)
	h.writeResponse(w, Response{Result: "event updated"})
}

// DeleteEvent - метод удаления события
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Debug("Deleting event",
		zappretty.Field("method", r.Method),
		zappretty.Field("path", r.URL.Path),
	)

	if r.Header.Get("Content-Type") != "application/json" {
		h.logger.Warn("Unsupported media type")
		h.writeError(w, errors.ErrUnsupportedMedia.Error(), http.StatusBadRequest)
		return
	}

	var request struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Warn("Invalid JSON format", zappretty.Field("error", err))
		h.writeError(w, errors.ErrInvalidJSON.Error(), http.StatusBadRequest)
		return
	}

	if err := h.eventUseCase.DeleteEvent(ctx, request.ID); err != nil {
		h.logger.Error("Failed to delete event",
			zappretty.Field("error", err),
			zappretty.Field("event_id", request.ID),
		)
		h.handleCalendarError(w, err)
		return
	}

	h.logger.Info("Event deleted successfully",
		zappretty.Field("event_id", request.ID),
	)
	h.writeResponse(w, Response{Result: "event deleted"})
}

// EventsForDay - метод получения событий за день
func (h *EventHandler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	h.logger.Debug("Getting events for day",
		zappretty.Field("method", r.Method),
		zappretty.Field("path", r.URL.Path),
		zappretty.Field("user_id", userID),
		zappretty.Field("date", date),
	)

	if userID == "" || date == "" {
		h.logger.Warn("Missing parameters for events for day")
		h.writeError(w, errors.ErrMissingParameters.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.eventUseCase.GetEventsForDay(ctx, userID, date)
	if err != nil {
		h.logger.Error("Failed to get events for day",
			zappretty.Field("error", err),
			zappretty.Field("user_id", userID),
			zappretty.Field("date", date),
		)
		h.handleCalendarError(w, err)
		return
	}

	h.logger.Debug("Retrieved events for day",
		zappretty.Field("user_id", userID),
		zappretty.Field("date", date),
		zappretty.Field("count", len(events)),
	)
	h.writeResponse(w, Response{Result: events})
}

// EventsForWeek - метод получения событий за неделю
func (h *EventHandler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	h.logger.Debug("Getting events for week",
		zappretty.Field("method", r.Method),
		zappretty.Field("path", r.URL.Path),
		zappretty.Field("user_id", userID),
		zappretty.Field("date", date),
	)

	if userID == "" || date == "" {
		h.logger.Warn("Missing parameters for events for week")
		h.writeError(w, errors.ErrMissingParameters.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.eventUseCase.GetEventsForWeek(ctx, userID, date)
	if err != nil {
		h.logger.Error("Failed to get events for week",
			zappretty.Field("error", err),
			zappretty.Field("user_id", userID),
			zappretty.Field("date", date),
		)
		h.handleCalendarError(w, err)
		return
	}

	h.logger.Debug("Retrieved events for week",
		zappretty.Field("user_id", userID),
		zappretty.Field("date", date),
		zappretty.Field("count", len(events)),
	)
	h.writeResponse(w, Response{Result: events})
}

// EventsForMonth - метод получения событий за месяц
func (h *EventHandler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	h.logger.Debug("Getting events for month",
		zappretty.Field("method", r.Method),
		zappretty.Field("path", r.URL.Path),
		zappretty.Field("user_id", userID),
		zappretty.Field("date", date),
	)

	if userID == "" || date == "" {
		h.logger.Warn("Missing parameters for events for month")
		h.writeError(w, errors.ErrMissingParameters.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.eventUseCase.GetEventsForMonth(ctx, userID, date)
	if err != nil {
		h.logger.Error("Failed to get events for month",
			zappretty.Field("error", err),
			zappretty.Field("user_id", userID),
			zappretty.Field("date", date),
		)
		h.handleCalendarError(w, err)
		return
	}

	h.logger.Debug("Retrieved events for month",
		zappretty.Field("user_id", userID),
		zappretty.Field("date", date),
		zappretty.Field("count", len(events)),
	)
	h.writeResponse(w, Response{Result: events})
}

// handleCalendarError - обработчик ошибок календаря
func (h *EventHandler) handleCalendarError(w http.ResponseWriter, err error) {
	switch {
	case stdErrors.Is(err, errors.ErrEventNotFound):
		h.writeError(w, err.Error(), http.StatusServiceUnavailable)

	case stdErrors.Is(err, errors.ErrInvalidDate),
		stdErrors.Is(err, errors.ErrEmptyEventID),
		stdErrors.Is(err, errors.ErrEmptyUserID),
		stdErrors.Is(err, errors.ErrEmptyTitle):
		h.writeError(w, err.Error(), http.StatusBadRequest)

	case stdErrors.Is(err, errors.ErrEventConflict):
		h.writeError(w, err.Error(), http.StatusConflict)

	default:
		h.writeError(w, "internal server error", http.StatusInternalServerError)
	}
}

// writeResponse - функция для записи ответа
func (h *EventHandler) writeResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// writeError - функция для записи ошибки
func (h *EventHandler) writeError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(Response{Error: errorMsg}); err != nil {
		h.logger.Error("Failed to encode error response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
