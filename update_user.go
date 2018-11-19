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

// UpdateUser is an type used for Updating user information
type UpdateUser = CreateUser

// updateUserResponse is an struct containing response from userengage update user endpoint
type updateUserResponse = CreateUserResponse

// UpdateUserTimeout is an method used for updating user information
func (c Client) UpdateUserTimeout(ctx context.Context, timeout time.Duration, userID int,
	user UpdateUser) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.UpdateUser(timeoutCtx, userID, user)
}

const updateUserEndpoint = "/users/%d/"

// UpdateUser is an method used for updating user information
func (c Client) UpdateUser(ctx context.Context, userID int, user UpdateUser) error {
	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// https://app.userengage.com/api/public/users/:id/
	endpoint := c.apiPrefix + updateUserEndpoint

	body := bytes.NewBuffer(payload)
	r, err := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, userID), body)
	if err != nil {
		return err
	}

	return c.requestUserUpdate(r.WithContext(ctx))
}

// UpdateUserCustomUserID is an type used for Updating user information
type UpdateUserCustomUserID = CreateOrUpdateUser

// UpdateUserCustomUserIDTimeout is an method used for updating user information
func (c Client) UpdateUserCustomUserIDTimeout(ctx context.Context, timeout time.Duration, user UpdateUserCustomUserID) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.UpdateUserCustomUserID(timeoutCtx, user)
}

const updateUserCustomUserIDEndpoint = "/users-by-id/%s/"

// UpdateUserCustomUserID updates user's information using custom user ID
func (c Client) UpdateUserCustomUserID(ctx context.Context, user UpdateUserCustomUserID) error {
	payload, err := json.Marshal(user)
	if err != nil {
		return errors.WithMessage(err, "error while marshaling update user payload")
	}

	endpoint := c.apiPrefix + updateUserCustomUserIDEndpoint

	body := bytes.NewBuffer(payload)
	r, err := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, user.UserID), body)
	if err != nil {
		return err
	}

	return c.requestUserUpdate(r.WithContext(ctx))
}

func (c Client) requestUserUpdate(r *http.Request) error {
	client := http.Client{}

	r.Header.Set("Authorization", "Token "+c.apiKey)
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var updateResponse updateUserResponse
	err = json.NewDecoder(resp.Body).Decode(&updateResponse)
	if err != nil {
		return err
	}

	if resp.StatusCode == 400 && updateResponse.Errors != nil {
		return errors.New(string(*updateResponse.Errors))
	}
	return statusErrors[resp.StatusCode]
}
