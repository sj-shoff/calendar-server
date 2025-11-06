package handler

import (
	"calendar-server/pkg/errors"
	"encoding/json"
	"net/http"

	"calendar-server/internal/domain"
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
	eventUseCase *uc.EventUseCase
	logger       *zap.Logger
}

// NewEventHandler - конструктор обработчика событий
func NewEventHandler(eventUseCase *uc.EventUseCase, logger *zap.Logger) *EventHandler {
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
		writeError(w, errors.ErrUnsupportedMedia.Error(), http.StatusBadRequest)
		return
	}

	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.logger.Warn("Invalid JSON format",
			zappretty.Field("error", err),
		)
		writeError(w, errors.ErrInvalidJSON.Error(), http.StatusBadRequest)
		return
	}

	if err := h.eventUseCase.CreateEvent(ctx, event); err != nil {
		h.logger.Error("Failed to create event",
			zappretty.Field("error", err),
			zappretty.Field("event_id", event.ID),
			zappretty.Field("user_id", event.UserID),
		)
		handleCalendarError(w, err)
		return
	}

	h.logger.Info("Event created successfully",
		zappretty.Field("event_id", event.ID),
		zappretty.Field("user_id", event.UserID),
	)
	writeResponse(w, Response{Result: "event created"})
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
		writeError(w, errors.ErrUnsupportedMedia.Error(), http.StatusBadRequest)
		return
	}

	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.logger.Warn("Invalid JSON format", zappretty.Field("error", err))
		writeError(w, errors.ErrInvalidJSON.Error(), http.StatusBadRequest)
		return
	}

	if err := h.eventUseCase.UpdateEvent(ctx, event); err != nil {
		h.logger.Error("Failed to update event",
			zappretty.Field("error", err),
			zappretty.Field("event_id", event.ID),
		)
		handleCalendarError(w, err)
		return
	}

	h.logger.Info("Event updated successfully",
		zappretty.Field("event_id", event.ID),
	)
	writeResponse(w, Response{Result: "event updated"})
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
		writeError(w, errors.ErrUnsupportedMedia.Error(), http.StatusBadRequest)
		return
	}

	var request struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Warn("Invalid JSON format", zappretty.Field("error", err))
		writeError(w, errors.ErrInvalidJSON.Error(), http.StatusBadRequest)
		return
	}

	if err := h.eventUseCase.DeleteEvent(ctx, request.ID); err != nil {
		h.logger.Error("Failed to delete event",
			zappretty.Field("error", err),
			zappretty.Field("event_id", request.ID),
		)
		handleCalendarError(w, err)
		return
	}

	h.logger.Info("Event deleted successfully",
		zappretty.Field("event_id", request.ID),
	)
	writeResponse(w, Response{Result: "event deleted"})
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
		writeError(w, errors.ErrMissingParameters.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.eventUseCase.GetEventsForDay(ctx, userID, date)
	if err != nil {
		h.logger.Error("Failed to get events for day",
			zappretty.Field("error", err),
			zappretty.Field("user_id", userID),
			zappretty.Field("date", date),
		)
		handleCalendarError(w, err)
		return
	}

	h.logger.Debug("Retrieved events for day",
		zappretty.Field("user_id", userID),
		zappretty.Field("date", date),
		zappretty.Field("count", len(events)),
	)
	writeResponse(w, Response{Result: events})
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
		writeError(w, errors.ErrMissingParameters.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.eventUseCase.GetEventsForWeek(ctx, userID, date)
	if err != nil {
		h.logger.Error("Failed to get events for week",
			zappretty.Field("error", err),
			zappretty.Field("user_id", userID),
			zappretty.Field("date", date),
		)
		handleCalendarError(w, err)
		return
	}

	h.logger.Debug("Retrieved events for week",
		zappretty.Field("user_id", userID),
		zappretty.Field("date", date),
		zappretty.Field("count", len(events)),
	)
	writeResponse(w, Response{Result: events})
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
		writeError(w, errors.ErrMissingParameters.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.eventUseCase.GetEventsForMonth(ctx, userID, date)
	if err != nil {
		h.logger.Error("Failed to get events for month",
			zappretty.Field("error", err),
			zappretty.Field("user_id", userID),
			zappretty.Field("date", date),
		)
		handleCalendarError(w, err)
		return
	}

	h.logger.Debug("Retrieved events for month",
		zappretty.Field("user_id", userID),
		zappretty.Field("date", date),
		zappretty.Field("count", len(events)),
	)
	writeResponse(w, Response{Result: events})
}

// handleCalendarError - обработчик ошибок календаря
func handleCalendarError(w http.ResponseWriter, err error) {
	switch err {
	case errors.ErrEventNotFound:
		writeError(w, err.Error(), http.StatusServiceUnavailable)

	case errors.ErrInvalidDate,
		errors.ErrEmptyEventID,
		errors.ErrEmptyUserID,
		errors.ErrEmptyTitle:
		writeError(w, err.Error(), http.StatusBadRequest)

	case errors.ErrEventConflict:
		writeError(w, err.Error(), http.StatusConflict)

	default:
		writeError(w, "internal server error", http.StatusInternalServerError)
	}
}

// writeResponse - функция для записи ответа
func writeResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// writeError - функция для записи ошибки
func writeError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{Error: errorMsg})
}
