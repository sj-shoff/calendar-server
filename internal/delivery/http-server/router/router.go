package router

import (
	"net/http"

	handler "calendar-server/internal/delivery/http-server/handler/event_handler"
	"calendar-server/internal/delivery/http-server/middleware"

	"go.uber.org/zap"
)

func NewRouter(eventHandler *handler.EventHandler, logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /create_event", eventHandler.CreateEvent)
	mux.HandleFunc("POST /update_event", eventHandler.UpdateEvent)
	mux.HandleFunc("POST /delete_event", eventHandler.DeleteEvent)
	mux.HandleFunc("GET /events_for_day", eventHandler.EventsForDay)
	mux.HandleFunc("GET /events_for_week", eventHandler.EventsForWeek)
	mux.HandleFunc("GET /events_for_month", eventHandler.EventsForMonth)

	return middleware.LoggingMiddleware(logger, mux)
}
