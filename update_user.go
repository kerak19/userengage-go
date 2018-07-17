package userengage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const updateUserEndpoint = "https://app.userengage.com/api/public/users/%d/"

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

// UpdateUser is an method used for updating user information
func (c Client) UpdateUser(ctx context.Context, userID int, user UpdateUser) error {
	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}
	requestBody := bytes.NewBuffer(payload)
	request, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf(updateUserEndpoint, userID), requestBody)
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
