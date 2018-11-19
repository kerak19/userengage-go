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

// Attribute represents single attribute
type Attribute struct {
	Attribute string      `json:"attribute"`
	Value     interface{} `json:"value"`
}

// SetAttributeTimeout set's provided attribute for provided user
func (c Client) SetAttributeTimeout(ctx context.Context, timeout time.Duration, userID int, attr Attribute) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.SetAttribute(timeoutCtx, userID, attr)
}

const setAttributeEndpoint = "/users/%d/set_attribute/"

// SetAttribute set's provided attribute for provided user
func (c Client) SetAttribute(ctx context.Context, userID int, attr Attribute) error {
	payload, err := json.Marshal(attr)
	if err != nil {
		return err
	}

	// https://app.userengage.com/api/public/users/:id/set_attribute/
	endpoint := c.apiPrefix + setAttributeEndpoint

	body := bytes.NewBuffer(payload)
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf(endpoint, userID), body)
	if err != nil {
		return err
	}

	return c.requestSetAttribute(r.WithContext(ctx))
}

// SetAttributeCustomUserIDTimeout set's provided attribute for provided
// user, but using custom user ID
func (c Client) SetAttributeCustomUserIDTimeout(ctx context.Context, timeout time.Duration, customUserID string, attr Attribute) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.SetAttributeCustomUserID(timeoutCtx, customUserID, attr)
}

const setAttributeWithCustomUserIDEndpoint = "/users-by-id/%s/set_attribute/"

// SetAttributeCustomUserID set's provided attribute for provided user, but using custom
// user ID
func (c Client) SetAttributeCustomUserID(ctx context.Context, customUserID string,
	attr Attribute) error {
	payload, err := json.Marshal(attr)
	if err != nil {
		return errors.WithMessage(err, "error while marshaling attribute")
	}

	// https://app.userengage.com/api/public/users-by-id/:user_id/set_attribute/
	endpoint := c.apiPrefix + setAttributeWithCustomUserIDEndpoint

	body := bytes.NewBuffer(payload)
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf(endpoint, customUserID), body)
	if err != nil {
		return errors.WithMessage(err, "error while creating set attribute request")
	}

	return c.requestSetAttribute(r.WithContext(ctx))
}

func (c Client) requestSetAttribute(r *http.Request) error {
	client := http.Client{}

	r.Header.Set("Authorization", "Token "+c.apiKey)
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return errors.WithMessage(err, "error while set attribute request")
	}
	resp.Body.Close()

	return statusErrors[resp.StatusCode]
}
