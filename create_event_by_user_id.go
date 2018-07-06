package userengage

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

const createEventEndpoint = "https://app.userengage.com/api/public/events/"

// CreateEvent contains data needed for creation of event
type CreateEvent struct {
	Name      string                 `json:"name"`
	Timestamp uint64                 `json:"timestamp"`
	Client    int                    `json:"client"`
	Data      map[string]interface{} `json:"data"`
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

	request.Header.Set("Authorization", "Token "+c.apikey)
	request.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(request.WithContext(ctx))
	return err
}
