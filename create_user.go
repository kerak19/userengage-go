package userengage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// CreateUser is an struct used for creation of user
type CreateUser struct {
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone_number,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Status    int       `json:"status,omitempty"`
	Gender    int       `json:"gender,omitempty"`
	LastIP    string    `json:"last_ip,omitempty"`
	FirstSeen time.Time `json:"first_seen,omitempty"`
	LastSeen  time.Time `json:"last_seen,omitempty"`
	City      string    `json:"city,omitempty"`
	Region    string    `json:"region,omitempty"`
	Country   string    `json:"country,omitempty"`
	Facebook  string    `json:"facebook_url,omitempty"`
	Linkedin  string    `json:"linkedin_url,omitempty"`
	Twitter   string    `json:"twitter_url,omitempty"`
	Google    string    `json:"google_url,omitempty"`
}

// CreateUserResponse is an struct containing response from userengage
type CreateUserResponse struct {
	ID     int              `json:"id"`
	Errors *json.RawMessage `json:"errors"`
}

// CreateUserTimeout creates user with data provided in CreateUser struct
func (c Client) CreateUserTimeout(ctx context.Context, timeout time.Duration,
	user CreateUser) (CreateUserResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.CreateUser(timeoutCtx, user)
}

const createUserEndpoint = "/users/"

// CreateUser creates user with data provided in CreateUser struct
func (c Client) CreateUser(ctx context.Context, user CreateUser) (CreateUserResponse, error) {
	payload, err := json.Marshal(user)
	if err != nil {
		return CreateUserResponse{}, err
	}

	// https://app.userengage.com/api/public/users/
	endpoint := c.apiPrefix + createUserEndpoint

	body := bytes.NewBuffer(payload)
	r, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return CreateUserResponse{}, err
	}

	return c.requestUserCreation(r.WithContext(ctx))
}

// CreateOrUpdateUser is an struct used for creation or updating an user
type CreateOrUpdateUser struct {
	UserID string `json:"user_id"` // custom user ID
	CreateUser
}

// CreateOrUpdateUserTimeout creates user with provided UserID. If user with this ID
// already exist, it'll update it instead
func (c Client) CreateOrUpdateUserTimeout(ctx context.Context, timeout time.Duration,
	user CreateOrUpdateUser) (CreateUserResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.CreateOrUpdateUser(timeoutCtx, user)
}

const createOrUpdateUserEndpoint = "/users/update_or_create/"

// CreateOrUpdateUser creates user with provided UserID. If user with this ID
// already exist, it'll update it instead
func (c Client) CreateOrUpdateUser(ctx context.Context, user CreateOrUpdateUser) (CreateUserResponse, error) {
	payload, err := json.Marshal(user)
	if err != nil {
		return CreateUserResponse{}, err
	}

	// https://app.userengage.com/api/public/users/update_or_create/
	endpoint := c.apiPrefix + createOrUpdateUserEndpoint

	body := bytes.NewBuffer(payload)
	r, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return CreateUserResponse{}, err
	}

	return c.requestUserCreation(r.WithContext(ctx))
}

func (c Client) requestUserCreation(r *http.Request) (CreateUserResponse, error) {
	var createResponse CreateUserResponse

	client := http.Client{}

	r.Header.Set("Authorization", "Token "+c.apiKey)
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)
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
