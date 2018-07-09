package userengage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

	request.Header.Set("Authorization", "Token "+c.apikey)
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request.WithContext(ctx))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	var setAttributesResponse setAttributesResponse
	err = json.NewDecoder(resp.Body).Decode(&setAttributesResponse)
	if err != nil {
		return err
	}

	if resp.StatusCode == 400 && setAttributesResponse.Errors != nil {
		return errors.New(string(*setAttributesResponse.Errors))
	}
	return statusErrors[resp.StatusCode]
}
