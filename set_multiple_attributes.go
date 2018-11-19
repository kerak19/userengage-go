package userengage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Attributes contains attributes which are going to be changed
type Attributes map[string]interface{}

// SetMultipleAttributesTimeout is an method used for setting multiple user attributes
func (c Client) SetMultipleAttributesTimeout(ctx context.Context, timeout time.Duration, userID int, attr Attributes) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.SetMultipleAttributes(timeoutCtx, userID, attr)
}

const setMultipleAttributesEndpoint = "/users/%d/set_multiple_attributes/"

// SetMultipleAttributes is an method used for setting multiple user attributes
func (c Client) SetMultipleAttributes(ctx context.Context, userID int, attr Attributes) error {
	payload, err := json.Marshal(attr)
	if err != nil {
		return err
	}

	// https://app.userengage.com/api/public/users/:id/set_multiple_attributes/
	endpoint := c.apiPrefix + setMultipleAttributesEndpoint

	requestBody := bytes.NewBuffer(payload)
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf(endpoint, userID), requestBody)
	if err != nil {
		return err
	}

	return c.requestSetMultipleAttributes(r.WithContext(ctx))
}

// SetMultipleAttributesCustomUserIDTimeout is an method used for setting multiple user attributes
func (c Client) SetMultipleAttributesCustomUserIDTimeout(ctx context.Context, timeout time.Duration, customUserID string, attr Attributes) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.SetMultipleAttributesCustomUserID(timeoutCtx, customUserID, attr)
}

const setMultipleAttributesCustomUserIDEndpoint = "/users-by-id/%s/set_multiple_attributes/"

// SetMultipleAttributesCustomUserID is an method used for setting multiple user attributes,
// but using custom user ID
func (c Client) SetMultipleAttributesCustomUserID(ctx context.Context, customUserID string, attr Attributes) error {
	payload, err := json.Marshal(attr)
	if err != nil {
		return errors.WithMessage(err, "error while marshaling multiple attributes")
	}

	endpoint := c.apiPrefix + setMultipleAttributesCustomUserIDEndpoint

	body := bytes.NewBuffer(payload)
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf(endpoint, customUserID), body)
	if err != nil {
		return errors.WithMessage(err, "error while creating set multiple attributes request")
	}

	return c.requestSetMultipleAttributes(r.WithContext(ctx))
}

func (c Client) requestSetMultipleAttributes(r *http.Request) error {
	client := http.Client{}

	r.Header.Set("Authorization", "Token "+c.apiKey)
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return statusErrors[resp.StatusCode]
}
