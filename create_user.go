package userengage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

const createUserEndpoint = "https://app.userengage.com/api/public/users/"

// CreateUser is an struct used for creation of user
type CreateUser struct {
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone_number,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Gender    *int   `json:"gender"`
	LastIP    string `json:"last_ip,omitempty"`
	City      string `json:"city,omitempty"`
	Region    string `json:"region,omitempty"`
	Country   string `json:"country,omitempty"`
	Facebook  string `json:"facebook_url,omitempty"`
	Linkedin  string `json:"linkedin_url,omitempty"`
	Twitter   string `json:"twitter_url,omitempty"`
	Google    string `json:"google_url,omitempty"`
}

// CreateUserResponse is an struct containing response from userengage
type CreateUserResponse struct {
	ID     int              `json:"id"`
	Errors *json.RawMessage `json:"errors"`
}

// CreateUser creates user with data provided in CreateUser struct
func (c Client) CreateUser(ctx context.Context, user CreateUser) (CreateUserResponse, error) {
	var createResponse CreateUserResponse

	payload, err := json.Marshal(user)
	if err != nil {
		return createResponse, err
	}

	requestBody := bytes.NewBuffer(payload)
	request, err := http.NewRequest(http.MethodPost, createUserEndpoint, requestBody)
	if err != nil {
		return createResponse, err
	}

	request.Header.Set("Authorization", "Token "+c.apikey)
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request.WithContext(ctx))
	if err != nil {
		return createResponse, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&createResponse)
	if err != nil {
		return createResponse, err
	}

	if resp.StatusCode == 400 && createResponse.Errors != nil {
		return createResponse, errors.New(string(*createResponse.Errors))
	}

	return createResponse, statusErrors[resp.StatusCode]
}
