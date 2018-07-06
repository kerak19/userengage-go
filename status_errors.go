package userengage

import "errors"

var statusErrors = map[int]error{
	401: errors.New("unauthorized api key"),
	404: errors.New("resource not found"),
	500: errors.New("internal server error"),
}
