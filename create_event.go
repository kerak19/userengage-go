package userengage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const createEventEndpoint = "https://app.userengage.com/api/public/events/"

// CreateEvent contains data needed for creation of event
type CreateEvent struct {
	Name      string                 `json:"name"`
	Timestamp uint64                 `json:"timestamp"`
	Client    int                    `json:"client"`
	Data      map[string]interface{} `json:"data"`
}

// createEventResponse contains possible errors in request
type createEventResponse struct {
	Errors *json.RawMessage `json:"errors"`
}

// CreateEventTimeout creates event using user id
func (c Client) CreateEventTimeout(ctx context.Context, timeout time.Duration,
	event CreateEvent) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.CreateEvent(timeoutCtx, event)
}

// CreateEvent creates event using user id
func (c Client) CreateEvent(ctx context.Context, event CreateEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	requestBody := bytes.NewBuffer(payload)
	request, err := http.NewRequest(http.MethodPost, createEventEndpoint, requestBody)
	if err != nil {
		return err
	}

	client := http.Client{}

	request.Header.Set("Authorization", "Token "+c.apikey)
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var eventResponse createEventResponse
	err = json.NewDecoder(resp.Body).Decode(&eventResponse)
	if err != nil {
		return err
	}

	if resp.StatusCode == 400 && eventResponse.Errors != nil {
		return errors.New(string(*eventResponse.Errors))
	}
	return statusErrors[resp.StatusCode]
}
