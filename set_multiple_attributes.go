package userengage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const setMultipleAttributes = "https://app.userengage.com/api/public/users/%d/set_multiple_attributes/"

// Attributes contains attributes which are going to be changed
type Attributes map[string]interface{}

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

	request.Header.Set("Authorization", "Token "+c.apikey)
	request.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(request.WithContext(ctx))
	return err
}
