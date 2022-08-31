package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	v1 "github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/api/v1"
	"github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/storage"

	"github.com/go-chi/chi/v5" //nolint:typecheck
)

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) (int, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, eventID int) error
	ListEvent(ctx context.Context) ([]storage.Event, error)
	GetEvent(ctx context.Context, eventID int) (*storage.Event, error)
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type Handlers struct {
	app    Application
	logger Logger
}

func NewHandlers(app Application, logger Logger) *Handlers {
	return &Handlers{
		app:    app,
		logger: logger,
	}
}

func (h *Handlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var request v1.EventCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	event := storage.Event{
		UserID:      request.UserID,
		Title:       request.Title,
		Description: request.Description,
		StartAt:     time.Unix(request.StartAt, 0),
		EndAt:       time.Unix(request.EndAt, 0),
		NotifyAt:    v1.GetNotifyAt(request.NotifyAt),
	}

	eventID, err := h.app.CreateEvent(ctx, event)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := v1.EventResponse{
		ID:          eventID,
		UserID:      event.UserID,
		Title:       event.Title,
		Description: event.Description,
		StartAt:     event.StartAt,
		EndAt:       event.EndAt,
		NotifyAt:    event.NotifyAt,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.Atoi(chi.URLParam(r, "eventID")) //nolint:typecheck
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	event, err := h.app.GetEvent(ctx, eventID)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if event == nil {
		h.logger.Error("event not found")
		http.Error(w, "event not found", http.StatusNotFound)
		return
	}

	var request v1.EventUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	event.Title = request.Title
	event.Description = request.Description
	event.StartAt = time.Unix(request.StartAt, 0)
	event.EndAt = time.Unix(request.EndAt, 0)
	event.NotifyAt = v1.GetNotifyAt(request.NotifyAt)

	if err := h.app.UpdateEvent(ctx, *event); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := v1.EventResponse{
		ID:          eventID,
		UserID:      event.UserID,
		Title:       event.Title,
		Description: event.Description,
		StartAt:     event.StartAt,
		EndAt:       event.EndAt,
		NotifyAt:    event.NotifyAt,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.Atoi(chi.URLParam(r, "eventID")) //nolint:typecheck
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := h.app.DeleteEvent(context.Background(), eventID); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) ListEvent(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	events, err := h.app.ListEvent(ctx)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response v1.EventListResponse
	response.Events = make([]v1.EventResponse, 0)
	for _, event := range events { //nolint:typecheck
		response.Events = append(response.Events, v1.EventResponse{
			ID:          event.ID,
			UserID:      event.UserID,
			Title:       event.Title,
			Description: event.Description,
			StartAt:     event.StartAt,
			EndAt:       event.EndAt,
			NotifyAt:    event.NotifyAt,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.Atoi(chi.URLParam(r, "eventID")) //nolint:typecheck
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	event, err := h.app.GetEvent(ctx, eventID)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if event == nil {
		h.logger.Error("event not found")
		http.Error(w, "event not found", http.StatusNotFound)
		return
	}

	response := v1.EventResponse{
		ID:          event.ID,
		UserID:      event.UserID,
		Title:       event.Title,
		Description: event.Description,
		StartAt:     event.StartAt,
		EndAt:       event.EndAt,
		NotifyAt:    event.NotifyAt,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
