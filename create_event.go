package userengage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// CreateEvent contains data needed for creation of event
type CreateEvent struct {
	Name      string                 `json:"name"`
	Timestamp uint64                 `json:"timestamp"`
	Client    int                    `json:"client"`
	Data      map[string]interface{} `json:"data"`
}

// CreateEventResponse contains possible errors in request
type CreateEventResponse struct {
	Errors *json.RawMessage `json:"errors"`
}

// CreateEventTimeout creates event
func (c Client) CreateEventTimeout(ctx context.Context, timeout time.Duration, event CreateEvent) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.CreateEvent(timeoutCtx, event)
}

const createEventEndpoint = "/events/"

// CreateEvent creates event
func (c Client) CreateEvent(ctx context.Context, event CreateEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// https://app.userengage.com/api/public/events/
	endpoint := c.apiPrefix + createEventEndpoint

	body := bytes.NewBuffer(payload)
	r, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}

	return c.requestEventCreate(r.WithContext(ctx))
}

// CreateEventCustomUserID creates new event using custom User ID
type CreateEventCustomUserID struct {
	UserID    string                 `json:"user_id"`
	Name      string                 `json:"name"`
	Timestamp uint64                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// CreateEventCustomUserIDTimeout creates event using custom user id
func (c Client) CreateEventCustomUserIDTimeout(ctx context.Context, timeout time.Duration, event CreateEventCustomUserID) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.CreateEventCustomUserID(timeoutCtx, event)
}

const createEventCustomUserIDEndpoint = "/users-by-id/%s/events/"

// CreateEventCustomUserID creates new event using custom User ID
func (c Client) CreateEventCustomUserID(ctx context.Context, event CreateEventCustomUserID) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// https://app.userengage.com/api/public/users-by-id/:user_id/events/
	endpoint := c.apiPrefix + createEventCustomUserIDEndpoint

	body := bytes.NewBuffer(payload)
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf(endpoint, event.UserID), body)
	if err != nil {
		return err
	}

	return c.requestEventCreate(r.WithContext(ctx))
}

func (c Client) requestEventCreate(r *http.Request) error {
	client := http.Client{}

	r.Header.Set("Authorization", "Token "+c.apiKey)
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var eventResponse CreateEventResponse
	err = json.NewDecoder(resp.Body).Decode(&eventResponse)
	if err != nil {
		return err
	}

	if resp.StatusCode == 400 && eventResponse.Errors != nil {
		return errors.New(string(*eventResponse.Errors))
	}
	return statusErrors[resp.StatusCode]
}
