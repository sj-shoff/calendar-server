package event_usecase

import (
	"calendar-server/internal/domain"
	"calendar-server/pkg/errors"
	"time"
)

// validateEvent проверяет валидность всей структуры Event.
func (uc *EventUseCase) validateEvent(event domain.Event) error {
	if event.ID == "" {
		return errors.ErrEmptyEventID
	}
	if event.UserID == "" {
		return errors.ErrEmptyUserID
	}
	if event.Title == "" {
		return errors.ErrEmptyTitle
	}
	if !isValidDate(event.Date) {
		return errors.ErrInvalidDate
	}
	return nil
}

// ValidateEventID проверяет валидность ID события.
func (uc *EventUseCase) validateEventID(id string) error {
	if id == "" {
		return errors.ErrEmptyEventID
	}
	return nil
}

// ValidateUserID проверяет валидность UserID.
func (uc *EventUseCase) validateUserID(userID string) error {
	if userID == "" {
		return errors.ErrEmptyUserID
	}
	return nil
}

// ValidateDate проверяет валидность формата даты.
func (uc *EventUseCase) validateDate(date string) error {
	if !isValidDate(date) {
		return errors.ErrInvalidDate
	}
	return nil
}

// isValidDate проверяет корректность формата даты "YYYY-MM-DD".
func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
