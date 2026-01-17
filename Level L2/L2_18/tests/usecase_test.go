package tests

import (
	"L2_18/internal/entity"
	"L2_18/internal/repository"
	"L2_18/internal/usecase"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
)

// TestSaveEvent_Success - успешное сохранение события
func TestSaveEvent_Success(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	event := entity.Calendar{
		UserID:    "user123",
		NameEvent: "Meeting",
		DataEvent: time.Date(2025, 9, 24, 10, 0, 0, 0, time.UTC),
		Text:      "Team meeting",
	}

	err := uc.SaveEvent(event)
	if err != nil {
		t.Fatalf("SaveEvent() failed: %v", err)
	}
}

// TestGetEventForDay_Success - получение события на день
func TestGetEventForDay_Success(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	event := entity.Calendar{
		UserID:    "user123",
		NameEvent: "Meeting",
		DataEvent: time.Date(2025, 9, 24, 10, 0, 0, 0, time.UTC),
		Text:      "Team meeting",
	}

	uc.SaveEvent(event)

	events, err := uc.GetEventForDay("user123", "2025-09-24")
	if err != nil {
		t.Fatalf("GetEventForDay() failed: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

// TestGetEventForDay_InvalidDate - некорректная дата
func TestGetEventForDay_InvalidDate(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	_, err := uc.GetEventForDay("user123", "invalid")
	if !errors.Is(err, entity.ErrParsing) {
		t.Errorf("expected ErrParsing, got %v", err)
	}
}

// TestGetEventForDay_NoEvents - событий нет
func TestGetEventForDay_NoEvents(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	_, err := uc.GetEventForDay("user123", "2025-09-24")
	if !errors.Is(err, entity.ErrNoEvents) {
		t.Errorf("expected ErrNoEvents, got %v", err)
	}
}

// TestUpdateEvent_Success - успешное обновление
func TestUpdateEvent_Success(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	event := entity.Calendar{
		UserID:    "user123",
		NameEvent: "Meeting",
		DataEvent: time.Date(2025, 9, 24, 10, 0, 0, 0, time.UTC),
		Text:      "Old text",
	}

	uc.SaveEvent(event)

	event.Text = "New text"
	err := uc.UpdateEvent(event)
	if err != nil {
		t.Fatalf("UpdateEvent() failed: %v", err)
	}

	events, _ := uc.GetEventForDay("user123", "2025-09-24")
	if events[0].Text != "New text" {
		t.Errorf("expected 'New text', got '%s'", events[0].Text)
	}
}

// TestUpdateEvent_NotFound - событие не найдено
func TestUpdateEvent_NotFound(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	event := entity.Calendar{
		UserID:    "user123",
		NameEvent: "NonExistent",
		DataEvent: time.Date(2025, 9, 24, 10, 0, 0, 0, time.UTC),
		Text:      "Text",
	}

	err := uc.UpdateEvent(event)
	if !errors.Is(err, entity.ErrNoEvents) {
		t.Errorf("expected ErrNoEvents, got %v", err)
	}
}

// TestDeleteEvent_Success - успешное удаление
func TestDeleteEvent_Success(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	event := entity.Calendar{
		UserID:    "user123",
		NameEvent: "Meeting",
		DataEvent: time.Date(2025, 9, 24, 10, 0, 0, 0, time.UTC),
		Text:      "Text",
	}

	uc.SaveEvent(event)

	err := uc.DeleteEvent(event)
	if err != nil {
		t.Fatalf("DeleteEvent() failed: %v", err)
	}

	_, err = uc.GetEventForDay("user123", "2025-09-24")
	if !errors.Is(err, entity.ErrNoEvents) {
		t.Errorf("event should be deleted")
	}
}

// TestDeleteEvent_NotFound - удаление несуществующего события
func TestDeleteEvent_NotFound(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	event := entity.Calendar{
		UserID:    "user123",
		NameEvent: "NonExistent",
		DataEvent: time.Date(2025, 9, 24, 10, 0, 0, 0, time.UTC),
	}

	err := uc.DeleteEvent(event)
	if !errors.Is(err, entity.ErrNoEvents) {
		t.Errorf("expected ErrNoEvents, got %v", err)
	}
}

// TestGetEventsForWeek_Success - события на неделю
func TestGetEventsForWeek_Success(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	baseDate := time.Date(2025, 9, 24, 10, 0, 0, 0, time.UTC)

	events := []entity.Calendar{
		{
			UserID:    "user123",
			NameEvent: "Event1",
			DataEvent: baseDate,
			Text:      "First",
		},
		{
			UserID:    "user123",
			NameEvent: "Event2",
			DataEvent: baseDate.AddDate(0, 0, 2),
			Text:      "Second",
		},
	}

	for _, e := range events {
		uc.SaveEvent(e)
	}

	result, err := uc.GetEventsForWeek("user123", "2025-09-24")
	if err != nil {
		t.Fatalf("GetEventsForWeek() failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 events, got %d", len(result))
	}
}

// TestGetEventsForWeek_NoEvents - нет событий на неделю
func TestGetEventsForWeek_NoEvents(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	_, err := uc.GetEventsForWeek("user123", "2025-09-24")
	if !errors.Is(err, entity.ErrNoEvents) {
		t.Errorf("expected ErrNoEvents, got %v", err)
	}
}

// TestGetEventForMonth_Success - события на месяц
func TestGetEventForMonth_Success(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	events := []entity.Calendar{
		{
			UserID:    "user123",
			NameEvent: "Event1",
			DataEvent: time.Date(2025, 9, 5, 10, 0, 0, 0, time.UTC),
			Text:      "First",
		},
		{
			UserID:    "user123",
			NameEvent: "Event2",
			DataEvent: time.Date(2025, 9, 20, 10, 0, 0, 0, time.UTC),
			Text:      "Second",
		},
	}

	for _, e := range events {
		uc.SaveEvent(e)
	}

	result, err := uc.GetEventForMonth("user123", "2025-09-15")
	if err != nil {
		t.Fatalf("GetEventForMonth() failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 events, got %d", len(result))
	}
}

// TestGetEventForMonth_NoEvents - нет событий на месяц
func TestGetEventForMonth_NoEvents(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	_, err := uc.GetEventForMonth("user123", "2025-10-15")
	if !errors.Is(err, entity.ErrNoEvents) {
		t.Errorf("expected ErrNoEvents, got %v", err)
	}
}

// TestParseDate_EmptyString - пустая строка
func TestParseDate_EmptyString(t *testing.T) {
	_, err := usecase.ParseDate("")
	if err == nil {
		t.Error("expected error for empty date")
	}
}

// TestParseDate_InvalidLength - некорректная длина
func TestParseDate_InvalidLength(t *testing.T) {
	_, err := usecase.ParseDate("2025-09")
	if err == nil {
		t.Error("expected error for short date")
	}
}

// TestParseDate_InvalidFormat - некорректный формат
func TestParseDate_InvalidFormat(t *testing.T) {
	_, err := usecase.ParseDate("2025-13-45")
	if err == nil {
		t.Error("expected error for invalid date")
	}
}

// TestParseDate_Success - успешный парсинг
func TestParseDate_Success(t *testing.T) {
	date, err := usecase.ParseDate("2025-09-24")
	if err != nil {
		t.Fatalf("ParseDate() failed: %v", err)
	}

	if date.Year() != 2025 || date.Month() != 9 || date.Day() != 24 {
		t.Errorf("unexpected date: %v", date)
	}
}
