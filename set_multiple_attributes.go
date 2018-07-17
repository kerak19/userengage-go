package userengage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const setMultipleAttributes = "https://app.userengage.com/api/public/users/%d/set_multiple_attributes/"

// Attributes contains attributes which are going to be changed
type Attributes map[string]interface{}

// setMultipleAttributesResponse contains possible errors in request
type setMultipleAttributesResponse struct {
	Errors *json.RawMessage `json:"errors"`
}

// SetMultipleAttributesTimeout is an method used for setting multiple user attributes
func (c Client) SetMultipleAttributesTimeout(ctx context.Context, timeout time.Duration,
	userID int, attr Attributes) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.SetMultipleAttributes(timeoutCtx, userID, attr)
}

// SetMultipleAttributes is an method used for setting multiple user attributes
func (c Client) SetMultipleAttributes(ctx context.Context, userID int,
	attr Attributes) error {
	payload, err := json.Marshal(attr)
	if err != nil {
		return err
	}
	requestBody := bytes.NewBuffer(payload)
	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf(setMultipleAttributes, userID), requestBody)
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
	resp.Body.Close()

	return statusErrors[resp.StatusCode]
}
