package tests

import (
	"L2_18/internal/entity"
	"L2_18/internal/repository"
	"L2_18/internal/usecase"
	"testing"
	"time"

	"go.uber.org/zap"
)

// TestUseCase_SaveAndGetEvent - тест сохранения и получения запроса
func TestUseCase_SaveAndGetEvent(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	event := entity.Calendar{
		UserID:    "123",
		NameEvent: "Meeting",
		DataEvent: time.Date(2025, 9, 24, 10, 0, 0, 0, time.UTC),
		Text:      "Team meeting",
	}

	// Сохраняем событие через UseCase
	if err := uc.SaveEvent(event); err != nil {
		t.Fatalf("SaveEvent failed: %v", err)
	}

	// Получаем события на день
	dateStr := event.DataEvent.Format("2006-01-02")
	events, err := uc.GetEventForDay("123", dateStr)
	if err != nil {
		t.Fatalf("GetEventForDay failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	if events[0].NameEvent != "Meeting" {
		t.Errorf("Expected event name 'Meeting', got %s", events[0].NameEvent)
	}
}

// TestUseCase_UpdateEvent - обновление эвента
func TestUseCase_UpdateEvent(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	event := entity.Calendar{
		UserID:    "123",
		NameEvent: "Meeting",
		DataEvent: time.Now(),
		Text:      "Old text",
	}

	err := uc.SaveEvent(event)
	if err != nil {
		t.Fatalf("SaveEvent failed: %v", err)
	}

	// Обновляем текст события через UseCase
	event.Text = "Updated text"
	if err := uc.UpdateEvent(event); err != nil {
		t.Fatalf("UpdateEvent failed: %v", err)
	}

	dateStr := event.DataEvent.Format("2006-01-02")
	events, _ := uc.GetEventForDay("123", dateStr)
	if events[0].Text != "Updated text" {
		t.Errorf("Expected updated text, got %s", events[0].Text)
	}
}

// TestUseCase_DeleteEvent - удаление события
func TestUseCase_DeleteEvent(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := repository.New(logger)
	uc := usecase.New(repo, logger)

	event := entity.Calendar{
		UserID:    "123",
		NameEvent: "Meeting",
		DataEvent: time.Now(),
		Text:      "Text",
	}

	err := uc.SaveEvent(event)
	if err != nil {
		t.Fatalf("SaveEvent failed: %v", err)
	}

	// Удаляем событие через UseCase
	if err := uc.DeleteEvent(event); err != nil {
		t.Fatalf("DeleteEvent failed: %v", err)
	}

	dateStr := event.DataEvent.Format("2006-01-02")
	events, err := uc.GetEventForDay("123", dateStr)
	if err == nil && len(events) > 0 {
		t.Errorf("Expected no events, but got %d", len(events))
	}
}
