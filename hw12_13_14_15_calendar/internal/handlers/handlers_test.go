package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/api/v1"
	"github.com/TOIIIA86/hw-otus/hw12_13_14_15_calendar/internal/storage"

	"github.com/stretchr/testify/require"
)

func TestHandlers_CreateEvent(t *testing.T) {
	server := httptest.NewTLSServer(router)
	defer server.Close()

	var body bytes.Buffer
	request := v1.EventCreateRequest{
		UserID:      100500,
		Title:       "title 100500",
		Description: stringPtr("description 100500"),
		StartAt:     1657128600,
		EndAt:       1657130400,
		NotifyAt:    int64Ptr(1657126800),
	}
	err := json.NewEncoder(&body).Encode(request)
	require.NoError(t, err)

	res, err := server.Client().Post(server.URL+"/events", "application/json", &body)
	require.NoError(t, err)

	require.Equal(t, 200, res.StatusCode)

	var response v1.EventResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	require.NoError(t, err)

	err = store.DeleteEvent(context.Background(), response.ID)
	require.NoError(t, err)

	expected := v1.EventResponse{
		ID:          response.ID,
		UserID:      100500,
		Title:       "title 100500",
		Description: stringPtr("description 100500"),
		StartAt:     time.Unix(1657128600, 0),
		EndAt:       time.Unix(1657130400, 0),
		NotifyAt:    timePtr(time.Unix(1657126800, 0)),
	}

	require.Equal(t, expected, response)
}

func TestHandlers_UpdateEvent(t *testing.T) {
	server := httptest.NewTLSServer(router)
	defer server.Close()

	eventID, err := store.CreateEvent(context.Background(), storage.Event{
		UserID:      100500,
		Title:       "title 100500",
		Description: stringPtr("description 100500"),
		StartAt:     time.Unix(1657128600, 0),
		EndAt:       time.Unix(1657130400, 0),
		NotifyAt:    timePtr(time.Unix(1657126800, 0)),
	})
	require.NoError(t, err)

	var body bytes.Buffer
	request := v1.EventUpdateRequest{
		Title:       "title new",
		Description: stringPtr("description new"),
		StartAt:     1657128601,
		EndAt:       1657130401,
		NotifyAt:    int64Ptr(1657126801),
	}
	err = json.NewEncoder(&body).Encode(request)
	require.NoError(t, err)

	client := server.Client()
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/events/%d", server.URL, eventID), &body)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	require.NoError(t, err)

	require.Equal(t, 200, res.StatusCode)

	var response v1.EventResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	require.NoError(t, err)

	err = store.DeleteEvent(context.Background(), response.ID)
	require.NoError(t, err)

	expected := v1.EventResponse{
		ID:          eventID,
		UserID:      100500,
		Title:       "title new",
		Description: stringPtr("description new"),
		StartAt:     time.Unix(1657128601, 0),
		EndAt:       time.Unix(1657130401, 0),
		NotifyAt:    timePtr(time.Unix(1657126801, 0)),
	}

	require.Equal(t, expected, response)
}

func TestHandlers_DeleteEvent(t *testing.T) {
	server := httptest.NewTLSServer(router)
	defer server.Close()

	eventID, err := store.CreateEvent(context.Background(), storage.Event{
		UserID:      100500,
		Title:       "title 100500",
		Description: stringPtr("description 100500"),
		StartAt:     time.Unix(1657128600, 0),
		EndAt:       time.Unix(1657130400, 0),
		NotifyAt:    timePtr(time.Unix(1657126800, 0)),
	})
	require.NoError(t, err)

	client := server.Client()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/events/%d", server.URL, eventID), nil)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	require.NoError(t, err)

	require.Equal(t, 200, res.StatusCode)
}

func TestHandlers_ListEvent(t *testing.T) {
	server := httptest.NewTLSServer(router)
	defer server.Close()

	event1ID, err := store.CreateEvent(context.Background(), storage.Event{
		UserID:      100501,
		Title:       "title 1",
		Description: stringPtr("description 1"),
		StartAt:     time.Unix(1657128601, 0),
		EndAt:       time.Unix(1657130401, 0),
		NotifyAt:    timePtr(time.Unix(1657126801, 0)),
	})
	require.NoError(t, err)

	event2ID, err := store.CreateEvent(context.Background(), storage.Event{
		UserID:      100502,
		Title:       "title 2",
		Description: stringPtr("description 2"),
		StartAt:     time.Unix(1657128602, 0),
		EndAt:       time.Unix(1657130402, 0),
		NotifyAt:    timePtr(time.Unix(1657126802, 0)),
	})
	require.NoError(t, err)

	res, err := server.Client().Get(server.URL + "/events")
	require.NoError(t, err)

	require.Equal(t, 200, res.StatusCode)

	var response v1.EventListResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	require.NoError(t, err)

	err = store.DeleteEvent(context.Background(), event1ID)
	require.NoError(t, err)

	err = store.DeleteEvent(context.Background(), event2ID)
	require.NoError(t, err)

	expected := v1.EventListResponse{Events: []v1.EventResponse{
		{
			ID:          event1ID,
			UserID:      100501,
			Title:       "title 1",
			Description: stringPtr("description 1"),
			StartAt:     time.Unix(1657128601, 0),
			EndAt:       time.Unix(1657130401, 0),
			NotifyAt:    timePtr(time.Unix(1657126801, 0)),
		},
		{
			ID:          event2ID,
			UserID:      100502,
			Title:       "title 2",
			Description: stringPtr("description 2"),
			StartAt:     time.Unix(1657128602, 0),
			EndAt:       time.Unix(1657130402, 0),
			NotifyAt:    timePtr(time.Unix(1657126802, 0)),
		},
	}}

	require.Equal(t, expected, response)
}

func TestHandlers_GetEvent(t *testing.T) {
	server := httptest.NewTLSServer(router)
	defer server.Close()

	eventID, err := store.CreateEvent(context.Background(), storage.Event{
		UserID:      100500,
		Title:       "title 100500",
		Description: stringPtr("description 100500"),
		StartAt:     time.Unix(1657128600, 0),
		EndAt:       time.Unix(1657130400, 0),
		NotifyAt:    timePtr(time.Unix(1657126800, 0)),
	})
	require.NoError(t, err)

	res, err := server.Client().Get(fmt.Sprintf("%s/events/%d", server.URL, eventID))
	require.NoError(t, err)

	require.Equal(t, 200, res.StatusCode)

	var response v1.EventResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	require.NoError(t, err)

	err = store.DeleteEvent(context.Background(), response.ID)
	require.NoError(t, err)

	expected := v1.EventResponse{
		ID:          response.ID,
		UserID:      100500,
		Title:       "title 100500",
		Description: stringPtr("description 100500"),
		StartAt:     time.Unix(1657128600, 0),
		EndAt:       time.Unix(1657130400, 0),
		NotifyAt:    timePtr(time.Unix(1657126800, 0)),
	}

	require.Equal(t, expected, response)
}

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(ts int64) *int64 {
	return &ts
}

func timePtr(t time.Time) *time.Time {
	return &t
}
