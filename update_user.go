package userengage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const updateUserEndpoint = "https://app.userengage.com/api/public/users/%d/"

// UpdateUser is an type used for Updating user information
type UpdateUser = CreateUser

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

	request.Header.Set("Authorization", "Token "+c.apikey)
	request.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(request.WithContext(ctx))
	return err
}
