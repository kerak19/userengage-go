package userengage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const setAttributeEndpoint = "https://app.userengage.com/api/public/users/%d/set_attribute/"

// Attribute type is used for setting custom user attributes
type Attribute struct {
	Attribute string      `json:"attribute"`
	Value     interface{} `json:"value"`
}

// setAttributesResponse contains possible errors in request
type setAttributesResponse struct {
	Errors *json.RawMessage `json:"errors"`
}

// SetAttributeTimeout set's provided attribute for provided user
func (c Client) SetAttributeTimeout(ctx context.Context, timeout time.Duration, userID int,
	attr Attribute) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.SetAttribute(timeoutCtx, userID, attr)
}

// SetAttribute set's provided attribute for provided user
func (c Client) SetAttribute(ctx context.Context, userID int, attr Attribute) error {
	payload, err := json.Marshal(attr)
	if err != nil {
		return err
	}

	requestBody := bytes.NewBuffer(payload)
	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf(setAttributeEndpoint, userID), requestBody)
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
